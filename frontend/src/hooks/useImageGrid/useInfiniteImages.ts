"use client";

import { useState, useCallback, useRef, useEffect, useMemo } from 'react';
import { fetchJobMetadata } from '@/services/api';
import { ImagePage } from '@/types/HomeImageGrid';
import { useImagePageLoader } from './useImagePageLoader';
import { logger } from '@/utils/logger';

// ============================================================================
// Types
// ============================================================================

interface PageInfo {
  item_dataset_name?: string;
}

// ============================================================================
// Constants
// ============================================================================

const MAX_CACHED_PAGES = 10; // Maximum number of pages to keep in memory

// ============================================================================
// Hook: useInfiniteImages
// ============================================================================

export function useInfiniteImages(
  selectedPages: string,
  selectedDataset: string,
  DatasetList: string[],
  selectedPageIndex: number,
  setSelectedDataset: (dataset: string) => void
) {
  // ========== State ==========
  const [totalPages, setTotalPages] = useState<number>(0);
  const [pagesData, setPagesData] = useState<Map<number, ImagePage>>(new Map());
  const [loading, setLoading] = useState<boolean>(false);
  const [allPagesInfo, setAllPagesInfo] = useState<PageInfo[]>([]);

  // ========== Refs ==========
  const currentJobDatasetKeyRef = useRef<string>('');
  const isInitializingRef = useRef<boolean>(false);
  const loadingPromisesRef = useRef<Map<number, Promise<unknown>>>(new Map());
  const pageAccessOrderRef = useRef<number[]>([]); // Track page access order for LRU
  const loadedPagesRef = useRef<Set<number>>(new Set()); // Track successfully loaded pages

  // ========== Custom Hooks ==========
  const { loadPageImages, clearLoadedPages, removePageFromCache } = useImagePageLoader();

  // ========== Helper Functions ==========
  
  /**
   * Clean up old pages using LRU (Least Recently Used) strategy
   * Keeps only the most recent MAX_CACHED_PAGES pages
   */
  const cleanupOldPages = useCallback((currentPageIndex: number) => {
    // Update access order
    pageAccessOrderRef.current = pageAccessOrderRef.current.filter(idx => idx !== currentPageIndex);
    pageAccessOrderRef.current.push(currentPageIndex);

    // If we have too many pages, remove the oldest ones
    if (pageAccessOrderRef.current.length > MAX_CACHED_PAGES) {
      const pagesToRemove = pageAccessOrderRef.current.slice(0, pageAccessOrderRef.current.length - MAX_CACHED_PAGES);
      
      setPagesData(prev => {
        const newMap = new Map(prev);
        pagesToRemove.forEach(pageIndex => {
          const pageData = newMap.get(pageIndex);
          if (pageData) {
            // Remove from cache
            removePageFromCache(pageData.dataset, pageIndex);
            newMap.delete(pageIndex);
            loadedPagesRef.current.delete(pageIndex); // Remove from loaded tracking
            logger.log(`[Memory] Cleaned up page ${pageIndex} (dataset: ${pageData.dataset})`);
          }
        });
        return newMap;
      });

      // Update access order
      pageAccessOrderRef.current = pageAccessOrderRef.current.filter(
        idx => !pagesToRemove.includes(idx)
      );
      
      logger.log(`[Memory] Total cached pages: ${pageAccessOrderRef.current.length}/${MAX_CACHED_PAGES}`);
    }
  }, [removePageFromCache]);

  // ========== Memoized Values ==========
  const jobDatasetKey = useMemo(
    () => `${selectedPages}-${DatasetList.join(',')}-${DatasetList.length}`,
    [selectedPages, DatasetList]
  );

  const hasPages = useMemo(() => totalPages > 0, [totalPages]);

  // ========== Initialize Total Pages ==========
  const initializeTotalPages = useCallback(async () => {
    if (!selectedPages) {
      logger.log('[useInfiniteImages] initializeTotalPages: No selectedPages, clearing state');
      setTotalPages(0);
      setAllPagesInfo([]);
      return;
    }

    logger.log('[useInfiniteImages] initializeTotalPages: Starting for job:', selectedPages);

    try {
      logger.log('[useInfiniteImages] Calling fetchJobMetadata for:', selectedPages);
      const response = await fetchJobMetadata(selectedPages);
      logger.log('[useInfiniteImages] fetchJobMetadata response:', response);
      
      const totalPagesCount = response.total_pages || 0;
      const datasetNames = response.dataset_names || [];

      logger.log('[useInfiniteImages] Setting totalPages to:', totalPagesCount);
      logger.log('[useInfiniteImages] Dataset list length:', datasetNames.length);
      setTotalPages(totalPagesCount);
      
      // Build allPagesInfo from dataset_names (one per page)
      const pagesInfo = datasetNames.map((datasetName) => ({
        item_dataset_name: datasetName
      }));
      setAllPagesInfo(pagesInfo);
      logger.log('[useInfiniteImages] Built allPagesInfo with', pagesInfo.length, 'pages');

      // Auto-select first dataset if none selected
      if (datasetNames.length > 0 && !selectedDataset && datasetNames[0]) {
        logger.log('[useInfiniteImages] Auto-selecting first dataset:', datasetNames[0]);
        setSelectedDataset(datasetNames[0]);
      } else {
        logger.log('[useInfiniteImages] Not auto-selecting dataset. Datasets:', datasetNames.length, 'selectedDataset:', selectedDataset);
      }
    } catch (error) {
      logger.error('[useInfiniteImages] Error fetching total pages:', error);
      setTotalPages(0);
      setAllPagesInfo([]);
    }
  }, [selectedPages, selectedDataset, setSelectedDataset]);

  // ========== Load Single Page ==========
  const loadPage = useCallback(
    async (pageIndex: number): Promise<ImagePage | null> => {
      logger.log(`[useInfiniteImages] loadPage called for pageIndex: ${pageIndex}`);
      
      // Validation
      if (!selectedPages || pageIndex < 0 || pageIndex >= totalPages) {
        logger.log(`[useInfiniteImages] loadPage validation failed:`, { selectedPages, pageIndex, totalPages });
        return null;
      }

      // Skip if already loading
      if (loadingPromisesRef.current.has(pageIndex)) {
        logger.log(`[useInfiniteImages] Skipping page ${pageIndex} - already loading`);
        return null;
      }
      
      // Skip if already successfully loaded
      if (loadedPagesRef.current.has(pageIndex)) {
        logger.log(`[useInfiniteImages] Skipping page ${pageIndex} - already loaded`);
        return pagesData.get(pageIndex) || null;
      }

      logger.log(`[useInfiniteImages] Loading page ${pageIndex} for job: ${selectedPages}`);

      // Determine target dataset
      let targetDataset = selectedDataset;
      if (allPagesInfo.length > pageIndex) {
        const pageInfo = allPagesInfo[pageIndex];
        if (pageInfo?.item_dataset_name) {
          targetDataset = pageInfo.item_dataset_name;
          logger.log(`[useInfiniteImages] Page ${pageIndex} dataset from pageInfo:`, targetDataset);
          if (targetDataset !== selectedDataset) {
            logger.log(`[useInfiniteImages] Changing selectedDataset from "${selectedDataset}" to "${targetDataset}"`);
            setSelectedDataset(targetDataset);
          }
        }
      } else {
        logger.log(`[useInfiniteImages] No pageInfo available for pageIndex ${pageIndex}, using selectedDataset:`, selectedDataset);
      }

      if (!targetDataset) {
        logger.error(`[useInfiniteImages] No dataset found for page ${pageIndex}`);
        return null;
      }

      // Set loading placeholder
      setPagesData(prev =>
        new Map(prev).set(pageIndex, {
          dataset: targetDataset,
          images: [],
          isNewDataset: false,
        })
      );

      // Start loading
      logger.log(`[useInfiniteImages] Starting loadPageImages for page ${pageIndex}, job: ${selectedPages}, dataset: ${targetDataset}`);
      const loadPromise = loadPageImages(selectedPages, targetDataset, pageIndex);
      loadingPromisesRef.current.set(pageIndex, loadPromise);

      try {
        const result = await loadPromise;
        logger.log(`[useInfiniteImages] loadPageImages result for page ${pageIndex}:`, result);
        
        if (result?.imagePage) {
          logger.log(`[useInfiniteImages] Setting pagesData for page ${pageIndex} with ${result.imagePage.images.length} images`);
          setPagesData(prev => new Map(prev).set(pageIndex, result.imagePage));
          loadedPagesRef.current.add(pageIndex); // Mark as successfully loaded (even if empty)
          
          // Clean up old pages after loading new one
          cleanupOldPages(pageIndex);
          
          if (result.isEmpty) {
            logger.log(`[useInfiniteImages] Page ${pageIndex} successfully loaded but is empty`);
          }
          
          return result.imagePage;
        } else {
          logger.log(`[useInfiniteImages] loadPageImages returned null for page ${pageIndex}`);
          // If load returned null (cancelled or validation failed), remove placeholder
          // Don't mark as loaded - allow retry later
          setPagesData(prev => {
            const newMap = new Map(prev);
            newMap.delete(pageIndex);
            return newMap;
          });
          logger.log(`[InfiniteImages] Page ${pageIndex} load returned null (cancelled or failed validation)`);
        }
      } catch (error) {
        logger.error(`Error loading page ${pageIndex}:`, error);
        // Remove failed page from cache
        setPagesData(prev => {
          const newMap = new Map(prev);
          newMap.delete(pageIndex);
          return newMap;
        });
      } finally {
        loadingPromisesRef.current.delete(pageIndex);
      }

      return null;
    },
    [selectedPages, totalPages, selectedDataset, allPagesInfo, loadPageImages, setSelectedDataset, cleanupOldPages, pagesData]
  );

  // ========== Preload Adjacent Pages ==========
  const preloadAdjacentPages = useCallback(
    (centerIndex: number, range: number = 1) => {
      if (totalPages === 0) return;

      const startIndex = Math.max(0, centerIndex - range);
      const endIndex = Math.min(totalPages - 1, centerIndex + range);

      logger.log(`[useInfiniteImages] Preloading pages ${startIndex} to ${endIndex} around center ${centerIndex}`);

      for (let i = startIndex; i <= endIndex; i++) {
        // Skip current page (already loaded in main effect)
        if (i === centerIndex) continue;
        
        // Skip if already loaded or loading
        if (!loadingPromisesRef.current.has(i) && !loadedPagesRef.current.has(i)) {
          logger.log(`[useInfiniteImages] Preloading page ${i}`);
          loadPage(i);
        } else {
          logger.log(`[useInfiniteImages] Skipping preload for page ${i} (already loaded or loading)`);
        }
      }
    },
    [totalPages, loadPage]
  );

  // ========== Get Page Data ==========
  const getPageData = useCallback(
    (pageIndex: number) => {
      const page = pagesData.get(pageIndex);
      const isLoading = loadingPromisesRef.current.has(pageIndex);

      // Get dataset for this page
      let dataset = selectedDataset;
      if (allPagesInfo.length > pageIndex) {
        const pageInfo = allPagesInfo[pageIndex];
        if (pageInfo?.item_dataset_name) {
          dataset = pageInfo.item_dataset_name;
        }
      }

      // If actively loading, show loading state
      if (isLoading) {
        return {
          pageIndex,
          dataset: dataset || 'Unknown',
          images: [],
          isLoading: true,
          isEmpty: false,
        };
      }

      // If page exists with images, show them
      if (page && page.images.length > 0) {
        // Update dataset in images to ensure they have the correct current dataset
        // This is important when selectedDataset changes but page is already loaded
        const updatedImages = page.images.map(img => ({
          ...img,
          dataset: dataset || img.dataset
        }));
        
        return {
          pageIndex,
          dataset: dataset || 'Unknown',
          images: updatedImages,
          isLoading: false,
          isEmpty: false,
        };
      }

      // If successfully loaded but empty
      if (loadedPagesRef.current.has(pageIndex)) {
        return {
          pageIndex,
          dataset: dataset || 'Unknown',
          images: [],
          isLoading: false,
          isEmpty: true,
        };
      }

      // Not yet loaded - show empty state (will trigger load in PageItem)
      return {
        pageIndex,
        dataset: dataset || 'Unknown',
        images: [],
        isLoading: false,
        isEmpty: false,
      };
    },
    [pagesData, selectedDataset, allPagesInfo]
  );

  // ========== Reset Images ==========
  const resetImages = useCallback(() => {
    loadingPromisesRef.current.clear();
    loadedPagesRef.current.clear(); // Clear loaded pages tracking
    setPagesData(new Map());
    setTotalPages(0);
    setAllPagesInfo([]);
    setLoading(false);
    currentJobDatasetKeyRef.current = '';
    pageAccessOrderRef.current = []; // Clear LRU tracking
    clearLoadedPages();
    isInitializingRef.current = false;
    logger.log('[Memory] Reset all pages and cache');
  }, [clearLoadedPages]);

  // ========== Initialize Images ==========
  const initializeImages = useCallback(async () => {
    if (!selectedPages || DatasetList.length === 0) {
      resetImages();
      return;
    }

    // Skip if already initialized with same data
    if (currentJobDatasetKeyRef.current === jobDatasetKey && totalPages > 0) {
      return;
    }

    // Prevent concurrent initializations
    if (isInitializingRef.current) {
      return;
    }

    try {
      isInitializingRef.current = true;
      setLoading(true);

      // Clear previous data
      setPagesData(new Map());
      loadedPagesRef.current.clear(); // Clear loaded pages tracking
      clearLoadedPages();

      // Fetch metadata
      await initializeTotalPages();

      // Mark as initialized
      currentJobDatasetKeyRef.current = jobDatasetKey;
    } catch (error) {
      logger.error('Error initializing images:', error);
    } finally {
      setLoading(false);
      isInitializingRef.current = false;
    }
  }, [selectedPages, DatasetList, jobDatasetKey, totalPages, resetImages, clearLoadedPages, initializeTotalPages]);

  // ========== Check if Page is Loaded ==========
  const isPageLoaded = useCallback(
    (pageIndex: number) => pagesData.has(pageIndex),
    [pagesData]
  );

  // ========== Effects ==========

  // Auto-initialize when job/dataset changes
  useEffect(() => {
    if (selectedPages && DatasetList.length > 0) {
      if (currentJobDatasetKeyRef.current !== jobDatasetKey) {
        initializeImages();
      }
    } else if (!selectedPages) {
      resetImages();
    }
  }, [selectedPages, DatasetList, jobDatasetKey, initializeImages, resetImages]);

  // Load current page and adjacent pages with debounce
  useEffect(() => {
    if (!(totalPages > 0 && selectedPageIndex >= 0 && selectedPageIndex < totalPages)) {
      return;
    }

    // Always load current page immediately
    loadPage(selectedPageIndex);
    
    // Debounce preloading to avoid excessive requests during rapid navigation
    const preloadTimer = setTimeout(() => {
      preloadAdjacentPages(selectedPageIndex, 1); // Load ±1 page after 300ms
    }, 300);
    
    return () => clearTimeout(preloadTimer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedPageIndex, totalPages]);

  // Load first page on initialization
  useEffect(() => {
    if (!(totalPages > 0 && selectedPageIndex === 0 && !pagesData.has(0))) {
      return;
    }

    loadPage(0);
    // Debounce preloading for initialization as well
    const preloadTimer = setTimeout(() => {
      preloadAdjacentPages(0, 1);
    }, 300);
    
    return () => clearTimeout(preloadTimer);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [totalPages, selectedPageIndex]);

  // ========== Return ==========
  return {
    totalPages,
    loading,
    getPageData,
    loadPage,
    preloadAdjacentPages,
    resetImages,
    initializeImages,
    hasPages,
    isPageLoaded,
  };
}

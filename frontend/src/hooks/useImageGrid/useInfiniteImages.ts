"use client";

import { useState, useCallback, useRef, useEffect } from 'react';
import { fetchALLPages } from '@/services/api';
import { ImagePage } from '@/types/HomeImageGrid';
import { useImagePageLoader } from './useImagePageLoader';

export function useInfiniteImages(
  selectedPages: string, 
  selectedDataset: string, 
  DatasetList: string[], 
  selectedPageIndex: number,
  setSelectedDataset: (dataset: string) => void
) {
  const [totalPages, setTotalPages] = useState<number>(0);
  const [pagesData, setPagesData] = useState<Map<number, ImagePage>>(new Map());
  const [loading, setLoading] = useState<boolean>(false);
  const [currentJobDatasetKey, setCurrentJobDatasetKey] = useState<string>('');
  const [allPagesInfo, setAllPagesInfo] = useState<{ item_dataset_name?: string }[]>([]);
  
  const isInitializingRef = useRef<boolean>(false);
  const loadingPromises = useRef<Map<number, Promise<unknown>>>(new Map());

  const { loadPageImages, clearLoadedPages } = useImagePageLoader();

  const initializeTotalPages = useCallback(async () => {
    if (!selectedPages) {
      setTotalPages(0);
      setAllPagesInfo([]);
      return;
    }

    try {
      const response = await fetchALLPages(selectedPages);
      const totalPagesCount = response.total_pages || 0;
      const pages = Array.isArray(response.pages) ? response.pages : Object.values(response.pages || {});
      
      setTotalPages(totalPagesCount);
      setAllPagesInfo(pages);
      
      if (pages.length > 0 && !selectedDataset) {
        const firstPage = pages[0];
        if (firstPage?.item_dataset_name) {
          setSelectedDataset(firstPage.item_dataset_name);
        }
      }
    } catch (error) {
      console.error('Error fetching total pages:', error);
      setTotalPages(0);
      setAllPagesInfo([]);
    }
  }, [selectedPages, selectedDataset, setSelectedDataset]);

  const loadPage = useCallback(async (pageIndex: number): Promise<ImagePage | null> => {
    if (!selectedPages || pageIndex < 0 || pageIndex >= totalPages) {
      return null;
    }
    const existingPage = pagesData.get(pageIndex);
    if (existingPage || loadingPromises.current.has(pageIndex)) {
      return existingPage || null;
    }

    let targetDataset = selectedDataset;
    if (allPagesInfo.length > pageIndex) {
      const pageInfo = allPagesInfo[pageIndex];
      if (pageInfo?.item_dataset_name) {
        targetDataset = pageInfo.item_dataset_name;
        if (targetDataset !== selectedDataset) {
          setSelectedDataset(targetDataset);
        }
      }
    }

    if (!targetDataset) {
      console.error(`No dataset found for page ${pageIndex}`);
      return null;
    }

    setPagesData((prev: Map<number, ImagePage>) => new Map(prev).set(pageIndex, {
      dataset: targetDataset,
      images: [],
      isNewDataset: false
    }));

    const loadPromise = loadPageImages(selectedPages, targetDataset, pageIndex);
    loadingPromises.current.set(pageIndex, loadPromise);

    try {
      const result = await loadPromise;
      if (result?.imagePage) {
        setPagesData((prev: Map<number, ImagePage>) => new Map(prev).set(pageIndex, result.imagePage));
        return result.imagePage;
      }
    } catch (error) {
      console.error(`Error loading page ${pageIndex}:`, error);
      setPagesData((prev: Map<number, ImagePage>) => {
        const newMap = new Map(prev);
        newMap.delete(pageIndex);
        return newMap;
      });
    } finally {
      loadingPromises.current.delete(pageIndex);
    }

    return null;
  }, [selectedPages, totalPages, selectedDataset, allPagesInfo, loadPageImages, setSelectedDataset, pagesData]);

  const preloadAdjacentPages = useCallback((centerIndex: number, range: number = 1) => {
    if (totalPages === 0) return;
    
    const startIndex = Math.max(0, centerIndex - range);
    const endIndex = Math.min(totalPages - 1, centerIndex + range);
    
    for (let i = startIndex; i <= endIndex; i++) {
      if (!pagesData.has(i) && !loadingPromises.current.has(i)) {
        loadPage(i);
      }
    }
  }, [totalPages, pagesData, loadPage]);

  const getPageData = useCallback((pageIndex: number) => {
    const page = pagesData.get(pageIndex);
    const isLoading = loadingPromises.current.has(pageIndex);
    
    // 獲取頁面對應的數據集
    let dataset = selectedDataset;
    if (allPagesInfo.length > pageIndex) {
      const pageInfo = allPagesInfo[pageIndex];
      if (pageInfo?.item_dataset_name) {
        dataset = pageInfo.item_dataset_name;
      }
    }
    
    return {
      pageIndex,
      dataset: dataset || 'Unknown',
      images: page?.images || [],
      isLoading: isLoading && !page,
      isEmpty: !isLoading && !page?.images?.length && page !== undefined
    };
  }, [pagesData, selectedDataset, allPagesInfo]);

  const resetImages = useCallback(() => {
    console.log('Resetting images...');
    
    loadingPromises.current.clear();
    
    setPagesData(new Map());
    setTotalPages(0);
    setAllPagesInfo([]);
    setLoading(false);
    setCurrentJobDatasetKey('');
    clearLoadedPages();
    isInitializingRef.current = false;
  }, [clearLoadedPages]);

  const initializeImages = useCallback(async () => {
    if (!selectedPages || DatasetList.length === 0) {
      resetImages();
      return;
    }

    const newJobDatasetKey = `${selectedPages}-${DatasetList.join(',')}-${DatasetList.length}`;
    
    if (currentJobDatasetKey === newJobDatasetKey && totalPages > 0) {
      return;
    }

    if (isInitializingRef.current) {
      return;
    }

    try {
      console.log('Initializing images for job:', selectedPages);
      isInitializingRef.current = true;
      setLoading(true);
      
      setPagesData(new Map());
      clearLoadedPages();
    
      await initializeTotalPages();
      
      setCurrentJobDatasetKey(newJobDatasetKey);
      
    } catch (error) {
      console.error('Error initializing images:', error);
    } finally {
      setLoading(false);
      isInitializingRef.current = false;
    }
  }, [selectedPages, DatasetList, currentJobDatasetKey, resetImages, clearLoadedPages, initializeTotalPages, totalPages]);

  useEffect(() => {
    if (selectedPages && DatasetList.length > 0) {
      const newJobDatasetKey = `${selectedPages}-${DatasetList.join(',')}-${DatasetList.length}`;
      
      if (currentJobDatasetKey !== newJobDatasetKey) {
        initializeImages();
      }
    } else if (!selectedPages) {
      resetImages();
    }
  }, [selectedPages, DatasetList, currentJobDatasetKey, initializeImages, resetImages]);

  useEffect(() => {
    if (totalPages > 0 && selectedPageIndex >= 0 && selectedPageIndex < totalPages) {
      loadPage(selectedPageIndex);
      preloadAdjacentPages(selectedPageIndex);
    }
  }, [selectedPageIndex, totalPages, loadPage, preloadAdjacentPages]);

  useEffect(() => {
    if (totalPages > 0 && selectedPageIndex === 0 && !pagesData.has(0)) {
      loadPage(0);
      preloadAdjacentPages(0);
    }
  }, [totalPages, selectedPageIndex, loadPage, preloadAdjacentPages, pagesData]);

  return {
    totalPages,
    loading,
    getPageData,
    loadPage,
    preloadAdjacentPages,
    resetImages,
    initializeImages,
    hasPages: totalPages > 0,
    isPageLoaded: (pageIndex: number) => pagesData.has(pageIndex)
  };
}
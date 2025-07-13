"use client";

import { useState, useCallback, useRef, useEffect } from 'react';
import { fetchALLPages } from '@/services/api';
import { ImagePage } from '@/types/HomeImageGrid';
import { usePageObserver } from './usePageObserver';
import { useImagePageLoader } from './useImagePageLoader';

export function useInfiniteImages(
  selectedJob: string | null, selectedDataset: string | null, allDatasets: string[], currentPage: number,
  setSelectedDataset: (dataset: string) => void, setCurrentPage: (page: number) => void
) {

  const [maxPageIndex, setMaxPageIndex] = useState<number>(0); 
  const [allImagePages, setAllImagePages] = useState<ImagePage[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [currentPageIndex, setCurrentPageIndex] = useState<number>(0);
  const [currentDatasetIndex, setCurrentDatasetIndex] = useState<number>(0);
  const [hasMore, setHasMore] = useState<boolean>(true);
  
  const [currentJobDatasetKey, setCurrentJobDatasetKey] = useState<string>('');
  const isInitializingRef = useRef<boolean>(false);

  const { registerPageElement, cleanupObservers } = usePageObserver();
  const { loadPageImages, clearLoadedPages } = useImagePageLoader();

  const updateCurrentPageIndex = useCallback((newPageIndex: number) => {
    if (newPageIndex !== currentPageIndex) {
      setCurrentPageIndex(newPageIndex);
      setCurrentPage(newPageIndex);
    }
  }, [currentPageIndex, setCurrentPage]);

  // This function registers the page element for the infinite scroll functionality
  const registerPageElementWithCallback = useCallback((pageIndex: number, element: HTMLElement | null) => {
    registerPageElement(pageIndex, element, updateCurrentPageIndex);
  }, [registerPageElement, updateCurrentPageIndex]);

  // This effect updates the current dataset index based on the current page index
  useEffect(() => {
    const fetchCurrentDatasetByCurrentPage = async () => {
      try {
        if (currentPageIndex <= 0 || !selectedJob || !selectedDataset) {
          console.warn('Invalid page index or no job selected, skipping dataset detail fetch.');
          return;
        }
        const response = await fetchALLPages(selectedJob);
        const pages = Array.isArray(response.pages) ? response.pages : Object.values(response.pages);
        setSelectedDataset(pages[currentPageIndex]);
        console.log(`Current dataset detail fetched for page index ${currentPageIndex}: ${pages[currentPageIndex]}`);
        
      } catch (error) {
        console.error(`Error fetching details for dataset ${selectedDataset}:`, error);
      }
    };

    fetchCurrentDatasetByCurrentPage();
  }, [currentPageIndex, selectedJob, selectedDataset, setSelectedDataset]);

  // This function resets the images and clears all observers
  const resetImages = useCallback(() => {
    console.log('Resetting images...');
    cleanupObservers();
    setAllImagePages([]);
    setCurrentPageIndex(0);
    setCurrentDatasetIndex(0);
    setHasMore(true);
    setCurrentJobDatasetKey('');
    setMaxPageIndex(0);
    clearLoadedPages();
    isInitializingRef.current = false;
  }, [cleanupObservers, clearLoadedPages]);

  // This function initializes the images for the selected job and datasets
  const initializeImages = useCallback(async () => {
    if (!selectedJob || allDatasets.length === 0) {
      resetImages();
      return;
    }

    const newJobDatasetKey = `${selectedJob}-${allDatasets.join(',')}-${allDatasets.length}`;
    
    if (currentJobDatasetKey === newJobDatasetKey) {
      console.log('Already initialized for this job and datasets, skipping...');
      return;
    }

    if (isInitializingRef.current) {
      console.log('Already initializing, skipping...');
      return;
    }

    try {
      console.log(`Initializing images for job: ${selectedJob}, datasets: ${allDatasets.join(',')}`);
      isInitializingRef.current = true;
      setLoading(true);
      
      cleanupObservers();
      setAllImagePages([]);
      setCurrentPageIndex(0);
      setCurrentDatasetIndex(0);
      setHasMore(true);
      clearLoadedPages();
      
      const firstDataset = allDatasets[0];
      const firstPage = await loadPageImages(selectedJob, firstDataset, 0);
      
      if (firstPage) {
        setAllImagePages([firstPage.imagePage]);
        setCurrentPageIndex(0);
        setCurrentDatasetIndex(0);
        setHasMore(firstPage.maxPage > 0);
        setCurrentJobDatasetKey(newJobDatasetKey);
        setCurrentPage(0);
      } else {
        setAllImagePages([]);
        setHasMore(false);
        setCurrentJobDatasetKey(newJobDatasetKey);
      }
      
    } catch (error) {
      console.error('Error initializing images:', error);
      setAllImagePages([]);
      setHasMore(false);
    } finally {
      setLoading(false);
      isInitializingRef.current = false;
    }
  }, [selectedJob, allDatasets, loadPageImages, currentJobDatasetKey, resetImages, cleanupObservers, clearLoadedPages, setCurrentPage]);

  // This effect initializes images when the selected job or datasets change
  useEffect(() => {
    if (selectedJob && allDatasets.length > 0) {
      const newJobDatasetKey = `${selectedJob}-${allDatasets.join(',')}-${allDatasets.length}`;
      
      if (currentJobDatasetKey !== newJobDatasetKey) {
        console.log('Job or datasets changed, initializing...');
        initializeImages();
      }
    } else if (!selectedJob) {
      resetImages();
    }
  }, [selectedJob, allDatasets, currentJobDatasetKey, initializeImages, resetImages]);

  // This function loads the next page of images for the current dataset
  const loadNextPage = useCallback(async () => {
    if (!selectedJob || loading || !hasMore || currentDatasetIndex >= allDatasets.length || isInitializingRef.current) {
      return;
    }

    try {
      setLoading(true);
      const currentDataset = allDatasets[currentDatasetIndex];
      const nextPageIndex = currentPageIndex + 1;
      
      if (nextPageIndex >= maxPageIndex && maxPageIndex > 0) {
        console.warn(`No more pages available for dataset ${currentDataset} at index ${currentDatasetIndex}`);
        setHasMore(false);
        return;
      }
      
      const response = await loadPageImages(selectedJob, currentDataset, nextPageIndex);
      if (response) {
        const { imagePage, maxPage } = response;
        setMaxPageIndex(maxPage);
        setAllImagePages(prev => [...prev, imagePage]);

        if (nextPageIndex >= maxPage) {
          setHasMore(false);
        }
      } else {
        setHasMore(false);
      }
    } catch (error) {
      console.error('Error loading next page:', error);
      setHasMore(false);
    } finally {
      setLoading(false);
    }
  }, [selectedJob, loading, hasMore, currentDatasetIndex, currentPageIndex, allDatasets, loadPageImages, maxPageIndex]);

  const getCurrentImagePages = useCallback(() => {
    return allImagePages;
  }, [allImagePages]);

  const hasMorePages = useCallback(() => {
    return hasMore;
  }, [hasMore]);

  useEffect(() => {
    return () => {
      cleanupObservers();
    };
  }, [cleanupObservers]);

  return {
    allImagePages,
    loading,
    currentPageIndex,
    getCurrentImagePages,
    hasMorePages,
    loadNextPage,
    resetImages,
    initializeImages,
    registerPageElement: registerPageElementWithCallback
  };
}

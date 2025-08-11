"use client";

import { useState, useCallback, useRef, useEffect } from 'react';
import { fetchALLPages } from '@/services/api';
import { ImagePage } from '@/types/HomeImageGrid';
import { usePageObserver } from './usePageObserver';
import { useImagePageLoader } from './useImagePageLoader';

export function useInfiniteImages(
  selectedPages: string, selectedDataset: string, DatasetList: string[], selectedPageIndex: number,
  setSelectedDataset: (dataset: string) => void, setselectedPageIndex: (page: number) => void
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

  // This function registers the page element for the infinite scroll functionality
  const updateCurrentPageIndex = useCallback((newPageIndex: number) => {
    if (newPageIndex !== currentPageIndex) {
      setCurrentPageIndex(newPageIndex);
      setselectedPageIndex(newPageIndex);
    }
  }, [currentPageIndex, setselectedPageIndex]);

  const registerPageElementWithCallback = useCallback((pageIndex: number, element: HTMLElement) => {
    registerPageElement(pageIndex, element, updateCurrentPageIndex);
  }, [registerPageElement, updateCurrentPageIndex]);

  // This effect updates the current dataset index based on the current page index
  // TODO: can add getdatasetnamebypageidex func 
  useEffect(() => {
    const fetchCurrentDatasetByCurrentPage = async () => {
      try {
        if (currentPageIndex <= 0 || !selectedPages || !selectedDataset) {
          console.warn('Invalid page index or no job selected, skipping dataset detail fetch.');
          return;
        }
        const response = await fetchALLPages(selectedPages);
        console.log(`Fetched ${response}`);
        const pages = Array.isArray(response.pages) ? response.pages : Object.values(response.pages);
        setSelectedDataset(pages[currentPageIndex].item_dataset_name);
        console.log(`Current dataset detail fetched for page index ${currentPageIndex}: ${pages[currentPageIndex].item_dataset_name}`);
        
      } catch (error) {
        console.error(`Error fetching details for dataset ${selectedDataset}:`, error);
      }
    };

    fetchCurrentDatasetByCurrentPage();
  }, [currentPageIndex]);

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
    if (!selectedPages || DatasetList.length === 0) {
      resetImages();
      return;
    }

    const newJobDatasetKey = `${selectedPages}-${DatasetList.join(',')}-${DatasetList.length}`;
    
    if (currentJobDatasetKey === newJobDatasetKey) {
      console.log('Already initialized for this job and datasets, skipping...');
      return;
    }

    if (isInitializingRef.current) {
      console.log('Already initializing, skipping...');
      return;
    }

    try { 
      console.log('Initializing images for job:', selectedPages);
      isInitializingRef.current = true;
      setLoading(true);
      
      cleanupObservers();
      setAllImagePages([]);
      setCurrentPageIndex(0);
      setCurrentDatasetIndex(0);
      setHasMore(true);
      clearLoadedPages();
      
      const firstDataset = DatasetList[0];
      const firstPage = await loadPageImages(selectedPages, firstDataset, 0);
      
      if (firstPage) {
        setAllImagePages([firstPage.imagePage]);
        setCurrentPageIndex(0);
        setCurrentDatasetIndex(0);
        setHasMore(firstPage.maxPage > 0);
        setCurrentJobDatasetKey(newJobDatasetKey);
        setselectedPageIndex(0);
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
  }, [selectedPages, DatasetList, loadPageImages, currentJobDatasetKey, resetImages, cleanupObservers, clearLoadedPages, setselectedPageIndex]);

  // This effect initializes images when the selected job or datasets change
  useEffect(() => {
    if (selectedPages && DatasetList.length > 0) {
      const newJobDatasetKey = `${selectedPages}-${DatasetList.join(',')}-${DatasetList.length}`;
      
      if (currentJobDatasetKey !== newJobDatasetKey) {
        console.log('Job or datasets changed, initializing...');
        initializeImages();
      }
    } else if (!selectedPages) {
      resetImages();
    }
  }, [selectedPages, DatasetList, currentJobDatasetKey, initializeImages, resetImages]);

  // This function loads the next page of images for the current dataset
  const loadNextPage = useCallback(async () => {
    if (!selectedPages || loading || !hasMore || currentDatasetIndex >= DatasetList.length || isInitializingRef.current) {
      return;
    }

    try {
      setLoading(true);
      const currentDataset = DatasetList[currentDatasetIndex];
      const nextPageIndex = currentPageIndex + 1;
      
      if (nextPageIndex >= maxPageIndex && maxPageIndex > 0) {
        console.warn(`No more pages available for dataset ${currentDataset} at index ${currentDatasetIndex}`);
        setHasMore(false);
        return;
      }
      
      const response = await loadPageImages(selectedPages, currentDataset, nextPageIndex);
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
  }, [selectedPages, loading, hasMore, currentDatasetIndex, currentPageIndex, DatasetList, loadPageImages, maxPageIndex]);

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

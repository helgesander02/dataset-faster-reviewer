"use client";

import { useState, useEffect, useCallback, useRef } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';
import { SidebarState, SidebarActions } from '@/types/HomeLeftSidebar';
import { fetchJobs, fetchDatasets } from '@/services/api';
import { DATASET_PER_PAGE } from '@/services/config';

export function useLeftSidebar(): SidebarState & SidebarActions {
  const { 
    selectedPages, selectedPageIndex,
    setSelectedDataset, setselectedPageIndex 
  } = useJobDataset();

  const [currentJobList, setJobList] = useState<string[]>([]);
  const [currentDatasetList, setDatasetList] = useState<string[]>([]);
  const [currentPagenation, setCurrentPagenation] = useState<number>(0);  
  const [loading, setLoading] = useState<boolean>(false);

  // Ref to hold the interval ID for auto-refreshing jobs
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  
  const loadJobs = useCallback(async () => {
    try {
      setLoading(true);
      const response = await fetchJobs();
      if (!response || !response.job_names) {
        throw new Error('Invalid response format');
      }
      setJobList(response.job_names);

    } catch (error) {
      console.error('Unable to load jobs:', error);

    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadJobs();
    
    intervalRef.current = setInterval(() => {
      loadJobs();
    }, 60000);

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [loadJobs]);

  // Load datasets when job changes
  useEffect(() => {
    const loadDatasets = async () => {
      if (!selectedPages) return;
      
      try {
        setLoading(true);
        const response = await fetchDatasets(selectedPages);
        if (!response || !response.dataset_names) {
          throw new Error('Invalid response format');
        }
        setDatasetList(response.dataset_names);
        setCurrentPagenation(0);

        if (response.dataset_names.length > 0) {
          setselectedPageIndex(0);
          setSelectedDataset(response.dataset_names[0]);
        }
      } catch (error) {
        throw new Error('Unable to load datasets: ' + error);

      } finally {
        setLoading(false);
      }
    };

    loadDatasets();
  }, [selectedPages, setSelectedDataset, setselectedPageIndex]); 

  // Auto-paginate when selectedPageIndex exceeds current page capacity
  useEffect(() => {
    if (selectedPageIndex < 0) return;
    
    const requiredPage = Math.floor(selectedPageIndex / DATASET_PER_PAGE);
    const totalPages = Math.ceil(currentDatasetList.length / DATASET_PER_PAGE);
    
    if (requiredPage !== currentPagenation && requiredPage < totalPages) {
      setCurrentPagenation(requiredPage);
    }
  }, [selectedPageIndex, currentDatasetList.length, currentPagenation]);

  // Pagination logic
  const previousPage = useCallback(() => {
    if (currentPagenation > 0) {
      const newPage = currentPagenation - 1;
      setCurrentPagenation(newPage);
      
      // Reset selectedPageIndex to first item of the new page
      const newIndex = newPage * DATASET_PER_PAGE;
      if (newIndex < currentDatasetList.length) {
        setselectedPageIndex(newIndex);
        setSelectedDataset(currentDatasetList[newIndex]);
      }
    }
  }, [currentPagenation, currentDatasetList, setselectedPageIndex, setSelectedDataset]);

  const nextPage = useCallback(() => {
    if ((currentPagenation + 1) * DATASET_PER_PAGE < currentDatasetList.length) {
      const newPage = currentPagenation + 1;
      setCurrentPagenation(newPage);
      
      // Reset selectedPageIndex to first item of the new page
      const newIndex = newPage * DATASET_PER_PAGE;
      if (newIndex < currentDatasetList.length) {
        setselectedPageIndex(newIndex);
        setSelectedDataset(currentDatasetList[newIndex]);
      }
    }
  }, [currentPagenation, currentDatasetList, setselectedPageIndex, setSelectedDataset]);

  return {
    currentJobList,
    currentDatasetList,
    currentPagenation,
    loading,
    previousPage,
    nextPage
  };
}

"use client";

import { useState, useEffect } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';
import { SidebarState, SidebarActions } from '@/types/HomeLeftSidebar';
import { fetchJobs, fetchDatasets } from '@/services/api';
import { DATASET_PER_PAGE } from '@/services/config';

export function useLeftSidebar(): SidebarState & SidebarActions {
  const { 
    selectedJob, selectedPages, selectedDataset, selectedPageIndex,
    setSelectedJob, setSelectedPages, setSelectedDataset, setselectedPageIndex 
  } = useJobDataset();

  const [currentJobList, setJobList] = useState<string[]>([]);
  const [currentDatasetList, setDatasetList] = useState<string[]>([]);
  const [currentPagenation, setCurrentPagenation] = useState<number>(0);  
  const [loading, setLoading] = useState<boolean>(false);
  

  // Load jobs on component mount
  useEffect(() => {
    const loadJobs = async () => {
      try {
        setLoading(true);
        const response = await fetchJobs();
        if (!response || !response.job_names) {
          throw new Error('Invalid response format');
        }
        setJobList(response.job_names);

      } catch (error) {
        throw new Error('Unable to load jobs: ' + error);

      } finally {
        setLoading(false);
      }
    };
    loadJobs();
  }, []);

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
  }, [selectedPages]);

  // Pagination logic
  const previousPage = () => {
    if (currentPagenation > 0) {
      setCurrentPagenation(currentPagenation - 1);
    }
  };

  const nextPage = () => {
    if ((currentPagenation + 1) * DATASET_PER_PAGE < currentDatasetList.length) {
      setCurrentPagenation(currentPagenation + 1);
    }
  };

  return {
    currentJobList,
    currentDatasetList,
    currentPagenation,
    loading,
    previousPage,
    nextPage
  };
}
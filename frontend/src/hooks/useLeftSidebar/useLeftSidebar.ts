"use client";

import { useState, useEffect } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';
import { SidebarState, SidebarActions } from '@/types/HomeLeftSidebar';
import { fetchJobs, fetchDatasets, fetchALLPages } from '@/services/api';
import { DATASET_PER_PAGE } from '@/services/config';

export function useLeftSidebar(): SidebarState & SidebarActions {
  const { 
    selectedJob, selectedDataset, currentPage, jobPages,
    setSelectedJob, setSelectedDataset, setCurrentPage 
  } = useJobDataset();

  const [currentJobs, setJobs] = useState<string[]>([]);
  const [currentDatasets, setDatasets] = useState<string[]>([]);  
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
        setJobs(response.job_names);

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
      if (!selectedJob) return;
      
      try {
        setLoading(true);
        const response = await fetchDatasets(selectedJob);
        if (!response || !response.dataset_names) {
          throw new Error('Invalid response format');
        }
        setDatasets(response.dataset_names);
        setCurrentPagenation(0);

        if (response.dataset_names.length > 0) {
          setSelectedDataset(response.dataset_names[0]);
        }
      } catch (error) {
        throw new Error('Unable to load datasets: ' + error);

      } finally {
        setLoading(false);
      }
    };

    loadDatasets();
  }, [selectedJob]);

  // Handle pagination for datasets
  useEffect(() => {
    const changCurrentPage = async () => {
      try {
        if (currentPage <= 0 || !selectedJob || !selectedDataset) {
          console.warn('No job or dataset selected, skipping pagination update.');
          return;
        }

        const firstIndex = jobPages.indexOf(selectedDataset);
        if (firstIndex !== currentPage) {
          setCurrentPage(firstIndex);
          console.log(`LiftSidebar: Page index updated to ${firstIndex} for dataset ${selectedDataset}`);
        }
        console.log(`Current dataset detail fetched for ${selectedDataset} at page index ${firstIndex}`);

      } catch (error) {
        throw new Error('Unable to update current page: ' + error);
      }
    };
  
    changCurrentPage();
  }, [selectedDataset, jobPages]);

  //  Handle job and dataset selection
  const handleJobSelect = (job: string) => {
    setSelectedJob(job);
    setSelectedDataset('');
    console.log(`LeftSidebar: Job selected: ${job}`);
  };

  const handleDatasetSelect = (dataset: string) => {
    setSelectedDataset(dataset);
    console.log(`LeftSidebar: Dataset selected: ${dataset}`);
  };

  // Pagination logic
  const previousPage = () => {
    if (currentPagenation > 0) {
      setCurrentPagenation(currentPagenation - 1);
    }
  };

  const nextPage = () => {
    if ((currentPagenation + 1) * DATASET_PER_PAGE < currentDatasets.length) {
      setCurrentPagenation(currentPagenation + 1);
    }
  };

  return {
    currentJobs,
    currentDatasets,
    currentPagenation,
    loading,
    handleJobSelect,
    handleDatasetSelect,
    previousPage,
    nextPage
  };
}
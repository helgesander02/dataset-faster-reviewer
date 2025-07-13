"use client";

import React, { createContext, useState, useContext, ReactNode, useEffect } from 'react';
import { CachedImage, JobDatasetContextType } from '@/types/JobDatasetContext';
import { getPendingReview, fetchALLPages } from '@/services/api';

const JobDatasetContext = createContext<JobDatasetContextType | undefined>(undefined);

/**
 * JobDatasetProvider component provides context for managing job and dataset state.
 * It fetches pending review images on initial render and updates job pages when selectedJob changes.
 * 
 * props:
 * - children: ReactNode - The child components that will have access to this context.
 * * This context includes:
 * - selectedJob:           string - The currently selected job.
 * - selectedDataset:       string - The currently selected dataset.
 * - currentPage:           number - The current page of datasets being displayed.
 * - setSelectedJob:        function - Function to set the selected job.
 * - setSelectedDataset:    function - Function to set the selected dataset.
 * - setCurrentPage:        function - Function to set the current page.
 * - cachedImages:          CachedImage[] - Array of cached images for the selected job and dataset.
 * - addImageToCache:       function - Function to add an image to the cache.
 * - removeImageFromCache:  function - Function to remove an image from the cache.
 * - getCache:              function - Function to retrieve cached images for a specific job and dataset.
 * - jobPages:              string[] - Array of pages for the selected job.
 */
export function JobDatasetProvider({ children }: { children: ReactNode }) {
  const [selectedJob, setSelectedJob] = useState<string>('');
  const [selectedDataset, setSelectedDataset] = useState<string>('');
  const [currentPage, setCurrentPage] = useState<number>(0);
  const [cachedImages, setCachedImages] = useState<CachedImage[]>([]);
  const [jobPages, setJobPages] = useState<string[]>([]);

  // Load pending review images on initial render
  useEffect(() => {
    async function loadPending() {
      try {
        const data = await getPendingReview(true);
        data.items.forEach((item: CachedImage) => {
          addImageToCache(item.job, item.dataset, item.imageName, item.imagePath);
        });
      } catch (error) {
        throw new Error('Error loading pending review: ' + error);
      }
    }

    loadPending();
  }, []);

  // Fetch job pages when selectedJob changes
  useEffect(() => {
    async function fetchJobPages() {
      if (!selectedJob) {
        setJobPages([]);
        return;
      }

      try {
        await new Promise(resolve => setTimeout(resolve, 1000)); // wait backend to be ready
        const response = await fetchALLPages(selectedJob);
        const pages = Array.isArray(response.pages) ? response.pages : Object.values(response.pages);
        setJobPages(pages);

      } catch (error) {
        setJobPages([]);
        throw new Error('Error fetching job pages: ' + error);
      }
    }

    fetchJobPages();
  }, [selectedJob]);

  // Function to add an image to the cache
  const addImageToCache = (job: string, dataset: string, imageName: string, imagePath: string) => {
    const exists = cachedImages.some(
      img => img.job === job && img.dataset === dataset && img.imagePath === imagePath && img.imageName === imageName
    );

    if (!exists) {
      const newCachedImage: CachedImage = { job, dataset, imageName, imagePath };
      setCachedImages(prev => [newCachedImage, ...prev]);
    }
  };

  // Function to remove an image from the cache
  const removeImageFromCache = (imagePath: string) => {
    setCachedImages(prev => prev.filter(img => img.imagePath !== imagePath));
  };

  const getCache = (job: string, dataset: string) => {
    return cachedImages
      .filter(img => img.job === job && img.dataset === dataset)
      .map(img => img.imagePath);
  };

  return (
    <JobDatasetContext.Provider
      value={{
        selectedJob,
        selectedDataset,
        currentPage,
        setSelectedJob,
        setSelectedDataset,
        setCurrentPage,
        cachedImages,
        addImageToCache,
        removeImageFromCache,
        getCache,
        jobPages,
      }}
    >
      {children}
    </JobDatasetContext.Provider>
  );
}

export function useJobDataset() {
  const context = useContext(JobDatasetContext);
  if (!context) {
    throw new Error('useJobDataset must be used within a JobDatasetProvider');
  }
  return context;
}
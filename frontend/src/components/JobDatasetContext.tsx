"use client";

import React, { createContext, useState, useContext, ReactNode, useEffect, useCallback } from 'react';
import { CachedImage, JobDatasetContextType } from '@/types/JobDatasetContext';
import { getPendingReview, updateALLPages } from '@/services/api';
import { IMAGES_PER_PAGE } from '@/services/config';

const JobDatasetContext = createContext<JobDatasetContextType | undefined>(undefined);

/**
 * JobDatasetProvider component provides context for managing job and dataset state.
 * It fetches pending review images on initial render and updates job pages when selectedJob changes.
 * 
 * props:
 * - children: ReactNode - The child components that will have access to this context.
 * * This context includes:
 * - selectedJob: string - The currently selected job.
 * - selectedPages: string - The currently selected pages for the job.
 * - selectedDataset: string - The currently selected dataset.
 * - selectedPageIndex: number - The index of the currently selected page.
 * - setSelectedJob: (job: string) => void - Function to set the selected job.
 * - setSelectedPages: (pages: string) => void - Function to set the selected pages.
 * - setSelectedDataset: (dataset: string) => void - Function to set the selected dataset.
 * - setselectedPageIndex: (index: number) => void - Function to set the selected page index.
 * - cachedImages: CachedImage[] - Array of cached images.
 * - addImageToCache: (job: string, dataset: string, imageName: string, imagePath: string) => void - Function to add an image to the cache.
 * - removeImageFromCache: (imagePath: string) => void - Function to remove an image from the cache.
 * - getCache: (job: string, dataset: string) => string[] - Function to get cached image paths for a specific job and dataset.s
 */
export function JobDatasetProvider({ children }: { children: ReactNode }) {
  const [selectedJob, setSelectedJob] = useState<string>('');
  const [selectedPages, setSelectedPages] = useState<string>('');
  const [selectedDataset, setSelectedDataset] = useState<string>('');
  const [selectedPageIndex, setselectedPageIndex] = useState<number>(0);
  const [cachedImages, setCachedImages] = useState<CachedImage[]>([]);
  
  // Function to add an image to the cache
  const addImageToCache = useCallback((item_job_name: string, item_dataset_name: string, item_image_name: string, item_image_path: string) => {
    setCachedImages(prev => {
      const exists = prev.some(
        img => img.item_job_name === item_job_name && 
               img.item_dataset_name === item_dataset_name && 
               img.item_image_path === item_image_path && 
               img.item_image_name === item_image_name
      );

      if (!exists) {
        const newCachedImage: CachedImage = { item_job_name, item_dataset_name, item_image_name, item_image_path };
        return [newCachedImage, ...prev];
      }
      return prev;
    });
  }, []);

  // Function to remove an image from the cache
  const removeImageFromCache = useCallback((imagePath: string) => {
    setCachedImages(prev => prev.filter(img => img.item_image_path !== imagePath));
  }, []);

  const getCache = useCallback((job: string, dataset: string) => {
    return cachedImages
      .filter(img => img.item_job_name === job && img.item_dataset_name === dataset)
      .map(img => img.item_image_path);
  }, [cachedImages]);

  // Load pending review images on initial render
  useEffect(() => {
    async function loadPending() {
      try {
        const data = await getPendingReview(true);
        data.items.forEach((item: CachedImage) => {
          addImageToCache(item.item_job_name, item.item_dataset_name, item.item_image_name, item.item_image_path);
        });
      } catch (error) {
        throw new Error('Error loading pending review: ' + error);
      }
    }

    loadPending();
  }, [addImageToCache]);

  // Update job pages when selectedJob changes
  useEffect(() => {
    async function updatePageData() {
      if (!selectedJob) return;

      try {
        await updateALLPages(selectedJob, IMAGES_PER_PAGE);
        setSelectedPages(selectedJob);

      } catch (error) {
        throw new Error('Error fetching job pages: ' + error);
      }
    }

    updatePageData();
  }, [selectedJob]);

  return (
    <JobDatasetContext.Provider
      value={{
        selectedJob,
        selectedPages,
        selectedDataset,
        selectedPageIndex,
        setSelectedJob,
        setSelectedPages,
        setSelectedDataset,
        setselectedPageIndex,
        cachedImages,
        addImageToCache,
        removeImageFromCache,
        getCache,
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

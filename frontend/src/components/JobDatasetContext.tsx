"use client";

import React, { createContext, useState, useContext, ReactNode, useEffect, useCallback, useMemo } from 'react';
import { CachedImage, JobDatasetContextType } from '@/types/JobDatasetContext';
import { getPendingReview, updateALLPages } from '@/services/api';
import { IMAGES_PER_PAGE } from '@/services/config';
import { logger } from '@/utils/logger';

// ============================================================================
// Context
// ============================================================================

const JobDatasetContext = createContext<JobDatasetContextType | undefined>(undefined);

// ============================================================================
// Provider Component
// ============================================================================

/**
 * JobDatasetProvider component provides context for managing job and dataset state.
 * 
 * Features:
 * - Manages job, dataset, and page selection state
 * - Caches pending review images
 * - Auto-fetches and updates pages when job changes
 * - Optimized with useMemo to prevent unnecessary re-renders
 * 
 * @param {ReactNode} children - Child components that will have access to this context
 */
export function JobDatasetProvider({ children }: { children: ReactNode }) {
  // ========== State ==========
  const [selectedJob, setSelectedJob] = useState<string>('');
  const [selectedPages, setSelectedPages] = useState<string>('');
  const [selectedDataset, setSelectedDataset] = useState<string>('');
  const [selectedPageIndex, setselectedPageIndex] = useState<number>(0);
  const [cachedImages, setCachedImages] = useState<CachedImage[]>([]);
  const [isLoadingPending, setIsLoadingPending] = useState<boolean>(false);

  // ========== Cached Image Management ==========

  /**
   * Add an image to the pending review cache
   * Prevents duplicates automatically
   */
  const addImageToCache = useCallback(
    (item_job_name: string, item_dataset_name: string, item_image_name: string, item_image_path: string) => {
      setCachedImages(prev => {
        // Check for duplicates
        const exists = prev.some(
          img =>
            img.item_job_name === item_job_name &&
            img.item_dataset_name === item_dataset_name &&
            img.item_image_path === item_image_path &&
            img.item_image_name === item_image_name
        );

        if (!exists) {
          const newCachedImage: CachedImage = {
            item_job_name,
            item_dataset_name,
            item_image_name,
            item_image_path,
          };
          return [newCachedImage, ...prev];
        }
        return prev;
      });
    },
    []
  );

  /**
   * Remove an image from the cache by its path
   */
  const removeImageFromCache = useCallback((imagePath: string) => {
    setCachedImages(prev => prev.filter(img => img.item_image_path !== imagePath));
  }, []);

  /**
   * Clear all cached images
   */
  const clearCache = useCallback(() => {
    setCachedImages([]);
  }, []);

  /**
   * Get all cached image paths for a specific job and dataset
   * Note: Kept for backward compatibility, but prefer using cachedImages directly
   * to avoid dependency issues in useEffect
   */
  const getCache = useCallback(
    (job: string, dataset: string): string[] => {
      return cachedImages
        .filter(img => img.item_job_name === job && img.item_dataset_name === dataset)
        .map(img => img.item_image_path);
    },
    [cachedImages]
  );

  /**
   * Get count of cached images for a specific job
   */
  const getCacheCountForJob = useCallback(
    (job: string): number => {
      return cachedImages.filter(img => img.item_job_name === job).length;
    },
    [cachedImages]
  );

  // ========== Effects ==========

  /**
   * Load pending review images on initial mount
   */
  useEffect(() => {
    let isMounted = true;

    async function loadPending() {
      if (isLoadingPending) return; // Prevent duplicate calls

      try {
        setIsLoadingPending(true);
        const data = await getPendingReview(true);
        
        if (!isMounted) return;

        // Batch add all pending images
        if (data.items && Array.isArray(data.items)) {
          data.items.forEach((item: CachedImage) => {
            addImageToCache(
              item.item_job_name,
              item.item_dataset_name,
              item.item_image_name,
              item.item_image_path
            );
          });
        }
      } catch (error) {
        if (isMounted) {
          logger.error('Error loading pending review:', error);
        }
      } finally {
        if (isMounted) {
          setIsLoadingPending(false);
        }
      }
    }

    loadPending();

    return () => {
      isMounted = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Only run once on mount, addImageToCache is stable via useCallback

  /**
   * Update job pages when selectedJob changes
   */
  useEffect(() => {
    let isMounted = true;

    async function updatePageData() {
      if (!selectedJob) return;

      setSelectedDataset('');
      setselectedPageIndex(0);

      try {
        await updateALLPages(selectedJob, IMAGES_PER_PAGE);
        if (isMounted) {
          setSelectedPages(selectedJob);
        }
      } catch (error) {
        if (isMounted) {
          logger.error('[JobDatasetContext] Error updating pages for job:', selectedJob, error);
        }
      }
    }

    updatePageData();

    return () => {
      isMounted = false;
    };
  }, [selectedJob]);

  // ========== Memoized Context Value ==========
  
  /**
   * Memoize context value to prevent unnecessary re-renders
   * Only recreates when state or callbacks change
   */
  const contextValue = useMemo<JobDatasetContextType>(
    () => ({
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
      clearCache,
      getCache,
      getCacheCountForJob,
    }),
    [
      selectedJob,
      selectedPages,
      selectedDataset,
      selectedPageIndex,
      cachedImages,
      addImageToCache,
      removeImageFromCache,
      clearCache,
      getCache,
      getCacheCountForJob,
    ]
  );

  return (
    <JobDatasetContext.Provider value={contextValue}>
      {children}
    </JobDatasetContext.Provider>
  );
}

// ============================================================================
// Hook: useJobDataset
// ============================================================================

/**
 * Custom hook to access JobDataset context
 * 
 * @throws {Error} If used outside of JobDatasetProvider
 * @returns {JobDatasetContextType} Context value with state and actions
 * 
 * @example
 * const { selectedJob, setSelectedJob, cachedImages } = useJobDataset();
 */
export function useJobDataset(): JobDatasetContextType {
  const context = useContext(JobDatasetContext);
  if (!context) {
    throw new Error('useJobDataset must be used within a JobDatasetProvider');
  }
  return context;
}

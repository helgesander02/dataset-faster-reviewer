"use client";

import { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';
import { SidebarState, SidebarActions } from '@/types/HomeLeftSidebar';
import { fetchJobs, fetchJobMetadata } from '@/services/api';
import { DATASET_PER_PAGE, JOB_REFRESH_INTERVAL } from '@/services/config';
import { logger } from '@/utils/logger';


export function useLeftSidebar(): SidebarState & SidebarActions {
  // ========== Context ==========
  const {
    selectedPages,
    selectedPageIndex,
    setSelectedDataset,
    setselectedPageIndex,
  } = useJobDataset();

  // ========== State ==========
  const [currentJobList, setJobList] = useState<string[]>([]);
  const [currentDatasetList, setDatasetList] = useState<string[]>([]);
  const [currentPagenation, setCurrentPagenation] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // ========== Refs ==========
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const mountedRef = useRef<boolean>(true);

  // ========== Memoized Values ==========
  /**
   * Calculate total number of pages for dataset pagination
   */
  const totalDatasetPages = useMemo(
    () => Math.ceil(currentDatasetList.length / DATASET_PER_PAGE),
    [currentDatasetList.length]
  );

  /**
   * Check if previous page is available
   */
  const hasPreviousPage = useMemo(
    () => currentPagenation > 0,
    [currentPagenation]
  );

  /**
   * Check if next page is available
   */
  const hasNextPage = useMemo(
    () => currentPagenation < totalDatasetPages - 1,
    [currentPagenation, totalDatasetPages]
  );

  /**
   * Get current page datasets
   */
  const currentPageDatasets = useMemo(() => {
    const startIndex = currentPagenation * DATASET_PER_PAGE;
    const endIndex = startIndex + DATASET_PER_PAGE;
    return currentDatasetList.slice(startIndex, endIndex);
  }, [currentDatasetList, currentPagenation]);

  // ========== Load Jobs ==========
  const loadJobs = useCallback(async () => {
    if (!mountedRef.current) return;

    try {
      setLoading(true);
      setError(null);

      const response = await fetchJobs();

      if (!mountedRef.current) return;

      if (!response?.job_names) {
        throw new Error('Invalid response format');
      }

      setJobList(response.job_names);
    } catch (err) {
      if (mountedRef.current) {
        const errorMessage = err instanceof Error ? err.message : 'Unable to load jobs';
        setError(errorMessage);
        logger.error('Error loading jobs:', err);
      }
    } finally {
      if (mountedRef.current) {
        setLoading(false);
      }
    }
  }, []);

  // ========== Load Datasets ==========
  const loadDatasets = useCallback(async () => {
    if (!selectedPages) {
      logger.log('[useLeftSidebar] loadDatasets: No selectedPages, skipping');
      return;
    }
    
    if (!mountedRef.current) {
      logger.log('[useLeftSidebar] loadDatasets: Component not mounted, skipping');
      return;
    }

    logger.log('[useLeftSidebar] loadDatasets: Starting for job:', selectedPages);

    try {
      setLoading(true);
      setError(null);

      logger.log('[useLeftSidebar] Calling fetchJobMetadata for:', selectedPages);
      const response = await fetchJobMetadata(selectedPages);

      if (!mountedRef.current) {
        logger.log('[useLeftSidebar] Component unmounted after fetchJobMetadata, returning');
        return;
      }

      logger.log('[useLeftSidebar] fetchJobMetadata response:', response);

      if (!response?.dataset_names) {
        throw new Error('Invalid response format');
      }

      logger.log('[useLeftSidebar] Setting dataset list with', response.dataset_names.length, 'datasets:', response.dataset_names);
      setDatasetList(response.dataset_names);
      setCurrentPagenation(0);

      if (response.dataset_names.length > 0) {
        logger.log('[useLeftSidebar] Auto-selecting first dataset:', response.dataset_names[0]);
        setselectedPageIndex(0);
        setSelectedDataset(response.dataset_names[0] || '');
      } else {
        logger.log('[useLeftSidebar] No datasets available to select');
      }
    } catch (err) {
      if (mountedRef.current) {
        const errorMessage = err instanceof Error ? err.message : 'Unable to load datasets';
        setError(errorMessage);
        logger.error('[useLeftSidebar] Error loading datasets:', err);
      }
    } finally {
      if (mountedRef.current) {
        logger.log('[useLeftSidebar] loadDatasets completed for:', selectedPages);
        setLoading(false);
      }
    }
  }, [selectedPages, setSelectedDataset, setselectedPageIndex]);

  // ========== Pagination Actions ==========
  /**
   * Navigate to previous page
   */
  const previousPage = useCallback(() => {
    if (!hasPreviousPage) return;

    const newPage = currentPagenation - 1;
    setCurrentPagenation(newPage);

    const newIndex = newPage * DATASET_PER_PAGE;
    if (newIndex < currentDatasetList.length) {
      setselectedPageIndex(newIndex);
      setSelectedDataset(currentDatasetList[newIndex] || '');
    }
  }, [hasPreviousPage, currentPagenation, currentDatasetList, setselectedPageIndex, setSelectedDataset]);

  /**
   * Navigate to next page
   */
  const nextPage = useCallback(() => {
    if (!hasNextPage) return;

    const newPage = currentPagenation + 1;
    setCurrentPagenation(newPage);

    const newIndex = newPage * DATASET_PER_PAGE;
    if (newIndex < currentDatasetList.length) {
      setselectedPageIndex(newIndex);
      setSelectedDataset(currentDatasetList[newIndex] || '');
    }
  }, [hasNextPage, currentPagenation, currentDatasetList, setselectedPageIndex, setSelectedDataset]);

  /**
   * Jump to specific page
   */
  const goToPage = useCallback(
    (page: number) => {
      if (page < 0 || page >= totalDatasetPages) return;

      setCurrentPagenation(page);

      const newIndex = page * DATASET_PER_PAGE;
      if (newIndex < currentDatasetList.length) {
        setselectedPageIndex(newIndex);
        setSelectedDataset(currentDatasetList[newIndex] || '');
      }
    },
    [totalDatasetPages, currentDatasetList, setselectedPageIndex, setSelectedDataset]
  );

  // ========== Effects ==========
  /**
   * Load jobs on mount and setup auto-refresh
   */
  useEffect(() => {
    mountedRef.current = true;

    loadJobs();

    intervalRef.current = setInterval(() => {
      loadJobs();
    }, JOB_REFRESH_INTERVAL);

    return () => {
      mountedRef.current = false;
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    };
  }, [loadJobs]);

  /**
   * Load datasets when selected job changes
   */
  useEffect(() => {
    if (selectedPages) {
      loadDatasets();
    }
  }, [selectedPages, loadDatasets]);

  /**
   * Auto-paginate when selectedPageIndex exceeds current page capacity
   */
  useEffect(() => {
    if (selectedPageIndex < 0) return;

    const requiredPage = Math.floor(selectedPageIndex / DATASET_PER_PAGE);

    if (requiredPage !== currentPagenation && requiredPage < totalDatasetPages) {
      setCurrentPagenation(requiredPage);
    }
  }, [selectedPageIndex, currentPagenation, totalDatasetPages]);

  // ========== Return ==========

  return {
    // State
    currentJobList,
    currentDatasetList,
    currentPagenation,
    loading,
    error,

    // Computed
    totalDatasetPages,
    hasPreviousPage,
    hasNextPage,
    currentPageDatasets,

    // Actions
    previousPage,
    nextPage,
    goToPage,
    loadJobs,
    loadDatasets,
  };
}

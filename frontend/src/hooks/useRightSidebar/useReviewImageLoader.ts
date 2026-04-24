"use client";

import { useState, useCallback, useRef } from 'react';
import { getPendingReview } from '@/services/api';
import { ReviewItem } from '@/types/HomeReview';
import { logger } from '@/utils/logger';
import { API_BASE_URL } from '@/services/config';

export interface ReviewImageWithUrl extends ReviewItem {
  displayUrl?: string;
}

/**
 * Hook for progressive loading of review images.
 * Loads metadata initially, then loads images page by page (9 items per page).
 */
export function useReviewImageLoader() {
  const [allItems, setAllItems] = useState<ReviewItem[]>([]);
  const [loadedItems, setLoadedItems] = useState<ReviewImageWithUrl[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(0);
  const [totalItems, setTotalItems] = useState(0);
  const [hasMore, setHasMore] = useState(false);
  
  const loadedPagesRef = useRef<Set<number>>(new Set());
  const currentPageRef = useRef(0);
  const allItemsRef = useRef<ReviewItem[]>([]);
  const itemsPerPage = 9;

  const getDisplayUrl = useCallback((item: ReviewItem): string => {
    if (item.item_image_path.startsWith('data:') || 
        item.item_image_path.startsWith('http')) {
      return item.item_image_path;
    }
    
    return `${API_BASE_URL}/api/getReviewImage?job=${encodeURIComponent(item.item_job_name)}&dataset=${encodeURIComponent(item.item_dataset_name)}&imageName=${encodeURIComponent(item.item_image_name)}`;
  }, []);

  /**
   * Load initial metadata (paths only, no images).
   * Fast operation even with thousands of items.
   */
  const loadMetadata = useCallback(async (): Promise<ReviewItem[]> => {
    try {
      setError(null);
      
      const data = await getPendingReview(true);
      
      if (!data || !Array.isArray(data.items)) {
        throw new Error('Invalid data format');
      }
      
      allItemsRef.current = data.items;
      setAllItems(data.items);
      setTotalItems(data.items.length);
      currentPageRef.current = 0;
      setCurrentPage(0);
      loadedPagesRef.current.clear();
      setLoadedItems([]);
      
      const maxPage = Math.ceil(data.items.length / itemsPerPage) - 1;
      setHasMore(0 < maxPage);
      
      return data.items;
    } catch (err) {
      logger.error('Failed to load review metadata:', err);
      setError('Failed to load review data');
      throw err;
    }
  }, []);

  const loadPage = useCallback(async (pageIndex: number, showLoading: boolean = false): Promise<void> => {
    if (loadedPagesRef.current.has(pageIndex)) {
      return;
    }

    const items = allItemsRef.current;
    const startIndex = pageIndex * itemsPerPage;
    const endIndex = Math.min(startIndex + itemsPerPage, items.length);
    
    if (startIndex >= items.length) {
      return;
    }

    try {
      if (showLoading) {
        setLoading(true);
      }
      
      const pageItems = items.slice(startIndex, endIndex);
      const itemsWithUrls: ReviewImageWithUrl[] = pageItems.map(item => ({
        ...item,
        displayUrl: getDisplayUrl(item)
      }));

      setLoadedItems(prev => [...prev, ...itemsWithUrls]);
      loadedPagesRef.current.add(pageIndex);
    } finally {
      if (showLoading) {
        setLoading(false);
      }
    }
  }, [itemsPerPage, getDisplayUrl]);

  const loadNextPage = useCallback(async (): Promise<void> => {
    const nextPage = currentPageRef.current + 1;
    const items = allItemsRef.current;
    const maxPage = Math.ceil(items.length / itemsPerPage) - 1;
    
    if (nextPage > maxPage) {
      setHasMore(false);
      return;
    }

    await loadPage(nextPage);
    currentPageRef.current = nextPage;
    setCurrentPage(nextPage);
    setHasMore(nextPage < maxPage);
  }, [itemsPerPage, loadPage]);

  const loadInitialPage = useCallback(async (): Promise<void> => {
    const items = allItemsRef.current;
    if (items.length === 0) {
      return;
    }

    await loadPage(0, true);
    currentPageRef.current = 0;
    setCurrentPage(0);
  }, [loadPage]);

  const reset = useCallback(() => {
    allItemsRef.current = [];
    setAllItems([]);
    setLoadedItems([]);
    currentPageRef.current = 0;
    setCurrentPage(0);
    setTotalItems(0);
    setError(null);
    setHasMore(false);
    loadedPagesRef.current.clear();
  }, []);

  return {
    allItems,
    loadedItems,
    loading,
    error,
    totalItems,
    currentPage,
    hasMorePages: hasMore,
    loadMetadata,
    loadInitialPage,
    loadNextPage,
    loadPage,
    reset,
    getDisplayUrl,
  };
}

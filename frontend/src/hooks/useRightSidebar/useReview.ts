"use client";

import { useState, useEffect, useCallback } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';
import { savePendingReview, deleteSelectedImages } from '@/services/api';
import { PendingReviewData, ReviewItem } from '@/types/HomeReview';
import { logger } from '@/utils/logger';
import { useReviewImageLoader } from './useReviewImageLoader';

const ERROR_MESSAGES = {
  LOAD_FAILED: 'Unable to load the data to be reviewed, please try again later.',
  INVALID_FORMAT: 'Invalid data format',
  REFRESH_FAILED: 'Failed to refresh data after deletion',
  DELETE_FAILED: 'Failed to delete images. Please try again.',
  NO_SELECTION: 'Please select at least one image to delete.',
} as const;

export function useHomeReview(isOpen: boolean) {
  const { cachedImages, addImageToCache, removeImageFromCache } = useJobDataset();
  const imageLoader = useReviewImageLoader();
  const [reviewData, setReviewData] = useState<PendingReviewData | null>(null);
  const [selectedImages, setSelectedImages] = useState<Set<string>>(new Set());
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [selectedJob, setSelectedJob] = useState<string>('');
  const [selectedDataset, setSelectedDataset] = useState<string>('');

  // Load metadata when modal opens
  useEffect(() => {
    if (isOpen) {
      fetchPendingReview();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen]);

  // Sync selected images with cached images
  useEffect(() => {
    if (reviewData?.items) {
      const cachedImageName = new Set(cachedImages.map(img => img.item_image_name));
      setSelectedImages(cachedImageName);
    }
  }, [reviewData, cachedImages]);

  const fetchPendingReview = async () => {
    try {
      const items = await imageLoader.loadMetadata();
      setReviewData({ items });
      await imageLoader.loadInitialPage();
    } catch (error) {
      logger.error('Failed to fetch pending review:', error);
    }
  };

  const createSaveData = useCallback(() => ({
    job: selectedJob,
    dataset: selectedDataset,
    images: cachedImages.map(img => ({
      job: img.item_job_name,
      dataset: img.item_dataset_name,
      imageName: img.item_image_name,
      imagePath: img.item_image_path
    })),
    timestamp: new Date().toISOString()
  }), [selectedJob, selectedDataset, cachedImages]);

  const saveToPendingReview = useCallback(async () => {
    try {
      setSaving(true);
      await savePendingReview(createSaveData());
    } finally {
      setSaving(false);
    }
  }, [createSaveData]);

  const toggleImageSelection = async (item: ReviewItem) => {
    const { item_job_name, item_dataset_name, item_image_name, item_image_path } = item;
    
    setSelectedImages(prev => {
      const isCurrentlySelected = prev.has(item_image_name);
      
      if (isCurrentlySelected) {
        const newSet = new Set(prev);
        newSet.delete(item_image_name);
        removeImageFromCache(item_image_path);
        return newSet;
      } else {
        const newSet = new Set(prev).add(item_image_name);
        addImageToCache(item_job_name, item_dataset_name, item_image_name, item_image_path);
        return newSet;
      }
    });
  };

  const selectAllImages = async () => {
    if (!reviewData?.items) return;
    
    reviewData.items.forEach(item => {
      addImageToCache(item.item_job_name, item.item_dataset_name, item.item_image_name, item.item_image_path);
    });
  };

  const deselectAllImages = async () => {
    if (!reviewData?.items) return;
    
    reviewData.items.forEach(item => {
      removeImageFromCache(item.item_image_path);
    });
  };

  const confirmDeletion = (count: number): boolean => {
    const message = `Are you sure you want to delete ${count} selected image(s)? This action cannot be undone.`;
    return window.confirm(message);
  };

  const clearCachedImages = useCallback(() => {
    cachedImages.forEach(img => removeImageFromCache(img.item_image_path));
    setSelectedImages(new Set());
  }, [cachedImages, removeImageFromCache]);

  const refreshReviewData = useCallback(async () => {
    try {
      imageLoader.reset();
      const items = await imageLoader.loadMetadata();
      setReviewData({ items });
      await imageLoader.loadInitialPage();
    } catch (error) {
      logger.error('Failed to refresh review data:', error);
    }
  }, [imageLoader]);

  const deleteSelectedImagesHandler = useCallback(async () => {
    if (cachedImages.length === 0) {
      alert(ERROR_MESSAGES.NO_SELECTION);
      return;
    }

    if (!confirmDeletion(cachedImages.length)) {
      return;
    }

    try {
      setDeleting(true);
      
      const imagesToDelete = cachedImages.map(img => ({
        job: img.item_job_name,
        dataset: img.item_dataset_name,
        imageName: img.item_image_name,
        imagePath: img.item_image_path
      }));

      const result = await deleteSelectedImages(imagesToDelete);
      
      clearCachedImages();
      await refreshReviewData();
      
      window.location.reload();
      
      alert(`Successfully deleted ${result.deleted_count} image(s)! Page will reload to refresh cache.`);
    } catch (err) {
      logger.error('Failed to delete images:', err);
      const errorMessage = err instanceof Error ? err.message : ERROR_MESSAGES.DELETE_FAILED;
      alert(`${ERROR_MESSAGES.DELETE_FAILED} Error: ${errorMessage}`);
    } finally {
      setDeleting(false);
    }
  }, [cachedImages, clearCachedImages, refreshReviewData]);

  return {
    reviewData,
    loading: imageLoader.loading,
    error: imageLoader.error,
    selectedImages,
    saving,
    deleting,
    selectedJob,
    selectedDataset,
    setSelectedJob,
    setSelectedDataset,
    fetchPendingReview,
    saveToPendingReview,
    toggleImageSelection,
    selectAllImages,
    deselectAllImages,
    deleteSelectedImagesHandler,
    loadedItems: imageLoader.loadedItems,
    loadNextPage: imageLoader.loadNextPage,
    hasMorePages: imageLoader.hasMorePages,
    totalItems: imageLoader.totalItems,
  };
}

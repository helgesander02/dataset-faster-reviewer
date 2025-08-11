"use client";

import { useState, useEffect, useCallback } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';
import { getPendingReview, savePendingReview } from '@/services/api';
import { PendingReviewData, ReviewItem } from '@/types/HomeReview';

export function useHomeReview(isOpen: boolean) {
  const { cachedImages, addImageToCache, removeImageFromCache } = useJobDataset();
  const [reviewData, setReviewData] = useState<PendingReviewData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedImages, setSelectedImages] = useState<Set<string>>(new Set());
  const [saving, setSaving] = useState(false);
  const [selectedJob, setSelectedJob] = useState<string>('');
  const [selectedDataset, setSelectedDataset] = useState<string>('');

  // Fetch pending review data when modal opens
  useEffect(() => {
    if (isOpen) {
      fetchPendingReview();
    }
  }, [isOpen]);

  // Sync selected images with cached images
  useEffect(() => {
    if (reviewData?.items) {
      const cachedImageName = new Set(cachedImages.map(img => img.item_image_name));
      setSelectedImages(cachedImageName);
    }
  }, [reviewData, cachedImages]);

  // Fetch pending review data
  const fetchPendingReview = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await getPendingReview(true);
      if (!data || !Array.isArray(data.items)) {
        throw new Error('Invalid data format');
      }
      setReviewData(data);

    } catch (err) {
      console.error('Error fetching pending review:', err);
      setError('Unable to load the data to be reviewed, please try again later.');

    } finally {
      setLoading(false);
    }
  };

  // Save to pending review
  const saveToPendingReview = useCallback(async () => {
    try {
      setSaving(true);
      const saveData = {
        job: selectedJob,
        dataset: selectedDataset,
        images: cachedImages.map(img => ({
          job: img.item_job_name,
          dataset: img.item_dataset_name,
          imageName: img.item_image_name,
          imagePath: img.item_image_path
        })),
        timestamp: new Date().toISOString()
      };

      await savePendingReview(saveData);

    } catch (err) {
      console.error('Error saving pending review:', err);
      
    } finally {
      setSaving(false);
    }
  }, [selectedJob, selectedDataset, cachedImages]);

  // Toggle individual image selection
  const toggleImageSelection = async (item: ReviewItem) => {
    const { item_job_name, item_dataset_name, item_image_name, item_image_path } = item;
    const isCurrentlySelected = selectedImages.has(item_image_name);
    
    if (isCurrentlySelected) {
      removeImageFromCache(item_image_path);
    } else {
      addImageToCache(item_job_name, item_dataset_name, item_image_name, item_image_path);
    }
  };

  // Select all images
  const selectAllImages = async () => {
    if (!reviewData?.items) return;
    
    reviewData.items.forEach(item => {
      addImageToCache(item.item_job_name, item.item_dataset_name, item.item_image_name, item.item_image_path);
    });
  };

  // Deselect all images
  const deselectAllImages = async () => {
    if (!reviewData?.items) return;
    
    reviewData.items.forEach(item => {
      removeImageFromCache(item.item_image_path);
    });
  };

  return {
    reviewData,
    loading,
    error,
    selectedImages,
    saving,
    selectedJob,
    selectedDataset,
    setSelectedJob,
    setSelectedDataset,
    fetchPendingReview,
    saveToPendingReview,
    toggleImageSelection,
    selectAllImages,
    deselectAllImages
  };
}

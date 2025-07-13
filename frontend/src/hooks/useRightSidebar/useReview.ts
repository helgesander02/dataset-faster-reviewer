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
      const cachedImageName = new Set(cachedImages.map(img => img.imageName));
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
          job: img.job,
          dataset: img.dataset,
          imageName: img.imageName,
          imagePath: img.imagePath
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
    const { job, dataset, imageName, imagePath } = item;
    const isCurrentlySelected = selectedImages.has(imageName);
    
    if (isCurrentlySelected) {
      removeImageFromCache(imagePath);
    } else {
      addImageToCache(job, dataset, imageName, imagePath);
    }
  };

  // Select all images
  const selectAllImages = async () => {
    if (!reviewData?.items) return;
    
    reviewData.items.forEach(item => {
      addImageToCache(item.job, item.dataset, item.imageName, item.imagePath);
    });
  };

  // Deselect all images
  const deselectAllImages = async () => {
    if (!reviewData?.items) return;
    
    reviewData.items.forEach(item => {
      removeImageFromCache(item.imagePath);
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

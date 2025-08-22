"use client";

import { useState } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';
import { savePendingReview } from '@/services/api';
import { RightSidebarState, RightSidebarActions, SaveData, CachedImage } from '@/types/HomeRightSidebar';

export function useRightSidebar(): RightSidebarState & RightSidebarActions {
  const { selectedPages, selectedDataset, cachedImages } = useJobDataset();
  const [isReviewOpen, setIsReviewOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);

  // Function to handle saving pending review
  const handleSave = async () => {
    if (!selectedPages || !selectedDataset || cachedImages.length === 0) {
      console.log('No job, dataset, or images selected for saving.');
      alert('Please select a job, dataset, and ensure there are images to save.');
      return;
    }

    try {
      setLoading(true);
      setSaveSuccess(false);
      
      const saveData: SaveData = {
        job: selectedPages,
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
      setSaveSuccess(true);
      
      setTimeout(() => {
        setSaveSuccess(false);
      }, 3000);
      
    } catch (error) {
      console.error('Error saving pending review:', error);
      alert('Save failed! Please try again later');

    } finally {
      setLoading(false);
    }
  };

  // Function to handle opening the review modal
  const handleReview = () => {
    setIsReviewOpen(true);
  };

  const handleCloseReview = () => {
    setIsReviewOpen(false);
  };

  // Group cached images by job and dataset
  const groupedImages = cachedImages.reduce((acc, img) => {
    const key = `${img.item_job_name}`;
    if (!acc[key]) {
      acc[key] = [];
    }
    acc[key].push(img);
    return acc;
  }, {} as Record<string, CachedImage[]>);

  return {
    isReviewOpen,
    loading,
    saveSuccess,
    groupedImages,
    handleSave,
    handleReview,
    handleCloseReview
  };
}

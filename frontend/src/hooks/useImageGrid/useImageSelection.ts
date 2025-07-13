"use client";

import { useState, useCallback, useEffect } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';

export function useImageSelection(
  selectedJob: string | null, allDatasets: string[]
) {

  const [selectedImages, setSelectedImages] = useState<Set<string>>(new Set());
  const { getCache, addImageToCache, removeImageFromCache } = useJobDataset();

  // This effect initializes the selected images based on the job and datasets
  useEffect(() => {
    if (!selectedJob || allDatasets.length === 0) {
      setSelectedImages(new Set());
      return;
    }

    const allCachedImages = new Set<string>();
    for (const dataset of allDatasets) {
      const cachedImagePaths = getCache(selectedJob, dataset);
      cachedImagePaths.forEach(path => allCachedImages.add(path));
    }
    setSelectedImages(allCachedImages);
  }, [selectedJob, allDatasets, getCache]);

  // This function handles the click event on an image
  const handleImageClick = useCallback((imageName: string, imageUrl: string, dataset: string) => {
    if (!selectedJob) return;
    
    if (selectedImages.has(imageUrl)) {
      setSelectedImages(prev => {
        const newSet = new Set(prev);
        newSet.delete(imageUrl);
        return newSet;
      });
      removeImageFromCache(imageUrl);
    } else {
      setSelectedImages(prev => new Set(prev).add(imageUrl));
      addImageToCache(selectedJob, dataset, imageName, imageUrl);
    }
  }, [selectedJob, selectedImages, addImageToCache, removeImageFromCache]);

  return { selectedImages, handleImageClick };
}

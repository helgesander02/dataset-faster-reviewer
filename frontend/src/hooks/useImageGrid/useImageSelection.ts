"use client";

import { useState, useCallback, useEffect } from 'react';
import { useJobDataset } from '@/components/JobDatasetContext';

export function useImageSelection(
  selectedJob: string, 
  allDatasets: string[]
) {
  const [selectedImages, setSelectedImages] = useState<Set<string>>(new Set());
  const { cachedImages, addImageToCache, removeImageFromCache } = useJobDataset();

  useEffect(() => {
    if (!selectedJob || allDatasets.length === 0) {
      setSelectedImages(new Set());
      return;
    }

    const allCachedImages = new Set<string>();
    for (const dataset of allDatasets) {
      cachedImages
        .filter(img => img.item_job_name === selectedJob && img.item_dataset_name === dataset)
        .forEach(img => allCachedImages.add(img.item_image_path));
    }
    
    setSelectedImages(new Set(allCachedImages));
  }, [selectedJob, allDatasets, cachedImages]);

  const handleImageClick = useCallback((
    imageName: string, 
    imageUrl: string, 
    dataset: string, 
    imagePath?: string, 
    imageJob?: string
  ) => {
    const jobToUse = imageJob || selectedJob;
    if (!jobToUse) return;
    
    const cacheKey = imagePath || imageUrl;
    
    setSelectedImages(prev => {
      if (prev.has(cacheKey)) {
        const newSet = new Set(prev);
        newSet.delete(cacheKey);
        removeImageFromCache(cacheKey);
        return newSet;
      } else {
        const newSet = new Set(prev).add(cacheKey);
        addImageToCache(jobToUse, dataset, imageName, cacheKey);
        return newSet;
      }
    });
  }, [selectedJob, addImageToCache, removeImageFromCache]);

  return { selectedImages, handleImageClick };
}

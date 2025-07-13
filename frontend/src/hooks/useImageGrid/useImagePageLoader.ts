"use client";

import { useCallback, useRef } from 'react';
import { fetchBase64Images } from '@/services/api';
import { IMAGES_PER_PAGE } from '@/services/config';
import { Image, ImagePage } from '@/types/HomeImageGrid';

export function useImagePageLoader() {
  const loadedPagesRef = useRef<Set<string>>(new Set());

  const loadPageImages = useCallback(async (
    job: string, dataset: string, pageIndex: number
  ): Promise<{ imagePage: ImagePage; maxPage: number } | null> => {

    const pageKey = `${dataset}-${pageIndex}`;
    
    if (loadedPagesRef.current.has(pageKey)) {
      console.log(`Page ${pageKey} already loaded, skipping...`);
      return null;
    }

    try {
      const response = await fetchBase64Images(job, dataset, pageIndex, IMAGES_PER_PAGE);
      console.log(`Loaded page ${pageIndex} for dataset ${dataset}:`, response);
      
      if (response && response.images.base64_image_set && response.images.base64_image_set.length > 0) {
        const pageImages: Image[] = response.images.base64_image_set.map((image: { name: string, path: string }) => ({
          name: image.name,
          url: `data:image/webp;base64,${image.path}`,
          dataset: dataset
        }));

        const imagePage: ImagePage = {
          dataset: dataset, images: pageImages,isNewDataset: pageIndex === 0
        };

        loadedPagesRef.current.add(pageKey);
        return { imagePage, maxPage: response.maxPage };
      }
      return null;
      
    } catch (error) {
      console.error(`Error loading page ${pageIndex} for dataset ${dataset}:`, error);
      return null;
    }
  }, []);

  const clearLoadedPages = useCallback(() => {
    loadedPagesRef.current.clear();
  }, []);

  const isPageLoaded = useCallback((dataset: string, pageIndex: number) => {
    const pageKey = `${dataset}-${pageIndex}`;
    return loadedPagesRef.current.has(pageKey);
  }, []);

  return {
    loadPageImages,
    clearLoadedPages,
    isPageLoaded
  };
}

"use client";

import { useCallback, useRef } from 'react';
import { fetchBase64Images, fetchImages, fetchALLPages } from '@/services/api';
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
      const response = await fetchBase64Images(job, pageIndex);
      const responseImage = await fetchImages(job, pageIndex);
      console.log(`Loaded page ${pageIndex} for dataset ${dataset}:`, response);
      const maxPage = await fetchALLPages(job)
      
      if (response && response.base64_image_set && response.base64_image_set.length > 0) {
        const pageImages: Image[] = response.base64_image_set.map((base64Image: string, index: number) => ({
            name: responseImage.image_name_set[index],
            url: `data:image/webp;base64,${base64Image}`,
            dataset: dataset
        }));
    

        const imagePage: ImagePage = {
          dataset: dataset, images: pageImages,isNewDataset: pageIndex === 0
        };

        loadedPagesRef.current.add(pageKey);
        return { imagePage, maxPage: maxPage.total_pages };
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

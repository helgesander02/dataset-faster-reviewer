"use client";

import { useCallback, useRef } from 'react';
import { fetchBase64Images, fetchImages } from '@/services/api';
import { Image, ImagePage } from '@/types/HomeImageGrid';
import axios from 'axios';
import { logger } from '@/utils/logger';

interface LoadPageResult {
  imagePage: ImagePage;
  isEmpty?: boolean;
}

export function useImagePageLoader() {
  const loadedPagesRef = useRef<Set<string>>(new Set());
  const loadingPagesRef = useRef<Set<string>>(new Set());
  const abortControllersRef = useRef<Map<string, AbortController>>(new Map());

  const isValidBase64 = (str: string): boolean => {
    if (!str || str.trim().length === 0) {
      return false;
    }
    
    const base64Regex = /^[A-Za-z0-9+/]*={0,2}$/;
    return base64Regex.test(str) && str.length > 10;
  };

  const loadPageImages = useCallback(async (
    job: string,
    dataset: string,
    pageIndex: number
  ): Promise<LoadPageResult | null> => {
    const pageKey = `${dataset}-${pageIndex}`;

    const existingController = abortControllersRef.current.get(pageKey);
    if (existingController) {
      existingController.abort();
      abortControllersRef.current.delete(pageKey);
    }

    if (loadedPagesRef.current.has(pageKey)) {
      return null;
    }

    if (loadingPagesRef.current.has(pageKey)) {
      return null;
    }

    loadingPagesRef.current.add(pageKey);

    const controller = new AbortController();
    abortControllersRef.current.set(pageKey, controller);

    try {
      const [base64Response, imageResponse] = await Promise.all([
        fetchBase64Images(job, pageIndex),
        fetchImages(job, pageIndex),
      ]);

      if (!base64Response?.base64_image) {
        return null;
      }

      if (base64Response.base64_image.length === 0) {
        loadedPagesRef.current.add(pageKey);
        return { 
          imagePage: {
            dataset: dataset,
            images: [],
            isNewDataset: pageIndex === 0,
          },
          isEmpty: true 
        };
      }

      const validBase64Images = base64Response.base64_image.every(
        (base64: string) => isValidBase64(base64)
      );

      if (!validBase64Images) {
        return null;
      }

      const pageImages: Image[] = base64Response.base64_image.map((base64Image: string, index: number) => ({
        name: imageResponse.image_name[index] || `image-${index}`,
        url: `data:image/webp;base64,${base64Image}`,
        dataset: dataset,
        path: imageResponse.image_path[index],
        job: job,
      }));

      const imagePage: ImagePage = {
        dataset,
        images: pageImages,
        isNewDataset: pageIndex === 0,
      };

      loadedPagesRef.current.add(pageKey);
      abortControllersRef.current.delete(pageKey);
      return { imagePage };
    } catch (error) {
      if (axios.isCancel(error) || (error as Error).name === 'AbortError') {
        return null;
      }
      logger.error(`Error loading page ${pageIndex} for dataset ${dataset}:`, error);
      return null;
    } finally {
      loadingPagesRef.current.delete(pageKey);
      abortControllersRef.current.delete(pageKey);
    }
  }, []);

  const clearLoadedPages = useCallback(() => {
    abortControllersRef.current.forEach((controller) => {
      controller.abort();
    });
    abortControllersRef.current.clear();
    loadedPagesRef.current.clear();
    loadingPagesRef.current.clear();
  }, []);

  const removePageFromCache = useCallback((dataset: string, pageIndex: number) => {
    const pageKey = `${dataset}-${pageIndex}`;
    
    const controller = abortControllersRef.current.get(pageKey);
    if (controller) {
      controller.abort();
      abortControllersRef.current.delete(pageKey);
    }
    
    loadedPagesRef.current.delete(pageKey);
    loadingPagesRef.current.delete(pageKey);
  }, []);

  const isPageLoaded = useCallback((dataset: string, pageIndex: number): boolean => {
    const pageKey = `${dataset}-${pageIndex}`;
    return loadedPagesRef.current.has(pageKey);
  }, []);

  const isPageLoading = useCallback((dataset: string, pageIndex: number): boolean => {
    const pageKey = `${dataset}-${pageIndex}`;
    return loadingPagesRef.current.has(pageKey);
  }, []);

  return {
    loadPageImages,
    clearLoadedPages,
    removePageFromCache,
    isPageLoaded,
    isPageLoading,
  };
}

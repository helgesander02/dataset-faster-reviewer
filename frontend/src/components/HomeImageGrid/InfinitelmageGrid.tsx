"use client";

import React, { useEffect } from 'react';
import { InfiniteImageGridProps } from '@/types/HomeImageGrid';
import { useDatasets } from '@/hooks/useImageGrid/useDatasets';
import { useInfiniteImages } from '@/hooks/useImageGrid/useInfiniteImages';
import { useImageSelection } from '@/hooks/useImageGrid/useImageSelection';
import { useIntersectionObserver } from '@/hooks/useImageGrid/useIntersectionObserver';
import HomeImageGrid from './ImageGrid';
import EmptyState from './EmptyState';
import LoadingState from './LoadingState';
import LoadingTrigger from './InfiniteLoadingTrigger';
import '@/styles/HomeImageGrid.css';

export default function InfiniteImageGrid({ 
  selectedPages, selectedDataset, selectedPageIndex,
  setSelectedDataset, setselectedPageIndex
}: InfiniteImageGridProps) {

  const { DatasetList, loading: datasetsLoading } = useDatasets(selectedPages); //get dataset list 
  const {
    loading: imagesLoading,
    getCurrentImagePages,
    hasMorePages,
    loadNextPage,
    resetImages,
    registerPageElement
  } = useInfiniteImages(selectedPages, selectedDataset, DatasetList, selectedPageIndex, setSelectedDataset, setselectedPageIndex);
  
  const { selectedImages, handleImageClick } = useImageSelection(selectedPages, DatasetList);
  const loadingRef = useIntersectionObserver(loadNextPage, imagesLoading);

  // This function registers the page element for the intersection observer
  const setPageRef = (pageIndex: number) => (element: HTMLDivElement | null) => {
    if (element) {
      registerPageElement(pageIndex, element);
    }
  };

  // 
  useEffect(() => {
    if (!selectedPages) {
      resetImages();
    }
  }, [selectedPages, resetImages]);

  if (!selectedPages) {
    return <EmptyState />;
  }
  
  if (datasetsLoading || (DatasetList.length === 0 && datasetsLoading)) {
    return <LoadingState message="Loading datasets..." />;
  }

  if (DatasetList.length === 0) {
    return (
      <EmptyState 
        title="No Datasets Found" 
        message="No datasets found for this job." 
      />
    );
  }

  //
  const currentImagePages = getCurrentImagePages();
  
  return (
    <div className="infinite-container">
      <div className="infinite-wrapper">
        {currentImagePages.length > 0 ? (
          <>
            {currentImagePages.map((page, pageIndex) => (
              <div
                key={`${page.dataset}-${pageIndex}`}
                ref={setPageRef(pageIndex)}
              >
                <HomeImageGrid
                  images={page.images}
                  selectedImages={selectedImages}
                  isLoading={false}
                  onImageClick={handleImageClick}
                />
              </div>
            ))}
            
            {hasMorePages() && (
              <div ref={loadingRef}>
                <LoadingTrigger isLoading={imagesLoading} />
              </div>
            )}
          </>
        ) : (
          imagesLoading ? (
            <LoadingState message="Loading images..." />
          ) : (
            <EmptyState 
              title="No Images Found" 
              message="No images found for this job." 
            />
          )
        )}
      </div>
    </div>
  );
}
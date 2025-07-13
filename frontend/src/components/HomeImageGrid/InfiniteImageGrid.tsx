"use client";

import React, { useEffect, useRef } from 'react';
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
  selectedJob, selectedDataset, currentPage,
  setSelectedDataset, setCurrentPage
}: InfiniteImageGridProps) {

  const { allDatasets, loading: datasetsLoading } = useDatasets(selectedJob); 
  const {
    loading: imagesLoading,
    getCurrentImagePages, hasMorePages, loadNextPage, resetImages,registerPageElement
  } = useInfiniteImages(selectedJob, selectedDataset, allDatasets, currentPage, setSelectedDataset, setCurrentPage);
  
  const { selectedImages, handleImageClick } = useImageSelection(selectedJob, allDatasets);
  const loadingRef = useIntersectionObserver(loadNextPage, imagesLoading);

  // This function registers the page element for the infinite scroll functionality
  const setPageRef = (pageIndex: number) => (element: HTMLDivElement | null) => {
    if (element) {
      registerPageElement(pageIndex, element);
    }
  };

  // Reset images when the selected job changes or becomes null
  useEffect(() => {
    if (!selectedJob) {
      resetImages();
    }
  }, [selectedJob, resetImages]);

  if (!selectedJob) {
    return <EmptyState />;
  } 
  if (datasetsLoading || (allDatasets.length === 0 && datasetsLoading)) {
    return <LoadingState message="Loading datasets..." />;
  }
  if (allDatasets.length === 0) {
    return (
      <EmptyState 
        title="No Datasets Found" 
        message="No datasets found for this job." 
      />
    );
  }
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
                  datasetName={page.dataset}
                  onImageClick={handleImageClick}
                />
              </div>
            ))}
            
            {hasMorePages() && imagesLoading && (
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
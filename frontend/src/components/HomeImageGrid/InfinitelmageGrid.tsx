"use client";

import React, { useEffect, useRef, useCallback, useMemo } from 'react';
import { FixedSizeList, ListChildComponentProps } from 'react-window';
import { InfiniteImageGridProps } from '@/types/HomeImageGrid';
import { useDatasets } from '@/hooks/useImageGrid/useDatasets';
import { useInfiniteImages } from '@/hooks/useImageGrid/useInfiniteImages';
import { useImageSelection } from '@/hooks/useImageGrid/useImageSelection';
import HomeImageGrid from './ImageGrid';
import EmptyState from './EmptyState';
import LoadingState from './LoadingState';
import '@/styles/HomeImageGrid.css';

const PAGE_HEIGHT = 800;

const PageItem: React.FC<ListChildComponentProps> = ({ index, style, data }) => {
  const { 
    getPageData, 
    selectedImages, 
    handleImageClick,
    loadPage,
    preloadAdjacentPages
  } = data;

  const pageData = useMemo(() => getPageData(index), [getPageData, index]);

  useEffect(() => {
    if (pageData && !pageData.isLoading && pageData.images.length === 0) {
      loadPage(index);
    }
    preloadAdjacentPages(index);
  }, [index, pageData, loadPage, preloadAdjacentPages]);

  if (pageData.isLoading) {
    return (
      <div style={style}>
        <div className="page-container">
          <LoadingState message={`Loading page ${index + 1}...`} />
        </div>
      </div>
    );
  }

  if (pageData.isEmpty) {
    return (
      <div style={style}>
        <div className="page-container">
          <div className="no-images">
            <p>No images found in this page</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={style}>
      <div className="page-container">
        <HomeImageGrid
          images={pageData.images}
          selectedImages={selectedImages}
          isLoading={pageData.isLoading}
          onImageClick={handleImageClick}
        />
      </div>
    </div>
  );
};

export default function InfiniteImageGrid({ 
  selectedPages, 
  selectedDataset, 
  selectedPageIndex,
  setSelectedDataset, 
  setselectedPageIndex
}: InfiniteImageGridProps) {
  
  const listRef = useRef<FixedSizeList>(null);
  const isScrollingRef = useRef(false);
  const scrollTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const { DatasetList, loading: datasetsLoading } = useDatasets(selectedPages);
  
  const {
    totalPages,
    loading: imagesLoading,
    getPageData,
    loadPage,
    preloadAdjacentPages,
    resetImages,
    hasPages
  } = useInfiniteImages(
    selectedPages, 
    selectedDataset, 
    DatasetList, 
    selectedPageIndex, 
    setSelectedDataset
  );
  
  const { selectedImages, handleImageClick } = useImageSelection(selectedPages, DatasetList);

  const scrollToPage = useCallback((pageIndex: number) => {
    if (listRef.current && pageIndex >= 0 && pageIndex < totalPages) {
      isScrollingRef.current = true;
      listRef.current.scrollToItem(pageIndex, 'start');
      
      if (scrollTimeoutRef.current) {
        clearTimeout(scrollTimeoutRef.current);
      }
      
      scrollTimeoutRef.current = setTimeout(() => {
        isScrollingRef.current = false;
      }, 500);
    }
  }, [totalPages]);

  const handleScroll = useCallback(({ scrollOffset }: { scrollOffset: number }) => {
    if (isScrollingRef.current) return;

    const currentIndex = Math.floor(scrollOffset / PAGE_HEIGHT);
    
    if (currentIndex !== selectedPageIndex && 
        currentIndex >= 0 && 
        currentIndex < totalPages) {
      setselectedPageIndex(currentIndex);
    }
  }, [selectedPageIndex, totalPages, setselectedPageIndex]);

  const goToNextPage = useCallback(() => {
    const nextIndex = Math.min(selectedPageIndex + 1, totalPages - 1);
    if (nextIndex !== selectedPageIndex) {
      setselectedPageIndex(nextIndex);
      scrollToPage(nextIndex);
    }
  }, [selectedPageIndex, totalPages, setselectedPageIndex, scrollToPage]);

  const goToPrevPage = useCallback(() => {
    const prevIndex = Math.max(selectedPageIndex - 1, 0);
    if (prevIndex !== selectedPageIndex) {
      setselectedPageIndex(prevIndex);
      scrollToPage(prevIndex);
    }
  }, [selectedPageIndex, setselectedPageIndex, scrollToPage]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'ArrowUp' || e.key === 'PageUp') {
        e.preventDefault();
        goToPrevPage();
      } else if (e.key === 'ArrowDown' || e.key === 'PageDown') {
        e.preventDefault();
        goToNextPage();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [goToPrevPage, goToNextPage]);

  useEffect(() => {
    if (selectedPageIndex >= 0 && selectedPageIndex < totalPages && !isScrollingRef.current) {
      scrollToPage(selectedPageIndex);
    }
  }, [selectedPageIndex, totalPages, scrollToPage]);

  useEffect(() => {
    if (!selectedPages) {
      resetImages();
    }
  }, [selectedPages, resetImages]);

  const itemData = useMemo(() => ({
    getPageData,
    selectedImages,
    handleImageClick,
    loadPage,
    preloadAdjacentPages
  }), [getPageData, selectedImages, handleImageClick, loadPage, preloadAdjacentPages]);

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

  if (!hasPages && !imagesLoading) {
    return (
      <EmptyState 
        title="No Pages Found" 
        message="No pages found for this dataset." 
      />
    );
  }

  return (
    <div className="main-container">
      <div className="infinite-container react-window-container">
        {totalPages > 0 ? (
          <FixedSizeList
            ref={listRef}
            height={typeof window !== 'undefined' ? window.innerHeight : 800}
            width="100%"
            itemCount={totalPages}
            itemSize={PAGE_HEIGHT}
            itemData={itemData}
            onScroll={handleScroll}
            className="react-window-list"
            overscanCount={2}
          >
            {PageItem}
          </FixedSizeList>
        ) : (
          <LoadingState message="Loading pages..." />
        )}
      </div>

      
    </div>
  );
}
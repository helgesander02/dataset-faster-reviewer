"use client";

import React from 'react';
import { HomeImageGridProps } from '@/types/HomeImageGrid';
import LoadingState from './LoadingState';
import ImageItem from './ImageItem';
import { memo } from 'react';

export default memo(function HomeImageGrid({ 
  images, 
  selectedImages, 
  isLoading,
  onImageClick, 
}: HomeImageGridProps) {
  
  if (isLoading && images.length === 0) {
    return <LoadingState />;
  }

  return (
    <div className="flex-1 flex flex-col overflow-hidden">
      <div 
        className="grid grid-cols-[repeat(auto-fill,minmax(clamp(180px,16vw,300px),1fr))] gap-3 flex-1 overflow-y-auto p-2"
        role="grid"
        aria-label="Image grid"
      >
        {images.map((image, index) => {
          // Check if image is selected using path (if available) or URL
          const cacheKey = image.path || image.url;
          const isSelected = selectedImages.has(cacheKey);
          return (
            <ImageItem
              key={`${image.dataset}-${image.name}-${index}`}
              image={image}
              index={index}
              isSelected={isSelected}
              allImages={images}
              onImageClick={onImageClick}
            />
          );
        })}
      </div>
      
      {!isLoading && images.length === 0 && (
        <div className="flex items-center justify-center p-12 text-center text-gray-500 text-base font-medium" role="alert">
          <p>No images found in this page</p>
        </div>
      )}
    </div>
  );
});

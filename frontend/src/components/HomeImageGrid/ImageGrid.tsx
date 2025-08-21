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
    <div className="image-grid-section">
      <div 
        className="image-grid"
        role="grid"
        aria-label="Image grid"
      >
        {images.map((image, index) => {
          const isSelected = selectedImages.has(image.url);
          return (
            <ImageItem
              key={`${image.dataset}-${image.name}-${index}`}
              image={image}
              index={index}
              isSelected={isSelected}
              onImageClick={onImageClick}
            />
          );
        })}
      </div>
      
      {!isLoading && images.length === 0 && (
        <div className="no-images" role="alert">
          <p>No images found in this page</p>
        </div>
      )}
    </div>
  );
});

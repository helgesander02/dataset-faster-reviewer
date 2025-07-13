"use client";

import React from 'react';
import { HomeImageGridProps } from '@/types/HomeImageGrid';
import LoadingState from './LoadingState';
import ImageItem from './ImageItem';

export default function HomeImageGrid({ 
  images, selectedImages, isLoading, datasetName,
  onImageClick,  
}: HomeImageGridProps) {

  if (isLoading) {
    return <LoadingState />;
  }

  return (
    <div className="image-grid-section">
      <div className="image-grid">
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
    </div>
  );
}

"use client";

import React from 'react';
import { ImageItemProps } from '@/types/HomeImageGrid';
import SelectionIndicator from './ImageSelectionIndicator';

export default function ImageItem({ 
  image, index, isSelected, 
  onImageClick 
}: ImageItemProps) {
    
  const handleImageError = (e: React.SyntheticEvent<HTMLImageElement>) => {
    (e.target as HTMLImageElement).src = "/api/placeholder/150/150";
  };

  const handleClick = () => {
    onImageClick(image.name, image.url, image.dataset);
  };

  return (
    <div 
      key={`${image.dataset}-${image.name}-${index}`} 
      className={`image-item ${isSelected ? 'selected' : ''}`}
      onClick={handleClick}
    >
      <img 
        src={image.url} 
        alt={`${image.dataset} - Image ${index}`} 
        className="grid-image"
        onError={handleImageError}
      />
      {isSelected && <SelectionIndicator />}
      <div className="image-name">
        {image.name}
      </div>
    </div>
  );
}

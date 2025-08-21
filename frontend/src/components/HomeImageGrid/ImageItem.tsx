"use client";

import React from 'react';
import Image from 'next/image';
import { ImageItemProps } from '@/types/HomeImageGrid';
import SelectionIndicator from './ImageSelectionIndicator';
import { memo } from 'react';
import { useCallback } from 'react';

export default memo(function ImageItem({ 
  image, index, isSelected, 
  onImageClick 
}: ImageItemProps) {
  const handleClick = useCallback(() => {
    onImageClick(image.name, image.url, image.dataset);
  }, [image.name, image.url, image.dataset, onImageClick]);

  return (
    <div 
      className={`image-item ${isSelected ? 'selected' : ''}`}
      onClick={handleClick}
    >
      <Image 
        src={image.url}
        alt={`${image.dataset} - Image ${index}`}
        width={150}
        height={150}
        className="grid-image"
        loading={index < 4 ? "eager" : "lazy"}
        sizes="(max-width: 768px) 100px, 150px"
      />
      {isSelected && <SelectionIndicator />}
      <div className="image-name">
        {image.name}
      </div>
    </div>
  );
});

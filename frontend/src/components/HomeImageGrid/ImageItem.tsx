"use client";

import React from 'react';
import Image from 'next/image';
import { ImageItemProps } from '@/types/HomeImageGrid';
import SelectionIndicator from './ImageSelectionIndicator';

export default function ImageItem({ 
  image, index, isSelected, 
  onImageClick 
}: ImageItemProps) {

  const handleClick = () => {
    onImageClick(image.name, image.url, image.dataset);
  };

  return (
    <div 
      key={`${image.dataset}-${image.name}-${index}`} 
      className={`image-item ${isSelected ? 'selected' : ''}`}
      onClick={handleClick}
    >
      <Image 
        src={image.url}
        alt={`${image.dataset} - Image ${index}`}
        width={150}
        height={150}
        className="grid-image"
      />
      {isSelected && <SelectionIndicator />}
      <div className="image-name">
        {image.name}
      </div>
    </div>
  );
}

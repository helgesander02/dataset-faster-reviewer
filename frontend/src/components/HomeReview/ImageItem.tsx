"use client";

import Image from 'next/image';
import { Check, ImageIcon } from 'lucide-react';
import { ImageItemProps } from '@/types/HomeReview';
import { useState } from 'react';

export function ImageItem({ 
  item, index, isSelected, 
  onToggle 
}: ImageItemProps) {
  const [imageError, setImageError] = useState(false);

  const handleImageError = () => {
    setImageError(true);
  };

  return (
    <div 
      className={`home-review-image-item ${isSelected ? 'selected' : ''}`}
      onClick={() => onToggle(item)}
    >
      {!imageError ? (
        <Image 
          src={item.item_image_path} 
          alt={`Image ${index + 1}`}
          width={150}
          height={150}
          className="home-review-image"
          onError={handleImageError}
        />
      ) : (
        <div className="home-review-image-fallback">
          <ImageIcon size={32} />
        </div>
      )}
      {isSelected && (
        <div className="home-review-selection-indicator">
          <Check size={16} />
        </div>
      )}
    </div>
  );
}

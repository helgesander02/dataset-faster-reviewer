"use client";

import { Check, ImageIcon } from 'lucide-react';
import { ImageItemProps } from '@/types/HomeReview';

export function ImageItem({ 
  item, index, isSelected, 
  onToggle 
}: ImageItemProps) {

  const handleImageError = (e: React.SyntheticEvent<HTMLImageElement>) => {
    e.currentTarget.style.display = 'none';
    const fallback = e.currentTarget.nextElementSibling as HTMLElement;
    if (fallback) {
      fallback.style.display = 'flex';
    }
  };

  return (
    <div 
      className={`home-review-image-item ${isSelected ? 'selected' : ''}`}
      onClick={() => onToggle(item)}
    >
      <img 
        src={item.imagePath} 
        alt={`Image ${index + 1}`}
        className="home-review-image"
        onError={handleImageError}
      />
      <div className="home-review-image-fallback">
        <ImageIcon size={32} />
      </div>
      {isSelected && (
        <div className="home-review-selection-indicator">
          <Check size={16} />
        </div>
      )}
    </div>
  );
}

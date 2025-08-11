"use client";

import { ImageItem } from './ImageItem';
import { ImagesGridProps } from '@/types/HomeReview';

export function ImagesGrid({ 
  items, selectedImages, 
  onToggleImage 
}: ImagesGridProps) {
  
  return (
    <div className="home-review-images-grid">
      {items.map((item, index) => {
        const isSelected = selectedImages.has(item.item_image_name);
        return (
          <ImageItem
            key={`${item.item_image_path}-${index}`}
            item={item}
            index={index}
            isSelected={isSelected}
            onToggle={onToggleImage}
          />
        );
      })}
    </div>
  );
}

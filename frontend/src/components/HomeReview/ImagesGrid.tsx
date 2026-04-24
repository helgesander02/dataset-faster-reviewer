"use client";

import { useEffect, useRef } from 'react';
import { ImageItem } from './ImageItem';
import { ImagesGridProps, ReviewItemWithUrl } from '@/types/HomeReview';
import { API_BASE_URL } from '@/services/config';

export function ImagesGrid({ 
  items, 
  selectedImages, 
  onToggleImage, 
  loadedItems, 
  onLoadMore, 
  hasMorePages
}: ImagesGridProps) {
  const loadMoreRef = useRef<HTMLDivElement>(null);
  
  const displayItems: ReviewItemWithUrl[] = loadedItems && loadedItems.length > 0 
    ? loadedItems.filter(item => item !== undefined && item !== null)
    : items.map(item => ({
        ...item,
        displayUrl: getDisplayUrl(item)
      }));

  function getDisplayUrl(item: ReviewItemWithUrl): string {
    if (item.item_image_path?.startsWith('data:') || 
        item.item_image_path?.startsWith('http')) {
      return item.item_image_path;
    }
    
    return `${API_BASE_URL}/api/getReviewImage?job=${encodeURIComponent(item.item_job_name)}&dataset=${encodeURIComponent(item.item_dataset_name)}&imageName=${encodeURIComponent(item.item_image_name)}`;
  }
  
  useEffect(() => {
    if (!hasMorePages || !onLoadMore || !loadMoreRef.current) return;
    
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          onLoadMore();
        }
      },
      { threshold: 0.1, rootMargin: '100px' }
    );
    
    const currentRef = loadMoreRef.current;
    observer.observe(currentRef);
    
    return () => {
      if (currentRef) {
        observer.unobserve(currentRef);
      }
      observer.disconnect();
    };
  }, [hasMorePages, onLoadMore]);
  
  return (
    <>
      <div className="grid grid-cols-[repeat(auto-fill,minmax(clamp(180px,16vw,300px),1fr))] gap-3 p-0">
        {displayItems.map((item, index) => {
          const isSelected = selectedImages.has(item.item_image_name);
          return (
            <ImageItem
              key={`${item.item_image_name}-${index}`}
              item={item}
              index={index}
              isSelected={isSelected}
              allItems={items}
              onToggle={onToggleImage}
            />
          );
        })}
      </div>
      
      {hasMorePages && (
        <div ref={loadMoreRef} className="h-20 flex items-center justify-center mt-6">
          <div className="text-gray-500 text-sm">Loading more...</div>
        </div>
      )}
    </>
  );
}

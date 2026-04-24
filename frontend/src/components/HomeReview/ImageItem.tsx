"use client";

import Image from 'next/image';
import { Check, ImageIcon } from 'lucide-react';
import { ImageItemProps, ReviewItemWithUrl } from '@/types/HomeReview';
import { useState, useCallback } from 'react';
import ImageLightbox from '@/components/ImageLightbox';
import { API_BASE_URL } from '@/services/config';

export function ImageItem({ 
  item, 
  index, 
  isSelected, 
  allItems,
  onToggle 
}: ImageItemProps) {
  const [imageError, setImageError] = useState(false);
  const [showLightbox, setShowLightbox] = useState(false);

  const itemWithUrl = item as ReviewItemWithUrl;
  const displayUrl = itemWithUrl.displayUrl || 
    `${API_BASE_URL}/api/getReviewImage?job=${encodeURIComponent(item.item_job_name)}&dataset=${encodeURIComponent(item.item_dataset_name)}&imageName=${encodeURIComponent(item.item_image_name)}`;

  const handleImageError = () => {
    setImageError(true);
  };

  const handleImageClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    setShowLightbox(true);
  }, []);

  const handleSelectClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onToggle(item);
  }, [item, onToggle]);

  return (
    <>
      <div 
        className={`relative rounded-lg overflow-hidden transition-all duration-200 bg-white border-2 ${
          isSelected 
            ? 'border-blue-500 shadow-lg shadow-blue-500/30' 
            : 'border-transparent shadow-md hover:shadow-lg'
        } hover:-translate-y-0.5`}
      >
        <div onClick={handleImageClick} className="relative cursor-pointer">
          {!imageError ? (
            <Image 
              src={displayUrl} 
              alt={`Image ${index + 1}`}
              width={150}
              height={150}
              className="w-full h-[200px] object-cover transition-transform duration-200 hover:scale-105"
              loading={index < 4 ? "eager" : "lazy"}
              sizes="(max-width: 768px) 100px, 150px"
              quality={85}
              placeholder="blur"
              blurDataURL="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTUwIiBoZWlnaHQ9IjE1MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTUwIiBoZWlnaHQ9IjE1MCIgZmlsbD0iI2YzZjRmNiIvPjwvc3ZnPg=="
              unoptimized
              onError={handleImageError}
            />
          ) : (
            <div className="flex items-center justify-center bg-gray-100 text-gray-400 h-[200px]">
              <ImageIcon size={32} />
            </div>
          )}
        </div>
        
        <button
          onClick={handleSelectClick}
          className={`absolute top-2 right-2 w-7 h-7 rounded-full flex items-center justify-center transition-all duration-200 shadow-md ${
            isSelected 
              ? 'bg-blue-500 scale-110' 
              : 'bg-white/90 hover:bg-white hover:scale-110'
          }`}
          aria-label={isSelected ? 'Deselect image' : 'Select image'}
        >
          {isSelected ? (
            <Check size={16} className="text-white" />
          ) : (
            <div className="w-4 h-4 border-2 border-gray-400 rounded-full" />
          )}
        </button>
      </div>

      <ImageLightbox
        isOpen={showLightbox}
        job={item.item_job_name}
        dataset={item.item_dataset_name}
        imagePath={item.item_image_path}
        imageName={item.item_image_name}
        allImages={allItems.map(i => ({
          name: i.item_image_name,
          url: i.item_image_path,
          path: i.item_image_path,
          dataset: i.item_dataset_name,
          job: i.item_job_name
        }))}
        onClose={() => setShowLightbox(false)}
        showAllImages={true}
      />
    </>
  );
}

"use client";

import React, { memo, useCallback, useState } from 'react';
import Image from 'next/image';
import { Check } from 'lucide-react';
import { ImageItemProps } from '@/types/HomeImageGrid';
import ImageLightbox from '@/components/ImageLightbox';
import { useJobDataset } from '@/components/JobDatasetContext';

export default memo(function ImageItem({ 
  image, 
  index, 
  isSelected, 
  allImages,
  onImageClick 
}: ImageItemProps) {
  const { selectedJob } = useJobDataset();
  const [showLightbox, setShowLightbox] = useState(false);
  
  const handleImageClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    setShowLightbox(true);
  }, []);
  
  const handleSelectClick = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onImageClick(image.name, image.url, image.dataset, image.path, image.job);
  }, [image.name, image.url, image.dataset, image.path, image.job, onImageClick]);

  return (
    <>
      <div 
        className={`relative rounded-lg overflow-hidden transition-all duration-200 cursor-pointer bg-white border-2 ${
          isSelected 
            ? 'border-blue-500 shadow-lg shadow-blue-500/30' 
            : 'border-transparent shadow-md hover:shadow-lg'
        } hover:-translate-y-0.5`}
      >
        <div onClick={handleImageClick} className="relative">
          <Image 
            src={image.url}
            alt={`${image.dataset} - Image ${index}`}
            width={150}
            height={150}
            className="w-full h-[200px] object-cover transition-transform duration-200 hover:scale-105"
            loading={index < 4 ? "eager" : "lazy"}
            sizes="(max-width: 768px) 100px, 150px"
            quality={85}
            placeholder="blur"
            blurDataURL="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTUwIiBoZWlnaHQ9IjE1MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTUwIiBoZWlnaHQ9IjE1MCIgZmlsbD0iI2YzZjRmNiIvPjwvc3ZnPg=="
            unoptimized={image.url.startsWith('data:')}
          />
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
        
        <div className="px-2 py-2 text-xs text-gray-600 bg-white text-center font-medium whitespace-nowrap overflow-hidden text-ellipsis">
          {image.name}
        </div>
      </div>
      
      <ImageLightbox
        isOpen={showLightbox}
        job={selectedJob}
        imagePath={image.path || image.url}
        imageName={image.name}
        dataset={image.dataset}
        allImages={allImages}
        onClose={() => setShowLightbox(false)}
      />
    </>
  );
});

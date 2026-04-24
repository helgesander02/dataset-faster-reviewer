"use client";

import { X, ZoomIn, ZoomOut, ChevronLeft, ChevronRight, Check } from 'lucide-react';
import { useEffect, useState, useCallback, useMemo } from 'react';
import { createPortal } from 'react-dom';
import { API_BASE_URL } from '@/services/config';
import { useJobDataset } from '@/components/JobDatasetContext';

interface ImageData {
  name: string;
  url: string;
  path?: string;
  dataset: string;
  job?: string;
}

interface ImageLightboxProps {
  isOpen: boolean;
  job: string;
  imagePath: string;
  imageName: string;
  dataset: string;
  allImages: ImageData[];
  onClose: () => void;
  showAllImages?: boolean;
}

export default function ImageLightbox({ 
  isOpen, job, imagePath, imageName, dataset, allImages, onClose,
  showAllImages = false
}: ImageLightboxProps) {
  const [originalImageUrl, setOriginalImageUrl] = useState<string>('');
  const [loading, setLoading] = useState(false);
  const [scale, setScale] = useState(1);
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  
  const { addImageToCache, removeImageFromCache, cachedImages } = useJobDataset();
  
  const datasetImages = useMemo(() => 
    showAllImages ? allImages : allImages.filter(img => img.dataset === dataset),
    [allImages, dataset, showAllImages]
  );
  
  const currentImage = useMemo(() => 
    datasetImages[currentImageIndex] || {
      name: imageName,
      url: imagePath,
      path: imagePath || undefined,
      dataset: dataset
    },
    [datasetImages, currentImageIndex, imageName, imagePath, dataset]
  );
  
  useEffect(() => {
    if (isOpen) {
      const index = datasetImages.findIndex(img => img.name === imageName);
      setCurrentImageIndex(index >= 0 ? index : 0);
      document.body.style.overflow = 'hidden';
      setScale(1);
    } else {
      document.body.style.overflow = 'unset';
      setOriginalImageUrl('');
      setScale(1);
      setCurrentImageIndex(0);
    }
    
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen, imageName, datasetImages]);
  
  useEffect(() => {
    if (isOpen && currentImage) {
      fetchOriginalImage();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isOpen, currentImageIndex]);
  
  const handlePrevious = useCallback(() => {
    if (currentImageIndex > 0) {
      setCurrentImageIndex(prev => prev - 1);
      setScale(1); // Reset zoom when changing image
    }
  }, [currentImageIndex]);
  
  const handleNext = useCallback(() => {
    if (currentImageIndex < datasetImages.length - 1) {
      setCurrentImageIndex(prev => prev + 1);
      setScale(1); // Reset zoom when changing image
    }
  }, [currentImageIndex, datasetImages.length]);
  
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      } else if (e.key === 'ArrowLeft') {
        handlePrevious();
      } else if (e.key === 'ArrowRight') {
        handleNext();
      }
    };
    
    if (isOpen) {
      window.addEventListener('keydown', handleKeyDown);
    }
    
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose, handlePrevious, handleNext]);
  
  const fetchOriginalImage = useCallback(async () => {
    const imageJob = currentImage.job || job;
    
    if (!imageJob || !currentImage.dataset) {
      setOriginalImageUrl(currentImage.url);
      return;
    }
    
    setLoading(true);
    
    try {
      const url = `${API_BASE_URL}/api/getReviewImage?job=${encodeURIComponent(imageJob)}&dataset=${encodeURIComponent(currentImage.dataset)}&imageName=${encodeURIComponent(currentImage.name)}`;
      setOriginalImageUrl(url);
    } catch {
      setOriginalImageUrl(currentImage.url);
    } finally {
      setLoading(false);
    }
  }, [job, currentImage]);
  
  const handleResetZoom = () => {
    setScale(1);
  };
  
  const isCurrentImageSelected = useMemo(() => {
    const imageJob = currentImage.job || job;
    const cacheKey = currentImage.path || currentImage.url;
    return cachedImages.some(cached => 
      cached.item_job_name === imageJob && 
      cached.item_image_path === cacheKey
    );
  }, [currentImage, cachedImages, job]);
  
  const handleToggleCurrentImage = useCallback(() => {
    const imageJob = currentImage.job || job;
    const cacheKey = currentImage.path || currentImage.url;
    
    if (isCurrentImageSelected) {
      removeImageFromCache(cacheKey);
    } else {
      addImageToCache(imageJob, currentImage.dataset, currentImage.name, cacheKey);
    }
  }, [currentImage, isCurrentImageSelected, job, addImageToCache, removeImageFromCache]);
  
  if (!isOpen) return null;
  
  const lightboxContent = (
    <div 
      className="fixed inset-0 z-[9999] bg-black/30 backdrop-blur-sm flex items-center justify-center overflow-auto"
      onClick={onClose}
    >
      <button
        onClick={onClose}
        className="fixed top-4 right-4 p-3 bg-red-500/80 hover:bg-red-600/90 rounded-lg backdrop-blur-sm transition-colors z-[10001]"
        title="Close (ESC)"
      >
        <X size={22} className="text-white" />
      </button>
      
      <div className="fixed top-4 left-4 px-4 py-2 bg-white/10 backdrop-blur-sm rounded-lg z-[10001] max-w-[50vw]">
        <p className="text-white text-sm font-medium truncate">{currentImage.name}</p>
        <p className="text-white/70 text-xs mt-1">
          {showAllImages ? (
            <>{currentImageIndex + 1} / {datasetImages.length} • {currentImage.dataset}</>
          ) : (
            <>{currentImageIndex + 1} / {datasetImages.length} in {dataset}</>
          )}
        </p>
      </div>
      
      {currentImageIndex > 0 && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            handlePrevious();
          }}
          className="fixed left-4 top-1/2 -translate-y-1/2 p-3 bg-white/20 hover:bg-white/30 rounded-lg backdrop-blur-sm transition-colors z-[10001]"
          title="Previous Image (←)"
        >
          <ChevronLeft size={32} className="text-white" />
        </button>
      )}
      
      {currentImageIndex < datasetImages.length - 1 && (
        <button
          onClick={(e) => {
            e.stopPropagation();
            handleNext();
          }}
          className="fixed right-4 top-1/2 -translate-y-1/2 p-3 bg-white/20 hover:bg-white/30 rounded-lg backdrop-blur-sm transition-colors z-[10001]"
          title="Next Image (→)"
        >
          <ChevronRight size={32} className="text-white" />
        </button>
      )}
      <div className="fixed bottom-4 left-1/2 transform -translate-x-1/2 flex items-center gap-3 px-6 py-3 bg-white/20 backdrop-blur-sm rounded-lg z-[10001]"
        onClick={(e) => e.stopPropagation()}
      >
        <ZoomOut size={18} className="text-white flex-shrink-0" />
        
        <input
          type="range"
          min="50"
          max="300"
          step="25"
          value={scale * 100}
          onChange={(e) => {
            e.stopPropagation();
            setScale(Number(e.target.value) / 100);
          }}
          className="w-48 h-2 bg-white/30 rounded-lg appearance-none cursor-pointer
            [&::-webkit-slider-thumb]:appearance-none
            [&::-webkit-slider-thumb]:w-4
            [&::-webkit-slider-thumb]:h-4
            [&::-webkit-slider-thumb]:rounded-full
            [&::-webkit-slider-thumb]:bg-white
            [&::-webkit-slider-thumb]:cursor-pointer
            [&::-webkit-slider-thumb]:hover:bg-blue-400
            [&::-webkit-slider-thumb]:transition-colors
            [&::-moz-range-thumb]:w-4
            [&::-moz-range-thumb]:h-4
            [&::-moz-range-thumb]:rounded-full
            [&::-moz-range-thumb]:bg-white
            [&::-moz-range-thumb]:border-0
            [&::-moz-range-thumb]:cursor-pointer
            [&::-moz-range-thumb]:hover:bg-blue-400
            [&::-moz-range-thumb]:transition-colors"
          title="Zoom Level"
        />
        
        <ZoomIn size={18} className="text-white flex-shrink-0" />
        
        <button
          onClick={(e) => {
            e.stopPropagation();
            handleResetZoom();
          }}
          className="ml-2 px-3 py-1 bg-white/30 hover:bg-white/40 rounded transition-colors"
          title="Reset Zoom (100%)"
        >
          <span className="text-white text-sm font-medium">{Math.round(scale * 100)}%</span>
        </button>
      </div>
      
      {loading && (
        <div className="fixed inset-0 flex items-center justify-center z-[10000]">
          <div className="px-6 py-3 bg-white/10 backdrop-blur-sm rounded-lg">
            <p className="text-white text-base font-medium">Loading original image...</p>
          </div>
        </div>
      )}
      
      {originalImageUrl && !loading && (
        <div 
          style={{
            width: '90vw',
            height: '90vh',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            position: 'relative'
          }}
          onClick={(e) => e.stopPropagation()}
        >
          <img
            src={originalImageUrl}
            alt={currentImage.name}
            style={{ 
              maxWidth: '100%',
              maxHeight: '100%',
              width: 'auto',
              height: 'auto',
              objectFit: 'contain',
              imageRendering: 'auto',
              transform: `scale(${scale})`,
              transformOrigin: 'center center',
              transition: 'transform 0.2s ease-in-out'
            }}
          />
          
          <button
            onClick={(e) => {
              e.stopPropagation();
              handleToggleCurrentImage();
            }}
            className={`absolute top-4 right-4 w-10 h-10 rounded-full flex items-center justify-center transition-all duration-200 shadow-lg ${
              isCurrentImageSelected 
                ? 'bg-blue-500 scale-110 hover:bg-blue-600' 
                : 'bg-white/90 hover:bg-white hover:scale-110'
            }`}
            title={isCurrentImageSelected ? 'Deselect this image' : 'Select this image'}
          >
            {isCurrentImageSelected ? (
              <Check size={20} className="text-white" />
            ) : (
              <div className="w-5 h-5 border-2 border-gray-400 rounded-full" />
            )}
          </button>
        </div>
      )}
    </div>
  );

  return createPortal(lightboxContent, document.body);
}

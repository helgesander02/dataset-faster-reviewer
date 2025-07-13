"use client";

import { LoadingState } from './ContentLoadingState';
import { ErrorState } from './ContentErrorState';
import { EmptyState } from './ContentEmptyState';
import { ImagesGrid } from './ImagesGrid';
import { ReviewContentProps } from '@/types/HomeReview';

export function ReviewContent({ 
  loading, error, reviewData, selectedImages, 
  onRetry, onToggleImage 
}: ReviewContentProps) {

  const hasImages = reviewData && reviewData.items.length > 0;

  return (
    <div className="home-review-content">
      {loading ? (
        <LoadingState />
      ) : error ? (
        <ErrorState error={error} onRetry={onRetry} />
      ) : !hasImages ? (
        <EmptyState />
      ) : (
        <ImagesGrid 
          items={reviewData.items}
          selectedImages={selectedImages}
          onToggleImage={onToggleImage}
        />
      )}
    </div>
  );
}

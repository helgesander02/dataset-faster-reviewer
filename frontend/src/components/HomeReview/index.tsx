"use client";

import { ReviewHeader } from './ReviewHeader';
import { ReviewContent } from './ReviewContent';
import { ReviewActions } from './ReviewActions';
import { HomeReviewProps } from '@/types/HomeReview';
import { useHomeReview } from '@/hooks/useRightSidebar/useReview';
import '@/styles/HomeReview.css';

export default function HomeReview({ 
  isOpen, onClose 
}: HomeReviewProps) {

  const {
    reviewData, loading, error, selectedImages, saving,
    fetchPendingReview, saveToPendingReview, toggleImageSelection,
    selectAllImages, deselectAllImages
  } = useHomeReview(isOpen);
  if (!isOpen) return null;

  const hasImages = reviewData && reviewData.items.length > 0;
  const selectedCount = selectedImages.size;
  const totalCount = reviewData?.items.length || 0;

  return (
    <div className="home-review-overlay">
      <div className="home-review-modal">
        <ReviewHeader 
          saving={saving}
          onClose={onClose}
        />

        <ReviewContent
          loading={loading}
          error={error}
          reviewData={reviewData}
          selectedImages={selectedImages}
          onRetry={fetchPendingReview}
          onToggleImage={toggleImageSelection}
        />

        {hasImages && (
          <ReviewActions
            selectedCount={selectedCount}
            totalCount={totalCount}
            saving={saving}
            onSelectAll={selectAllImages}
            onDeselectAll={deselectAllImages}
            onSave={saveToPendingReview}
          />
        )}
      </div>
    </div>
  );
}

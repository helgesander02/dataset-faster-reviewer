"use client";

import { ReviewHeader } from './ReviewHeader';
import { ReviewContent } from './ReviewContent';
import { ReviewActions } from './ReviewActions';
import { HomeReviewProps } from '@/types/HomeReview';
import { useHomeReview } from '@/hooks/useRightSidebar/useReview';

export default function HomeReview({ 
  isOpen, onClose 
}: HomeReviewProps) {

  const {
    reviewData, loading, error, selectedImages, saving, deleting,
    fetchPendingReview, saveToPendingReview, toggleImageSelection,
    selectAllImages, deselectAllImages, deleteSelectedImagesHandler,
    loadedItems, loadNextPage, hasMorePages, totalItems
  } = useHomeReview(isOpen);
  if (!isOpen) return null;

  const hasImages = reviewData && reviewData.items.length > 0;
  const selectedCount = selectedImages.size;
  const totalCount = reviewData?.items.length || 0;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-[1000] p-4">
      <div className="bg-white rounded-lg w-full max-w-[1000px] max-h-[90vh] flex flex-col shadow-2xl">
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
          loadedItems={loadedItems}
          onLoadMore={loadNextPage}
          hasMorePages={hasMorePages}
          totalItems={totalItems}
        />

        {hasImages && (
          <ReviewActions
            selectedCount={selectedCount}
            totalCount={totalCount}
            saving={saving}
            deleting={deleting}
            onSelectAll={selectAllImages}
            onDeselectAll={deselectAllImages}
            onSave={saveToPendingReview}
            onDelete={deleteSelectedImagesHandler}
          />
        )}
      </div>
    </div>
  );
}

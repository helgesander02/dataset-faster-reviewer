"use client";

import { LoadingState } from './ContentLoadingState';
import { ErrorState } from './ContentErrorState';
import { EmptyState } from './ContentEmptyState';
import { ImagesGrid } from './ImagesGrid';
import { ReviewContentProps } from '@/types/HomeReview';

export function ReviewContent({ 
  loading, error, reviewData, selectedImages, 
  onRetry, onToggleImage, loadedItems, onLoadMore, hasMorePages, totalItems
}: ReviewContentProps) {

  const hasImages = reviewData && reviewData.items.length > 0;
  const hasLoadedImages = loadedItems && loadedItems.length > 0;
  
  return (
    <div className="flex-1 overflow-y-auto p-6">
      {error ? (
        <ErrorState error={error} onRetry={onRetry} />
      ) : !hasImages ? (
        loading ? <LoadingState /> : <EmptyState />
      ) : (
        <>
          {!hasLoadedImages && loading ? (
            <LoadingState />
          ) : (
            <>
              <ImagesGrid 
                items={reviewData.items}
                selectedImages={selectedImages}
                onToggleImage={onToggleImage}
                loadedItems={loadedItems}
                onLoadMore={onLoadMore}
                hasMorePages={hasMorePages}
              />
              {totalItems && (
                <div className="mt-4 text-center text-sm text-gray-600">
                  Showing {loadedItems?.length || 0} of {totalItems} items
                </div>
              )}
            </>
          )}
        </>
      )}
    </div>
  );
}

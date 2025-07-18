"use client";

import { ReviewHeader } from './ReviewHeader';
import { ReviewContent } from './ReviewContent';
import { ReviewActions } from './ReviewActions';
import { HomeReviewProps } from '@/types/HomeReview';
import { useHomeReview } from '@/hooks/useRightSidebar/useReview';
import '@/styles/HomeReview.css';

/*
 * HomeReview component for the HomeRighSidebar page.
 * This component displays a modal for reviewing images
 * and allows users to select, deselect, and save images.
 * 
 * Props:
 * - isOpen:                boolean - Indicates if the review modal is open.
 * - onClose:               function - Callback to close the review modal.
 * - loading:               boolean - Indicates if the review data is being loaded. 
 * - reviewData:            object - Contains the data for the images to be reviewed.
 * - error:                 string - Contains any error message if the review data fails to load.
 * - selectedImages:        Set - Contains the currently selected images.
 * - saving:                boolean - Indicates if the save operation is in progress.
 * - fetchPendingReview:    function - Function to fetch the pending review data.
 * - saveToPendingReview:   function - Function to save the selected images.
 * - toggleImageSelection:  function - Function to toggle the selection of an image.
 * - selectAllImages:       function - Function to select all images.
 * - deselectAllImages:     function - Function to deselect all images. 
 */
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

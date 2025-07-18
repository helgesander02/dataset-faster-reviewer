"use client";

import { useRightSidebar } from '@/hooks/useRightSidebar/useRightSidebar';
import '@/styles/HomeRightSidebar.css';

import { useJobDataset } from '../JobDatasetContext';
import ReviewButton from './ReviewButton';
import FileChangeLog from './FileChangeLog';
import SaveButton from './SaveButton';
import HomeReview from '../HomeReview';


/** 
 * RightSidebar component for the Home page.
 * This component displays the right sidebar with review and save functionalities.
 * It utilizes context to manage the selected job, dataset, and cached images.
 * 
 * Props:
 * - isReviewOpen:      boolean - Indicates if the review modal is open.
 * - loading:           boolean - Indicates if the save or review operation is in progress.
 * - saveSuccess:       boolean - Indicates if the save operation was successful.
 * - groupedImages:     Record<string, CachedImage[]> - Grouped images for the file change log.
 * - handleSave:        function - Function to handle saving pending review.
 * - handleReview:      function - Function to handle opening the review modal.
 * - handleCloseReview: function - Function to close the review modal.
 */
export default function RightSidebar() {
  const { selectedJob, selectedDataset, cachedImages } = useJobDataset();
  const {
    isReviewOpen, loading, saveSuccess, groupedImages,
    handleSave, handleReview, handleCloseReview
  } = useRightSidebar();

  return (
    <div className="right-sidebar">
      <ReviewButton onReview={handleReview} loading={loading} />
      <FileChangeLog groupedImages={groupedImages} cachedImages={cachedImages} />
      <div className="mt-auto">
        <SaveButton 
          onSave={handleSave} 
          loading={loading} 
          saveSuccess={saveSuccess} 
          cachedImages={cachedImages} 
          disabled={!selectedJob || !selectedDataset || cachedImages.length === 0 || loading} 
        />
      </div>
      
      <HomeReview 
        isOpen={isReviewOpen} 
        onClose={handleCloseReview} 
      />
    </div>
  );
}

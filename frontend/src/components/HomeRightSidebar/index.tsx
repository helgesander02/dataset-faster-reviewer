"use client";

import { useRightSidebar } from '@/hooks/useRightSidebar/useRightSidebar';

import { useJobDataset } from '../JobDatasetContext';
import ReviewButton from './ReviewButton';
import FileChangeLog from './FileChangeLog';
import SaveButton from './SaveButton';
import HomeReview from '../HomeReview';


export default function RightSidebar() {
  const { selectedJob, selectedDataset, cachedImages } = useJobDataset();
  const {
    isReviewOpen, loading, saveSuccess, groupedImages,
    handleSave, handleReview, handleCloseReview
  } = useRightSidebar();

  return (
    <div className="w-[16.666667%] bg-gray-100 p-2 flex flex-col h-full">
      <ReviewButton onReview={handleReview} loading={loading} />
      <FileChangeLog groupedImages={groupedImages} cachedImages={cachedImages} />
      <div className="mt-auto">
        <SaveButton 
          onSave={handleSave} 
          loading={loading} 
          saveSuccess={saveSuccess} 
          cachedImages={cachedImages} 
          disabled={!selectedJob || !selectedDataset || loading} 
        />
      </div>
      
      <HomeReview 
        isOpen={isReviewOpen} 
        onClose={handleCloseReview} 
      />
    </div>
  );
}

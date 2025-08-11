"use client";

import React from 'react';
import { DatasetGridProps } from '@/types/HomeLeftSidebar';

export default function DatasetGrid({ 
  currentPagenation, datasetsPerPage, currentDatasetList, selectedPageIndex, 
  onDatasetSelect 
}: DatasetGridProps) {

  const getCurrentPageDatasets = () => {
    const startIndex = currentPagenation * datasetsPerPage;
    return currentDatasetList.slice(startIndex, startIndex + datasetsPerPage);
  };

  return (
    <div className="dataset-container">
      <p className="dataset-label">Select a Dataset:</p>
      <div className="dataset-grid">
        {getCurrentPageDatasets().map((dataset, idx) => {
          const absoluteIndex = currentPagenation * datasetsPerPage + idx;
          const isSelected = absoluteIndex === selectedPageIndex;
          
          return (
            <div 
              key       = {idx} 
              className = {`dataset-item ${isSelected ? 'dataset-item-selected' : ''}`}
              onClick   = {() => onDatasetSelect(dataset, idx)}
              title     = {dataset}
            >
              {absoluteIndex + 1}
            </div>
          );
        })}
      </div>
    </div>
  );
}
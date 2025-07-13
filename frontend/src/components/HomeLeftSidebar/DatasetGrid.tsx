"use client";

import React from 'react';
import { DatasetGridProps } from '@/types/HomeLeftSidebar';

export default function DatasetGrid({ 
  currentPagenation, datasetsPerPage, currentDatasets, selectedDataset, 
  onDatasetSelect 
}: DatasetGridProps) {

  const getCurrentPageDatasets = () => {
    const startIndex = currentPagenation * datasetsPerPage;
    return currentDatasets.slice(startIndex, startIndex + datasetsPerPage);
  };

  return (
    <div className="dataset-container">
      <p className="dataset-label">Select a Dataset:</p>
      <div className="dataset-grid">
        {getCurrentPageDatasets().map((dataset, i) => {
          const isSelected = dataset === selectedDataset;
          return (
            <div 
              key       = {dataset} 
              className = {`dataset-item ${isSelected ? 'dataset-item-selected' : ''}`}
              onClick   = {() => onDatasetSelect(dataset)}
              title     = {dataset}
            >
              {i + 1 + currentPagenation * datasetsPerPage}
            </div>
          );
        })}
      </div>
    </div>
  );
}

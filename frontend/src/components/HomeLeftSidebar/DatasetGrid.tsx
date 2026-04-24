"use client";

import React, { useEffect, useRef } from 'react';
import { DatasetGridProps } from '@/types/HomeLeftSidebar';

export default function DatasetGrid({ 
  currentPagenation, datasetsPerPage, currentDatasetList, selectedPageIndex, 
  onDatasetSelect 
}: DatasetGridProps) {

  const selectedItemRef = useRef<HTMLDivElement>(null);

  const getCurrentPageDatasets = () => {
    const startIndex = currentPagenation * datasetsPerPage;
    return currentDatasetList.slice(startIndex, startIndex + datasetsPerPage);
  };

  useEffect(() => {
    if (selectedItemRef.current) {
      selectedItemRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
      });
    }
  }, [selectedPageIndex]);

  return (
    <div className="overflow-auto mb-4 bg-white border border-gray-200 rounded p-1" style={{ height: 'calc(100vh - 180px)', minHeight: '300px', maxHeight: '1400px' }}>
      <p className="text-sm mb-2 text-gray-600 font-medium">Select a Dataset:</p>
      <div className="flex flex-col gap-2 max-w-full">
        {getCurrentPageDatasets().map((dataset, idx) => {
          const absoluteIndex = currentPagenation * datasetsPerPage + idx;
          const isSelected = absoluteIndex === selectedPageIndex;
          
          return (
            <div 
              key       = {idx}
              ref       = {isSelected ? selectedItemRef : null}
              className = {`w-auto min-h-8 h-auto border flex items-center justify-center text-xs cursor-pointer rounded-sm transition-all duration-200 font-medium px-2 py-1 text-center whitespace-normal break-words ${isSelected ? 'border-blue-500 bg-blue-100 text-blue-800 font-semibold hover:bg-blue-200' : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-100 hover:border-gray-400 hover:-translate-y-px hover:shadow-md'}`}
              onClick   = {() => onDatasetSelect(dataset, idx)}
              title     = {dataset}
            >
              {dataset + "_" + (absoluteIndex + 1)}
            </div>
          );
        })}
      </div>
    </div>
  );
}
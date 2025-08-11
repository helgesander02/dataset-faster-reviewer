"use client";

import React from 'react';
import { DATASET_PER_PAGE } from '@/services/config';
import { DatasetSectionProps } from '@/types/HomeLeftSidebar';

import DatasetGrid from './DatasetGrid';
import Pagination from './DatasetPagination';
import Status from './Status';


export function DatasetSection({
  currentPagenation, currentDatasetList, selectedPageIndex, selectedDataset, selectedJob,
  onDatasetSelect, onPrevious, onNext
}: DatasetSectionProps) {
  
  return (
    <>
      <DatasetGrid
        currentPagenation   = {currentPagenation}
        datasetsPerPage     = {DATASET_PER_PAGE}
        currentDatasetList  = {currentDatasetList}
        selectedPageIndex   = {selectedPageIndex}
        selectedDataset     = {selectedDataset}
        onDatasetSelect     = {onDatasetSelect}
      />
      <Status selectedJob={selectedJob} selectedDataset={selectedDataset} />
      {currentDatasetList.length > 0 && (
        <Pagination
          currentPagenation = {currentPagenation}
          totalDatasets     = {currentDatasetList.length}
          datasetsPerPage   = {DATASET_PER_PAGE} 
          onPrevious        = {onPrevious}
          onNext            = {onNext}
        />
      )}
    </>
  );
}

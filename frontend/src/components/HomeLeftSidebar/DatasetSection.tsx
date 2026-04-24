"use client";

import React from 'react';
import { DATASET_PER_PAGE } from '@/services/config';
import { DatasetSectionProps } from '@/types/HomeLeftSidebar';

import DatasetGrid from './DatasetGrid';
import Pagination from './DatasetPagination';


export function DatasetSection({
  currentPagenation, currentDatasetList, selectedPageIndex, selectedDataset,
  onDatasetSelect, onPrevious, onNext, onGoToPage
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
      
      {currentDatasetList.length > 0 && (
        <Pagination
          currentPagenation = {currentPagenation}
          totalDatasets     = {currentDatasetList.length}
          datasetsPerPage   = {DATASET_PER_PAGE} 
          onPrevious        = {onPrevious}
          onNext            = {onNext}
          onGoToPage        = {onGoToPage}
        />
      )}
    </>
  );
}

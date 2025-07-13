"use client";

import React from 'react';
import { DATASET_PER_PAGE } from '@/services/config';
import { DatasetSectionProps } from '@/types/HomeLeftSidebar';

import DatasetGrid from './DatasetGrid';
import Pagination from './DatasetPagination';


export function DatasetSection({
  currentPagenation, currentDatasets, selectedDataset,
  onDatasetSelect, onPrevious, onNext
}: DatasetSectionProps) {
  
  return (
    <>
      <DatasetGrid
        currentPagenation = {currentPagenation}
        datasetsPerPage   = {DATASET_PER_PAGE}
        currentDatasets   = {currentDatasets}
        selectedDataset   = {selectedDataset}
        onDatasetSelect   = {onDatasetSelect}
      />

      {currentDatasets.length > 0 && (
        <Pagination
          currentPagenation = {currentPagenation}
          totalDatasets     = {currentDatasets.length}
          datasetsPerPage   = {DATASET_PER_PAGE} 
          onPrevious        = {onPrevious}
          onNext            = {onNext}
        />
      )}
    </>
  );
}

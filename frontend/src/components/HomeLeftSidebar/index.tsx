"use client";

import { useLeftSidebar } from '@/hooks/useLeftSidebar/useLeftSidebar';

import { DATASET_PER_PAGE } from '@/services/config';
import { useJobDataset } from '../JobDatasetContext';
import { SidebarHeader } from './Header';
import { DatasetSection } from './DatasetSection';
import JobSelect from './JobSelect';
import LoadingIndicator from './LoadingIndicator';

export default function LeftSidebar() {

  const { 
    selectedJob, selectedDataset, selectedPageIndex,
    setSelectedJob, setSelectedDataset, setselectedPageIndex
   } = useJobDataset();

  const {
    currentJobList, currentDatasetList, currentPagenation, loading,
    previousPage, nextPage, goToPage
  } = useLeftSidebar();

  const onDatasetSelect = (dataset: string, relativeIdx: number) => {
    const absoluteIndex = currentPagenation * DATASET_PER_PAGE + relativeIdx;
    setselectedPageIndex(absoluteIndex);
    setSelectedDataset(dataset);
  };

  return (
    <div className="w-[17%] bg-gray-100 p-2 min-w-[280px]">
      <SidebarHeader />
      
      <JobSelect 
        currentJobList = {currentJobList} 
        selectedJob    = {selectedJob} 
        loading        = {loading} 
        onJobSelect    = {setSelectedJob} 
      />
      
      {selectedJob && (
        <DatasetSection
          currentPagenation   = {currentPagenation}
          currentDatasetList  = {currentDatasetList}
          selectedPageIndex   = {selectedPageIndex}
          selectedDataset     = {selectedDataset} 
          selectedJob         = {selectedJob}
          onDatasetSelect     = {onDatasetSelect}
          onPrevious          = {previousPage}
          onNext              = {nextPage}
          onGoToPage          = {goToPage}
        />
      )}
      {loading && <LoadingIndicator />}
    </div>
  );
}
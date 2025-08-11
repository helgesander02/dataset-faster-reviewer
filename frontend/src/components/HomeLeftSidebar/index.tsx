"use client";

import { useLeftSidebar } from '@/hooks/useLeftSidebar/useLeftSidebar';
import '@/styles/HomeLeftSidebar.css';

import { useJobDataset } from '../JobDatasetContext';
import { SidebarHeader } from './Header';
import { DatasetSection } from './DatasetSection';
import JobSelect from './JobSelect';
import LoadingIndicator from './LoadingIndicator';

/**
 * LeftSidebar component for the Home page.
 * This component displays a sidebar with job selection and dataset management features.
 * It utilizes context to manage the selected job and dataset.
 * 
 * Props:
 * - currentJobList: Array of job names to display.
 * - currentDatasetList: Array of dataset names to display.
 * - currentPagenation: Current pagination index for datasets.
 * - loading: Boolean indicating if data is being loaded.
 * - previousPage: Function to navigate to the previous page of datasets.
 * - nextPage: Function to navigate to the next page of datasets.
 */
export default function LeftSidebar() {

  const { 
    selectedJob, selectedDataset, selectedPageIndex,
    setSelectedJob, setSelectedDataset, setselectedPageIndex
   } = useJobDataset();

  const {
    currentJobList, currentDatasetList, currentPagenation, loading,
    previousPage, nextPage
  } = useLeftSidebar();

  const onDatasetSelect = (dataset: string, relativeIdx: number) => {
    // Calculate absolute index based on current pagination and relative index
    const absoluteIndex = currentPagenation * 12 + relativeIdx; // Assuming DATASET_PER_PAGE is 12
    setselectedPageIndex(absoluteIndex);
    setSelectedDataset(dataset);
  };

  return (
    <div className="sidebar-container">
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
        />
      )}
      {loading && <LoadingIndicator />}
    </div>
  );
}
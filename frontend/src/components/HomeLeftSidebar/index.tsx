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
 * - currentJobs:         string[] - List of jobs available for selection.
 * - currentDatasets:     string[] - List of datasets available for the selected job.
 * - currentPagenation:   number - Pagination information for the datasets.
 * - loading:             boolean - Indicates if data is currently being loaded.
 * - handleJobSelect:     function - Callback function to handle job selection.
 * - handleDatasetSelect: function - Callback function to handle dataset selection.
 * - previousPage:        function - Callback function to navigate to the previous page of datasets.
 * - nextPage:            function - Callback function to navigate to the next page of datasets.
 */
export default function LeftSidebar() {

  const { selectedJob, selectedDataset } = useJobDataset();
  const {
    currentJobs, currentDatasets, currentPagenation, loading,
    handleJobSelect, handleDatasetSelect, previousPage, nextPage
  } = useLeftSidebar();

  return (
    <div className="sidebar-container">
      <SidebarHeader />
      
      <JobSelect 
        currentJobs = {currentJobs} 
        selectedJob = {selectedJob} 
        loading     = {loading} 
        onJobSelect = {handleJobSelect} 
      />
      
      {selectedJob && (
        <DatasetSection
          currentPagenation = {currentPagenation}
          currentDatasets   = {currentDatasets}
          selectedDataset   = {selectedDataset}
          onDatasetSelect   = {handleDatasetSelect}
          onPrevious        = {previousPage}
          onNext            = {nextPage}
        />
      )}
      
      {loading && <LoadingIndicator />}
    </div>
  );
}

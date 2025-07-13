"use client";

import { useJobDataset } from '@/components/JobDatasetContext';
import InfiniteImageGrid from '@/components/HomeImageGrid/InfinitelmageGrid';

export default function Home() {
  const { selectedJob, selectedDataset, currentPage, setSelectedDataset, setCurrentPage } = useJobDataset();

  return (
    <div className="main-container">
      <InfiniteImageGrid 
        selectedJob={selectedJob}
        selectedDataset={selectedDataset}
        currentPage={currentPage}
        setCurrentPage={setCurrentPage}
        setSelectedDataset={setSelectedDataset}
      />
    </div>
  );
}

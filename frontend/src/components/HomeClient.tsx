"use client";

import { JobDatasetProvider, useJobDataset } from '@/components/JobDatasetContext';
import LeftSidebar from '@/components/HomeLeftSidebar/index';
import RightSidebar from '@/components/HomeRightSidebar/index';
import InfiniteImageGrid from '@/components/HomeImageGrid/InfinitelmageGrid';

function HomeContent() {
  const { 
    selectedPages, 
    selectedDataset, 
    selectedPageIndex, 
    setSelectedDataset, 
    setselectedPageIndex 
  } = useJobDataset();

  return (
    <main className="flex h-screen">
      <LeftSidebar />
      <div className="flex-1 main-container">
        <InfiniteImageGrid 
          selectedPages={selectedPages}
          selectedDataset={selectedDataset}
          selectedPageIndex={selectedPageIndex}
          setselectedPageIndex={setselectedPageIndex}
          setSelectedDataset={setSelectedDataset}
        />
      </div>
      <RightSidebar />
    </main>
  );
}

export default function HomeClient() {
  return (
    <JobDatasetProvider>
      <HomeContent />
    </JobDatasetProvider>
  );
}

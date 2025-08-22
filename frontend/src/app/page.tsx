"use client";

import { useJobDataset } from '@/components/JobDatasetContext';
import InfiniteImageGrid from '@/components/HomeImageGrid/InfinitelmageGrid';

export default function Home() {
  const { 
    selectedPages, selectedDataset, selectedPageIndex, 
    setSelectedDataset, setselectedPageIndex 
  } = useJobDataset();

  return (
    <div className="main-container">
      <InfiniteImageGrid 
        selectedPages         ={selectedPages}
        selectedDataset       ={selectedDataset}
        selectedPageIndex     ={selectedPageIndex}
        setselectedPageIndex  ={setselectedPageIndex}
        setSelectedDataset    ={setSelectedDataset}
      />
    </div>
  );
}

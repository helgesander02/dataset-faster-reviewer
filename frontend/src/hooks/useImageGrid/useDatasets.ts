"use client";

import { useState, useEffect } from 'react';
import { fetchDatasets } from '@/services/api';

export function useDatasets(
  selectedJob: string | null
) {

  const [DatasetList, setDatasetList] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  // This effect fetches datasets when the selected job changes
  useEffect(() => {
    const loadDatasets = async () => {
      if (!selectedJob) {
        setDatasetList([]);
        return;
      }

      try {
        setLoading(true);
        const response = await fetchDatasets(selectedJob);
        const datasets = response.dataset_names || [];
        setDatasetList(datasets);

      } catch (error) {
        console.error('Error loading datasets:', error);
        setDatasetList([]);

      } finally {
        setLoading(false);
      }
    };

    loadDatasets();
  }, [selectedJob]);

  return { 
    DatasetList, 
    loading 
  };
}

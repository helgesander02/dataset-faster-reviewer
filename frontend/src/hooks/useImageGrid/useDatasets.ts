"use client";

import { useState, useEffect } from 'react';
import { fetchDatasets } from '@/services/api';

export function useDatasets(
  selectedJob: string | null
) {

  const [allDatasets, setAllDatasets] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  // This effect fetches datasets when the selected job changes
  useEffect(() => {
    const loadDatasets = async () => {
      if (!selectedJob) {
        setAllDatasets([]);
        return;
      }

      try {
        setLoading(true);
        const response = await fetchDatasets(selectedJob);
        const datasets = response.dataset_names || [];
        setAllDatasets(datasets);

      } catch (error) {
        console.error('Error loading datasets:', error);
        setAllDatasets([]);

      } finally {
        setLoading(false);
      }
    };

    loadDatasets();
  }, [selectedJob]);

  return { 
    allDatasets, 
    loading 
  };
}

// ============================================================================
// Cached Image Types
// ============================================================================

export interface CachedImage {
  item_job_name: string;
  item_dataset_name: string;
  item_image_name: string;
  item_image_path: string;
}

// ============================================================================
// Context Type
// ============================================================================

export interface JobDatasetContextType {
  // Selection State
  selectedJob: string;
  selectedPages: string;
  selectedDataset: string;
  selectedPageIndex: number;

  // Selection Actions
  setSelectedJob: (job: string) => void;
  setSelectedPages: (pages: string) => void;
  setSelectedDataset: (dataset: string) => void;
  setselectedPageIndex: (page: number) => void;

  // Cached Images State
  cachedImages: CachedImage[];

  // Cache Management Actions
  addImageToCache: (
    job: string,
    dataset: string,
    imageName: string,
    imagePath: string
  ) => void;
  removeImageFromCache: (imagePath: string) => void;
  clearCache: () => void;
  getCache: (job: string, dataset: string) => string[];
  getCacheCountForJob: (job: string) => number;
}

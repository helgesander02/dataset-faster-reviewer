export interface CachedImage {
  job: string;
  dataset: string;
  imageName: string;
  imagePath: string;
}

export interface JobDatasetContextType {
  selectedJob: string;
  selectedDataset: string;
  currentPage: number;
  setSelectedJob: (job: string) => void;
  setSelectedDataset: (dataset: string) => void;
  setCurrentPage: (page: number) => void;
  cachedImages: CachedImage[];
  addImageToCache: (job: string, dataset: string, imageName:string, imagePath: string) => void;
  removeImageFromCache: (imagePath: string) => void;
  getCache: (job: string, dataset: string) => string[];
  jobPages: string[];
}

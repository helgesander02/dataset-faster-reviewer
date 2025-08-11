export interface CachedImage {
  item_job_name: string;
  item_dataset_name: string;
  item_image_name: string;
  item_image_path: string;
}

export interface JobDatasetContextType {
  selectedJob: string;
  selectedPages: string;
  selectedDataset: string;
  selectedPageIndex: number;
  setSelectedJob: (job: string) => void;
  setSelectedPages: (pages: string) => void;
  setSelectedDataset: (dataset: string) => void;
  setselectedPageIndex: (page: number) => void;
  cachedImages: CachedImage[];
  addImageToCache: (job: string, dataset: string, imageName:string, imagePath: string) => void;
  removeImageFromCache: (imagePath: string) => void;
  getCache: (job: string, dataset: string) => string[];
}

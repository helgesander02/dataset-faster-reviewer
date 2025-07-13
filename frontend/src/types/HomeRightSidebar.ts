export interface CachedImage {
  job:              string;
  dataset:          string;
  imageName:        string;
  imagePath:        string;
}

export interface SaveData {
  job:              string;
  dataset:          string;
  images:           Array<{
    job:              string;
    dataset:           string;
    imageName:        string;
    imagePath:        string;
  }>;
  timestamp:        string;
}

export interface FileChangeLogProps {
  groupedImages:    Record<string, CachedImage[]>;
  cachedImages:     CachedImage[];
}

export interface ReviewButtonProps {
  loading:          boolean;
  onReview:         () => void;
}

export interface SaveButtonProps {
  loading:          boolean;
  saveSuccess:      boolean;
  cachedImages:     CachedImage[];
  disabled:         boolean;
  onSave:           () => void;
}

export interface StatusProps {
  selectedJob:      string | null;
  selectedDataset:  string | null;
}

export interface RightSidebarState {
  isReviewOpen:     boolean;
  loading:          boolean;
  saveSuccess:      boolean;
  groupedImages:    Record<string, CachedImage[]>;
}

export interface RightSidebarActions {
  handleSave:        () => Promise<void>;
  handleReview:      () => void;
  handleCloseReview: () => void;
}

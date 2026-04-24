// useRightSidebar/useRightSidebar.ts
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


export interface CachedImage {
  item_job_name:          string;
  item_dataset_name:      string;
  item_image_name:        string;
  item_image_path:        string;
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

// FileChangeLog.tsx
export interface FileChangeLogProps {
  groupedImages:    Record<string, CachedImage[]>;
  cachedImages:     CachedImage[];
}

// ReviewButton.tsx
export interface ReviewButtonProps {
  loading:          boolean;
  onReview:         () => void;
}

// SaveButton.tsx
export interface SaveButtonProps {
  loading:          boolean;
  saveSuccess:      boolean;
  cachedImages:     CachedImage[];
  disabled:         boolean;
  onSave:           () => void;
}

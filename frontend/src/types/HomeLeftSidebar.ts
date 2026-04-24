// ============================================================================
// useLeftSidebar Hook Types
// ============================================================================

export interface SidebarState {
  currentJobList: string[];
  currentDatasetList: string[];
  currentPagenation: number;
  loading: boolean;
  error: string | null;
  
  // Computed values
  totalDatasetPages: number;
  hasPreviousPage: boolean;
  hasNextPage: boolean;
  currentPageDatasets: string[];
}

export interface SidebarActions {
  previousPage: () => void;
  nextPage: () => void;
  goToPage: (page: number) => void;
  loadJobs: () => Promise<void>;
  loadDatasets: () => Promise<void>;
}

// ============================================================================
// Component Props
// ============================================================================

// JobSelect.ts
export interface JobSelectProps {
  currentJobList: string[];
  selectedJob: string;
  loading: boolean;
  error?: string | null;
  onJobSelect: (job: string) => void;
}

// DatasetSection.ts
export interface DatasetSectionProps {
  currentPagenation: number;
  currentDatasetList: string[];
  selectedPageIndex: number;
  selectedDataset: string;
  selectedJob: string;
  onDatasetSelect: (dataset: string, idx: number) => void;
  onPrevious: () => void;
  onNext: () => void;
  onGoToPage?: (page: number) => void;
}

// DatasetGrid.ts
export interface DatasetGridProps {
  currentPagenation: number;
  datasetsPerPage: number;
  currentDatasetList: string[];
  selectedPageIndex: number;
  selectedDataset: string;
  onDatasetSelect: (dataset: string, idx: number) => void;
}

// DatasetPagination.ts
export interface PaginationProps {
  currentPagenation: number;
  totalDatasets: number;
  datasetsPerPage: number;
  hasPrevious?: boolean;
  hasNext?: boolean;
  onPrevious: () => void;
  onNext: () => void;
  onGoToPage?: (page: number) => void;
}

// Status.ts
export interface StatusProps {
  selectedJob: string | null;
  selectedDataset: string | null;
}

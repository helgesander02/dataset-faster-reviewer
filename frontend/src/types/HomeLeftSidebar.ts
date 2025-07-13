// useLeftSidebar.ts
export interface SidebarState {
    currentJobs:            string[];
    currentDatasets:        string[];
    currentPagenation:      number;
    loading:                boolean;  
}

export interface SidebarActions {
    handleJobSelect:        (job: string) => void;
    handleDatasetSelect:    (dataset: string) => void;
    previousPage:           () => void;
    nextPage:               () => void;
}

// JobSelect.ts
export interface JobSelectProps {
    currentJobs:            string[];
    selectedJob:            string;
    loading:                boolean;
    onJobSelect:            (job: string) => void;
}

// DatasetSection.ts
export interface DatasetSectionProps {
    currentPagenation:      number; 
    currentDatasets:        string[]; 
    selectedDataset:        string;
    onDatasetSelect:        (dataset: string) => void; 
    onPrevious:             () => void; 
    onNext:                 () => void;
}

// DatasetGrid.ts
export interface DatasetGridProps {
    currentPagenation:      number;
    datasetsPerPage:        number;
    currentDatasets:        string[];
    selectedDataset:        string;
    onDatasetSelect:        (dataset: string) => void;
}

// DatasetPagination.ts
export interface PaginationProps {
    currentPagenation:      number;
    totalDatasets:          number;
    datasetsPerPage:        number;
    onPrevious:             () => void;
    onNext:                 () => void;
}

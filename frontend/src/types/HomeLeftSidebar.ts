// useLeftSidebar.ts
export interface SidebarState {
    currentJobList:         string[];
    currentDatasetList:     string[];
    currentPagenation:      number;
    loading:                boolean;  
}

export interface SidebarActions {
    previousPage:           () => void;
    nextPage:               () => void;
}

// JobSelect.ts
export interface JobSelectProps {
    currentJobList:         string[];
    selectedJob:            string;
    loading:                boolean;
    onJobSelect:            (job: string) => void;
}

// DatasetSection.ts
export interface DatasetSectionProps {
    currentPagenation:      number; 
    currentDatasetList:     string[];
    selectedPageIndex:      number; 
    selectedDataset:        string;
    selectedJob:            string;
    onDatasetSelect:        (dataset: string, idx: number) => void;
    onPrevious:             () => void; 
    onNext:                 () => void;
}

// DatasetGrid.ts
export interface DatasetGridProps {
    currentPagenation:      number;
    datasetsPerPage:        number;
    currentDatasetList:     string[];
    selectedPageIndex:      number;
    selectedDataset:        string;
    onDatasetSelect:        (dataset: string, idx: number) => void;
}

// DatasetPagination.ts
export interface PaginationProps {
    currentPagenation:      number;
    totalDatasets:          number;
    datasetsPerPage:        number;
    onPrevious:             () => void;
    onNext:                 () => void;
}

// Status.ts
export interface StatusProps {
    selectedJob:      string | null;
    selectedDataset:  string | null;
}

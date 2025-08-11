export interface InfiniteImageGridProps {
    selectedPages:        string;
    selectedDataset:      string;
    selectedPageIndex:    number;
    setSelectedDataset:   (dataset: string) => void;
    setselectedPageIndex: (page: number) => void;
}

export interface UseInfiniteImagesReturn {
    allImagePages: ImagePage[];
    loading: boolean;
    currentPageIndex: number; 
    getCurrentImagePages: () => ImagePage[];
    hasMorePages: () => boolean;
    loadNextPage: () => Promise<void>;
    resetImages: () => void;
    initializeImages: () => Promise<void>;
    registerPageElement: (pageIndex: number, element: HTMLElement | null) => void; 
}

export interface Image {
    name: string;
    url: string;
    dataset: string;
}

export interface HomeImageGridProps {
    images: Image[];
    selectedImages: Set<string>;
    isLoading: boolean;
    onImageClick: (imageName: string, imageUrl: string, dataset: string) => void;
}

export interface ImageItemProps {
    image: Image;
    index: number;
    isSelected: boolean;
    onImageClick: (imageName: string, imageUrl: string, dataset: string) => void;
}

export interface LoadingStateProps {
    message?: string;
}

export interface SelectionIndicatorProps {
    className?: string;
}

export interface ImagePage {
    dataset: string;
    images: Image[];
    isNewDataset: boolean;
}

export interface EmptyStateProps {
    title?: string;
    message?: string;
}

export interface LoadingTriggerProps {
    isLoading: boolean;
    loadingMessage?: string;
}

export type DatasetImageCountsMap = Map<string, number>;
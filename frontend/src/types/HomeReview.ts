// index.tsx
export interface HomeReviewProps {
    isOpen: boolean;
    onClose: () => void;
}

// ReviewHeader.tsx
export interface ReviewHeaderProps {
    saving: boolean;
    onClose: () => void;
}

// ReviewContent.tsx
export interface ReviewContentProps {
    loading: boolean;
    error: string | null;
    reviewData: PendingReviewData | null;
    selectedImages: Set<string>;
    onRetry: () => void;
    onToggleImage: (item: ReviewItem) => void;
    loadedItems?: ReviewItemWithUrl[];
    onLoadMore?: () => void;
    hasMorePages?: boolean;
    totalItems?: number;
}

// ReviewItem with display URL
export interface ReviewItemWithUrl extends ReviewItem {
    displayUrl?: string;
}

// ImagesGrid.tsx
export interface ImagesGridProps {
    items: ReviewItem[];
    selectedImages: Set<string>;
    onToggleImage: (item: ReviewItem) => void;
    loadedItems?: ReviewItemWithUrl[];
    onLoadMore?: () => void;
    hasMorePages?: boolean;
}

// ImageItem.tsx
export interface ImageItemProps {
    item: ReviewItem;
    index: number;
    isSelected: boolean;
    allItems: ReviewItem[];
    onToggle: (item: ReviewItem) => void;
}

// ReviewActions.tsx
export interface ReviewActionsProps {
    selectedCount: number;
    totalCount: number;
    saving: boolean;
    deleting: boolean;
    onSelectAll: () => void;
    onDeselectAll: () => void;
    onSave: () => void;
    onDelete: () => void;
}

// useHomeReview.ts
export interface ReviewItem {
    item_job_name:       string;
    item_dataset_name:   string;
    item_image_name:     string;
    item_image_path:     string;      // Base64 thumbnail
}

export interface PendingReviewData {
    items: ReviewItem[];
}

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
}

// ImagesGrid.tsx
export interface ImagesGridProps {
    items: ReviewItem[];
    selectedImages: Set<string>;
    onToggleImage: (item: ReviewItem) => void;
}

// ImageItem.tsx
export interface ImageItemProps {
    item: ReviewItem;
    index: number;
    isSelected: boolean;
    onToggle: (item: ReviewItem) => void;
}

// ReviewActions.tsx
export interface ReviewActionsProps {
    selectedCount: number;
    totalCount: number;
    saving: boolean;
    onSelectAll: () => void;
    onDeselectAll: () => void;
    onSave: () => void;
}

// useHomeReview.ts
export interface ReviewItem {
    job: string;
    dataset: string;
    imageName: string;
    imagePath: string;
}

export interface PendingReviewData {
    items: ReviewItem[];
}

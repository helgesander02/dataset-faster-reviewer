/**
 * EmptyState 元件 - 空狀態顯示 (改進版)
 * 
 * 特性:
 * - 友善的視覺設計
 * - 清晰的行動指引
 * - 支援自訂圖標
 * - 可選操作按鈕
 */

import React from 'react';
import { 
  FileQuestion, 
  FolderOpen, 
  Search, 
  Image as ImageIcon,
  Inbox,
  LucideIcon 
} from 'lucide-react';
import Button from './Button';

export interface EmptyStateProps {
  variant?: 'default' | 'no-data' | 'no-results' | 'no-images' | 'no-selection';
  icon?: LucideIcon | React.ReactNode;
  title?: string;
  description?: string;
  action?: {
    label: string;
    onClick: () => void;
    icon?: React.ReactNode;
  };
  className?: string;
}

const EmptyState: React.FC<EmptyStateProps> = ({
  variant = 'default',
  icon,
  title,
  description,
  action,
  className = '',
}) => {
  // 預設圖標配置
  const variantConfig = {
    default: {
      icon: Inbox,
      title: '目前沒有內容',
      description: '這裡還沒有任何資料。',
    },
    'no-data': {
      icon: FolderOpen,
      title: '沒有資料',
      description: '目前沒有可顯示的資料。',
    },
    'no-results': {
      icon: Search,
      title: '找不到結果',
      description: '請嘗試其他搜尋條件。',
    },
    'no-images': {
      icon: ImageIcon,
      title: '沒有圖片',
      description: '此資料集中沒有圖片。',
    },
    'no-selection': {
      icon: FileQuestion,
      title: '請選擇項目',
      description: '從左側邊欄選擇一個項目以開始。',
    },
  };

  const config = variantConfig[variant];
  const IconComponent = icon || config.icon;
  const displayTitle = title || config.title;
  const displayDescription = description || config.description;

  return (
    <div
      className={`flex flex-col items-center justify-center p-12 text-center min-h-[400px] ${className}`}
      role="status"
      aria-label={displayTitle}
    >
      {/* 圖標 */}
      <div className="mb-6 text-gray-300" aria-hidden="true">
        {typeof IconComponent === 'function' ? (
          <IconComponent size={64} strokeWidth={1.5} />
        ) : (
          IconComponent
        )}
      </div>

      {/* 標題 */}
      <h3 className="text-xl font-semibold text-gray-900 mb-3">
        {displayTitle}
      </h3>

      {/* 描述 */}
      <p className="text-sm text-gray-600 mb-6 max-w-md">
        {displayDescription}
      </p>

      {/* 操作按鈕 */}
      {action && (
        <Button
          variant="primary"
          icon={action.icon}
          onClick={action.onClick}
        >
          {action.label}
        </Button>
      )}
    </div>
  );
};

/**
 * 特定場景的空狀態元件
 */

export const NoJobSelected: React.FC = () => (
  <EmptyState
    variant="no-selection"
    title="歡迎使用 Image Verify Viewer"
    description="請從左側邊欄選擇一個 Job 以開始審查圖片。"
  />
);

export const NoDatasetFound: React.FC = () => (
  <EmptyState
    variant="no-data"
    title="找不到資料集"
    description="此 Job 中沒有可用的資料集。"
  />
);

export const NoImagesFound: React.FC = () => (
  <EmptyState
    variant="no-images"
    title="沒有圖片"
    description="此資料集中沒有圖片可供審查。"
  />
);

export const NoSelectionMade: React.FC = () => (
  <EmptyState
    variant="no-selection"
    title="尚未選擇圖片"
    description="請從圖片網格中選擇至少一張圖片。"
  />
);

export default EmptyState;

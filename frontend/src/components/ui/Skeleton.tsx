/**
 * Skeleton 元件 - 骨架屏載入狀態
 * 
 * 用途:
 * - 內容載入時顯示佔位符
 * - 改善感知性能
 * - 減少用戶焦慮
 * 
 * 最佳實踐:
 * - 與實際內容相似的形狀和大小
 * - 使用脈衝動畫
 */

import React from 'react';

export interface SkeletonProps {
  variant?: 'text' | 'circular' | 'rectangular' | 'image';
  width?: string | number;
  height?: string | number;
  className?: string;
  count?: number;
}

const Skeleton: React.FC<SkeletonProps> = ({
  variant = 'text',
  width,
  height,
  className = '',
  count = 1,
}) => {
  // 變體樣式
  const variantStyles = {
    text: 'h-4 w-full rounded',
    circular: 'rounded-full',
    rectangular: 'rounded-md',
    image: 'aspect-square rounded-lg',
  };

  // 生成樣式
  const style: React.CSSProperties = {};
  if (width) style.width = typeof width === 'number' ? `${width}px` : width;
  if (height) style.height = typeof height === 'number' ? `${height}px` : height;

  const skeletonElement = (
    <div
      className={`skeleton ${variantStyles[variant]} ${className}`}
      style={style}
      aria-hidden="true"
    />
  );

  // 如果需要多個骨架屏
  if (count > 1) {
    return (
      <div className="space-y-2">
        {Array.from({ length: count }).map((_, index) => (
          <React.Fragment key={index}>{skeletonElement}</React.Fragment>
        ))}
      </div>
    );
  }

  return skeletonElement;
};

/**
 * SkeletonImage - 圖片骨架屏
 */
export const SkeletonImage: React.FC<{ className?: string }> = ({ className = '' }) => (
  <div className={`skeleton aspect-square rounded-lg ${className}`} aria-hidden="true">
    <div className="w-full h-full flex items-center justify-center">
      <svg
        className="w-12 h-12 text-gray-300"
        fill="currentColor"
        viewBox="0 0 20 20"
      >
        <path
          fillRule="evenodd"
          d="M4 3a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V5a2 2 0 00-2-2H4zm12 12H4l4-8 3 6 2-4 3 6z"
          clipRule="evenodd"
        />
      </svg>
    </div>
  </div>
);

/**
 * SkeletonCard - 卡片骨架屏
 */
export const SkeletonCard: React.FC<{ className?: string }> = ({ className = '' }) => (
  <div className={`card-base p-4 ${className}`} aria-hidden="true">
    <Skeleton variant="rectangular" height={120} className="mb-3" />
    <Skeleton variant="text" className="mb-2" />
    <Skeleton variant="text" width="60%" />
  </div>
);

/**
 * SkeletonImageGrid - 圖片網格骨架屏
 */
export const SkeletonImageGrid: React.FC<{ count?: number }> = ({ count = 12 }) => (
  <div 
    className="grid grid-cols-[repeat(auto-fill,minmax(180px,1fr))] gap-3"
    aria-busy="true"
    aria-label="Loading images"
  >
    {Array.from({ length: count }).map((_, index) => (
      <div key={index} className="card-base overflow-hidden">
        <SkeletonImage />
        <div className="p-2">
          <Skeleton variant="text" height={12} />
        </div>
      </div>
    ))}
  </div>
);

export default Skeleton;

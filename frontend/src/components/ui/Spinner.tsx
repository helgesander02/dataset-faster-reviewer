/**
 * Spinner 元件 - 載入動畫指示器
 * 
 * 用途:
 * - 顯示載入中狀態
 * - 非阻塞性的進度指示
 * 
 * 使用場景:
 * - 按鈕內部載入狀態
 * - 卡片內容載入
 * - 頁面局部刷新
 */

import React from 'react';
import { Loader2 } from 'lucide-react';

export interface SpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  color?: 'primary' | 'white' | 'gray';
  className?: string;
  label?: string;
}

const Spinner: React.FC<SpinnerProps> = ({
  size = 'md',
  color = 'primary',
  className = '',
  label = 'Loading...',
}) => {
  // 大小映射
  const sizeMap = {
    sm: 16,
    md: 24,
    lg: 32,
    xl: 48,
  };

  // 顏色樣式
  const colorStyles = {
    primary: 'text-primary',
    white: 'text-white',
    gray: 'text-gray-500',
  };

  return (
    <div className={`inline-flex items-center justify-center ${className}`} role="status">
      <Loader2 
        className={`spinner ${colorStyles[color]}`} 
        size={sizeMap[size]}
        aria-hidden="true"
      />
      <span className="sr-only">{label}</span>
    </div>
  );
};

export default Spinner;

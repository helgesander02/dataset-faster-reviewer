/**
 * ErrorState 元件 - 錯誤狀態顯示
 * 
 * 特性:
 * - 清晰的錯誤圖標和訊息
 * - 提供建議性的操作指引
 * - 支援重試功能
 * - 可訪問性優化
 */

import React from 'react';
import { AlertCircle, RefreshCw } from 'lucide-react';
import Button from './Button';
import { logger } from '@/utils/logger';

export interface ErrorStateProps {
  title?: string;
  message?: string;
  error?: Error | string;
  onRetry?: () => void;
  retryLabel?: string;
  showIcon?: boolean;
  className?: string;
}

const ErrorState: React.FC<ErrorStateProps> = ({
  title = '載入失敗',
  message = '無法載入內容,請稍後再試。',
  error,
  onRetry,
  retryLabel = '重試',
  showIcon = true,
  className = '',
}) => {
  // 處理錯誤訊息
  const errorMessage = error 
    ? typeof error === 'string' 
      ? error 
      : error.message 
    : message;

  return (
    <div
      className={`flex flex-col items-center justify-center p-8 text-center ${className}`}
      role="alert"
      aria-live="assertive"
    >
      {showIcon && (
        <div className="mb-4 text-error" aria-hidden="true">
          <AlertCircle size={48} strokeWidth={1.5} />
        </div>
      )}

      <h3 className="text-lg font-semibold text-gray-900 mb-2">
        {title}
      </h3>

      <p className="text-sm text-gray-600 mb-6 max-w-md">
        {errorMessage}
      </p>

      {onRetry && (
        <Button
          variant="primary"
          icon={<RefreshCw size={16} />}
          onClick={onRetry}
        >
          {retryLabel}
        </Button>
      )}

      {/* 開發模式下顯示完整錯誤 */}
      {process.env.NODE_ENV === 'development' && error && typeof error !== 'string' && (
        <details className="mt-4 text-left w-full max-w-md">
          <summary className="text-xs text-gray-500 cursor-pointer hover:text-gray-700">
            顯示錯誤詳情
          </summary>
          <pre className="mt-2 text-xs bg-gray-100 p-3 rounded overflow-auto max-h-40">
            {error.stack || error.message}
          </pre>
        </details>
      )}
    </div>
  );
};

/**
 * ErrorBoundary 組合 - 用於捕獲 React 錯誤
 */
export class ErrorBoundary extends React.Component<
  { children: React.ReactNode; fallback?: React.ReactNode },
  { hasError: boolean; error?: Error }
> {
  constructor(props: { children: React.ReactNode; fallback?: React.ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  override componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    logger.error('ErrorBoundary caught an error:', error, errorInfo);
  }

  override render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }
      
      return (
        <ErrorState
          title="發生錯誤"
          message="應用程式遇到了意外錯誤"
          error={this.state.error}
          onRetry={() => this.setState({ hasError: false })}
        />
      );
    }

    return this.props.children;
  }
}

export default ErrorState;

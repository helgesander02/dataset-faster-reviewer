/**
 * Button 元件 - 統一按鈕樣式
 * 
 * 特性:
 * - 多種變體: primary, secondary, success, danger
 * - 多種大小: sm, md, lg
 * - 載入狀態支援
 * - 圖標支援
 * - 完整的可訪問性
 */

import React from 'react';
import { Loader2 } from 'lucide-react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'success' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  icon?: React.ReactNode;
  fullWidth?: boolean;
  children: React.ReactNode;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      variant = 'primary',
      size = 'md',
      loading = false,
      icon,
      fullWidth = false,
      disabled,
      children,
      className = '',
      ...props
    },
    ref
  ) => {
    // 基礎樣式
    const baseStyles = 'btn-base focus-ring';

    // 變體樣式
    const variantStyles = {
      primary: 'bg-primary text-on-primary hover:bg-primary-hover active:bg-primary-active focus-visible:outline-primary',
      secondary: 'bg-gray-100 text-gray-700 hover:bg-gray-200 active:bg-gray-300 border border-gray-300 focus-visible:outline-gray-500',
      success: 'bg-success text-on-success hover:bg-success-hover active:bg-success-hover focus-visible:outline-success',
      danger: 'bg-error text-on-error hover:bg-error-hover active:bg-error-hover focus-visible:outline-error',
      ghost: 'bg-transparent text-gray-700 hover:bg-gray-100 active:bg-gray-200 focus-visible:outline-gray-500',
    };

    // 大小樣式
    const sizeStyles = {
      sm: 'text-xs px-3 py-1.5',
      md: 'text-sm px-4 py-2',
      lg: 'text-base px-6 py-3',
    };

    // 寬度樣式
    const widthStyles = fullWidth ? 'w-full' : '';

    // 合併樣式
    const combinedClassName = `
      ${baseStyles}
      ${variantStyles[variant]}
      ${sizeStyles[size]}
      ${widthStyles}
      ${className}
    `.trim();

    return (
      <button
        ref={ref}
        disabled={disabled || loading}
        className={combinedClassName}
        aria-busy={loading}
        {...props}
      >
        {loading && (
          <Loader2 className="spinner" size={16} aria-label="Loading" />
        )}
        {!loading && icon && <span aria-hidden="true">{icon}</span>}
        <span>{children}</span>
      </button>
    );
  }
);

Button.displayName = 'Button';

export default Button;

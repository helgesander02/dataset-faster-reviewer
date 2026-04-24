/**
 * Error handling utilities
 */

import { ApiError } from '@/services/api';
import { logger } from '@/utils/logger';

// ============================================================================
// Custom Error Classes
// ============================================================================

/**
 * Network error
 */
export class NetworkError extends Error {
  constructor(message: string = 'Network error occurred') {
    super(message);
    this.name = 'NetworkError';
  }
}

/**
 * API error
 */
export class APIError extends Error {
  public status?: number;
  public data?: unknown;

  constructor(message: string, status?: number, data?: unknown) {
    super(message);
    this.name = 'APIError';
    this.status = status;
    this.data = data;
  }
}

/**
 * Validation error
 */
export class ValidationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'ValidationError';
  }
}

// ============================================================================
// Error Type Guards
// ============================================================================

/**
 * Check if error is an API error
 */
export function isApiError(error: unknown): error is ApiError {
  return (
    typeof error === 'object' &&
    error !== null &&
    'message' in error &&
    typeof (error as ApiError).message === 'string'
  );
}

/**
 * Check if error is a network error
 */
export function isNetworkError(error: unknown): error is NetworkError {
  return error instanceof NetworkError || 
         (error instanceof Error && error.message.includes('Network'));
}

/**
 * Check if error is a validation error
 */
export function isValidationError(error: unknown): error is ValidationError {
  return error instanceof ValidationError;
}

/**
 * Check if error is a cancelled request
 */
export function isCancelledError(error: unknown): boolean {
  return (
    typeof error === 'object' &&
    error !== null &&
    'cancelled' in error &&
    (error as { cancelled?: boolean }).cancelled === true
  );
}

// ============================================================================
// Error Formatting
// ============================================================================

/**
 * Get user-friendly error message
 */
export function getErrorMessage(error: unknown): string {
  if (isCancelledError(error)) {
    return 'Request was cancelled';
  }

  if (isApiError(error)) {
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  if (typeof error === 'string') {
    return error;
  }

  return 'An unexpected error occurred';
}

/**
 * Format error for logging
 */
export function formatErrorForLogging(error: unknown): {
  name: string;
  message: string;
  stack?: string;
  status?: number;
  data?: unknown;
} {
  if (error instanceof APIError) {
    return {
      name: error.name,
      message: error.message,
      stack: error.stack,
      status: error.status,
      data: error.data,
    };
  }

  if (error instanceof Error) {
    return {
      name: error.name,
      message: error.message,
      stack: error.stack,
    };
  }

  return {
    name: 'Unknown',
    message: String(error),
  };
}

// ============================================================================
// Error Handling
// ============================================================================

/**
 * Safe error handler that won't throw
 */
export function safeErrorHandler(
  error: unknown,
  fallback: string = 'An error occurred'
): string {
  try {
    return getErrorMessage(error);
  } catch {
    return fallback;
  }
}

/**
 * Log error to console (dev only)
 */
export function logError(error: unknown, context?: string): void {
  const formatted = formatErrorForLogging(error);
  logger.error(
    `[Error${context ? ` in ${context}` : ''}]`,
    formatted
  );
}

/**
 * Handle error with retry logic
 */
export async function withRetry<T>(
  fn: () => Promise<T>,
  options: {
    retries?: number;
    delay?: number;
    shouldRetry?: (error: unknown) => boolean;
    onRetry?: (attempt: number, error: unknown) => void;
  } = {}
): Promise<T> {
  const {
    retries = 3,
    delay = 1000,
    shouldRetry = () => true,
    onRetry,
  } = options;

  let lastError: unknown;

  for (let attempt = 0; attempt < retries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error;

      // Don't retry if cancelled or if shouldRetry returns false
      if (isCancelledError(error) || !shouldRetry(error)) {
        throw error;
      }

      // Don't retry on last attempt
      if (attempt === retries - 1) {
        throw error;
      }

      // Call onRetry callback
      onRetry?.(attempt + 1, error);

      // Wait before retrying
      await new Promise(resolve => setTimeout(resolve, delay * (attempt + 1)));
    }
  }

  throw lastError;
}

// ============================================================================
// Error Boundary Helpers
// ============================================================================

/**
 * Check if error should show fallback UI
 */
export function shouldShowFallback(error: unknown): boolean {
  // Don't show fallback for cancelled requests
  if (isCancelledError(error)) {
    return false;
  }

  return true;
}

/**
 * Get fallback message for error boundary
 */
export function getFallbackMessage(error: unknown): string {
  if (isNetworkError(error)) {
    return 'Unable to connect to server. Please check your internet connection.';
  }

  if (isValidationError(error)) {
    return error.message;
  }

  if (isApiError(error)) {
    return error.message;
  }

  return 'Something went wrong. Please try again.';
}

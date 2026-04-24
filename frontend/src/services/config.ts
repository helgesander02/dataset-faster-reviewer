import { logger } from '@/utils/logger';

// ============================================================================
// Environment Variables
// ============================================================================

/**
 * Get environment variable with fallback
 */
function getEnvVar(key: string, fallback: string): string {
  if (typeof window !== 'undefined') {
    return process.env[key] || fallback;
  }
  return process.env[key] || fallback;
}

/**
 * Get numeric environment variable with fallback
 */
function getEnvNumber(key: string, fallback: number): number {
  const value = getEnvVar(key, '');
  const parsed = parseInt(value, 10);
  return isNaN(parsed) ? fallback : parsed;
}

/**
 * Get boolean environment variable with fallback
 */
function getEnvBoolean(key: string, fallback: boolean): boolean {
  const value = getEnvVar(key, '').toLowerCase();
  if (value === 'true' || value === '1') return true;
  if (value === 'false' || value === '0') return false;
  return fallback;
}

// ============================================================================
// API Configuration
// ============================================================================

/**
 * Base URL for API requests
 * 使用空字符串，因為 API 調用已經包含 /api 前綴
 * Traefik 會將 /api/* 路由到 backend
 */
export const API_BASE_URL = getEnvVar(
  'NEXT_PUBLIC_API_BASE_URL',
  ''
);

/**
 * API request timeout in milliseconds
 */
export const API_TIMEOUT = getEnvNumber('NEXT_PUBLIC_API_TIMEOUT', 30000);

// ============================================================================
// Pagination Configuration
// ============================================================================

/**
 * Number of datasets to display per page in sidebar
 */
export const DATASET_PER_PAGE = getEnvNumber('NEXT_PUBLIC_DATASET_PER_PAGE', 100);

/**
 * Number of images to load per page
 */
export const IMAGES_PER_PAGE = getEnvNumber('NEXT_PUBLIC_IMAGES_PER_PAGE', 9);

// ============================================================================
// Performance Configuration
// ============================================================================

/**
 * Job refresh interval in milliseconds
 */
export const JOB_REFRESH_INTERVAL = getEnvNumber(
  'NEXT_PUBLIC_JOB_REFRESH_INTERVAL',
  60000
);

/**
 * Number of adjacent pages to preload
 */
export const PRELOAD_RANGE = getEnvNumber('NEXT_PUBLIC_PRELOAD_RANGE', 1);

// ============================================================================
// Feature Flags
// ============================================================================

/**
 * Enable automatic retry for failed requests
 */
export const ENABLE_RETRY = getEnvBoolean('NEXT_PUBLIC_ENABLE_RETRY', true);

/**
 * Maximum number of retry attempts
 */
export const MAX_RETRIES = getEnvNumber('NEXT_PUBLIC_MAX_RETRIES', 2);

/**
 * Enable request cancellation
 */
export const ENABLE_REQUEST_CANCELLATION = getEnvBoolean(
  'NEXT_PUBLIC_ENABLE_REQUEST_CANCELLATION',
  true
);

// ============================================================================
// Development Configuration
// ============================================================================

/**
 * Check if running in development mode
 */
export const IS_DEV = process.env.NODE_ENV === 'development';

/**
 * Check if running in production mode
 */
export const IS_PROD = process.env.NODE_ENV === 'production';

/**
 * Enable debug logging
 */
export const DEBUG = IS_DEV && getEnvBoolean('NEXT_PUBLIC_DEBUG', false);

// ============================================================================
// Validation
// ============================================================================

if (typeof window === 'undefined' && API_BASE_URL && !API_BASE_URL.startsWith('http') && !API_BASE_URL.startsWith('/')) {
  logger.warn('⚠️  API_BASE_URL should start with http://, https://, or /');
}

if (IMAGES_PER_PAGE <= 0) {
  throw new Error('IMAGES_PER_PAGE must be greater than 0');
}

if (DATASET_PER_PAGE <= 0) {
  throw new Error('DATASET_PER_PAGE must be greater than 0');
}

// ============================================================================
// Export all configuration
// ============================================================================

export const config = {
  api: {
    baseURL: API_BASE_URL,
    timeout: API_TIMEOUT,
  },
  pagination: {
    datasetsPerPage: DATASET_PER_PAGE,
    imagesPerPage: IMAGES_PER_PAGE,
  },
  performance: {
    jobRefreshInterval: JOB_REFRESH_INTERVAL,
    preloadRange: PRELOAD_RANGE,
  },
  features: {
    enableRetry: ENABLE_RETRY,
    maxRetries: MAX_RETRIES,
    enableRequestCancellation: ENABLE_REQUEST_CANCELLATION,
  },
  env: {
    isDev: IS_DEV,
    isProd: IS_PROD,
    debug: DEBUG,
  },
} as const;

export default config;

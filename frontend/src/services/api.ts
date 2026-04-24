import axios, { AxiosError, AxiosResponse } from 'axios';
import { API_BASE_URL } from '@/services/config';
import { logger } from '@/utils/logger';

// ============================================================================
// Types
// ============================================================================

export interface JobsResponse {
  job_names: string[];
}

export interface DatasetsResponse {
  dataset_names: string[];
}

export interface PageItem {
  item_dataset_name: string;
  item_image_set: string[];
}

export interface GetAllPagesResponse {
  total_pages: number;
  pages: PageItem[];
}

export interface JobMetadataResponse {
  job_name: string;
  total_pages: number;
  dataset_names: string[];
}

export interface ImageSetResponse {
  image_name: string[];
  image_path: string[];
}

export interface Base64ImageSetResponse {
  image_path: string[];
  base64_image: string[];
}

export interface ImageToSave {
  job: string;
  dataset: string;
  imageName: string;
  imagePath: string;
}

export interface SavePendingReviewPayload {
  images: ImageToSave[];
}

export interface PendingReviewItem {
  item_job_name: string;
  item_dataset_name: string;
  item_image_name: string;
  item_image_path: string;
}

export interface PendingReviewResponse {
  items: PendingReviewItem[];
}

export interface ApiError {
  message: string;
  status?: number;
  data?: unknown;
}

// ============================================================================
// API Configuration
// ============================================================================

const api = axios.create({
  baseURL: API_BASE_URL?.endsWith('/') ? API_BASE_URL : `${API_BASE_URL}/`,
  timeout: 30000, // 30 seconds
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request tracking for cancellation (disabled by default to prevent issues)
// Can be enabled via config if needed
const pendingRequests = new Map<string, AbortController>();

// ============================================================================
// Interceptors
// ============================================================================

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Log requests in development
    if (process.env.NODE_ENV === 'development') {
      logger.log(`[API] ${config.method?.toUpperCase()} ${config.url}`);
    }
    
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor
api.interceptors.response.use(
  (response: AxiosResponse) => {
    // Log responses in development
    if (process.env.NODE_ENV === 'development') {
      logger.log(`[API] Response from ${response.config.url}:`, response.status);
    }
    
    return response;
  },
  async (error: AxiosError) => {
    // Don't retry if request was cancelled
    if (axios.isCancel(error)) {
      return Promise.reject({ message: 'Request cancelled', cancelled: true });
    }

    // Type guard to ensure error is AxiosError
    const axiosError: AxiosError = error;

    // Handle different error types
    const apiError: ApiError = {
      message: 'An unexpected error occurred',
      status: axiosError.response?.status,
      data: axiosError.response?.data,
    };

    if (axiosError.response) {
      // Server responded with error status
      switch (axiosError.response.status) {
        case 400:
          apiError.message = 'Invalid request parameters';
          break;
        case 404:
          apiError.message = 'Resource not found';
          break;
        case 500:
          apiError.message = 'Server error occurred';
          break;
        case 503:
          apiError.message = 'Service temporarily unavailable';
          break;
        default:
          apiError.message = `Request failed with status ${axiosError.response.status}`;
      }
    } else if (axiosError.request) {
      // Request made but no response received
      apiError.message = 'No response from server. Please check your connection.';
    } else {
      // Request setup error
      apiError.message = axiosError.message || 'Failed to make request';
    }

    // Log error in development only
    if (process.env.NODE_ENV === 'development') {
      logger.error('[API Error]', {
        url: axiosError.config?.url,
        method: axiosError.config?.method,
        params: axiosError.config?.params,
        error: apiError,
      });
    }

    return Promise.reject(apiError);
  }
);

// ============================================================================
// Retry Logic
// ============================================================================

async function withRetry<T>(
  requestFn: () => Promise<AxiosResponse<T>>,
  retries: number = 2,
  delay: number = 1000
): Promise<T> {
  let lastError: unknown;

  for (let attempt = 0; attempt <= retries; attempt++) {
    try {
      const response = await requestFn();
      return response.data;
    } catch (error) {
      lastError = error;
      
      // Don't retry if cancelled
      if ((error as { cancelled?: boolean }).cancelled) {
        throw error;
      }

      // Don't retry client errors (4xx)
      const status = (error as ApiError).status;
      if (status && status >= 400 && status < 500) {
        throw error;
      }

      // Wait before retry (except on last attempt)
      if (attempt < retries) {
        await new Promise(resolve => setTimeout(resolve, delay * (attempt + 1)));
      }
    }
  }

  throw lastError;
}

// ============================================================================
// API Functions
// ============================================================================

export const fetchJobs = async (): Promise<JobsResponse> => {
  return withRetry<JobsResponse>(() => api.get('/api/getJobs'));
};

export const fetchDatasets = async (job: string): Promise<DatasetsResponse> => {
  logger.log('[API] fetchDatasets called for job:', job);
  
  if (!job?.trim()) {
    throw new Error('The "job" parameter is required to fetch datasets.');
  }

  logger.log('[API] Calling GET /api/getAllPages with params:', { job });
  const response = await withRetry<GetAllPagesResponse>(() => 
    api.get('/api/getAllPages', { params: { job } })
  );

  logger.log('[API] getAllPages response:', { 
    totalPages: response.pages?.length,
    samplePage: response.pages?.[0] 
  });

  const datasetNames = response.pages
    .map(page => page.item_dataset_name)
    .filter(name => name && name.trim() !== '');

  logger.log('[API] Extracted dataset names:', datasetNames);

  return { dataset_names: datasetNames };
};

export const fetchBase64Images = async (
  job: string,
  pageIndex: number
): Promise<Base64ImageSetResponse> => {
  logger.log('[API] fetchBase64Images called:', { job, pageIndex });
  
  if (!job?.trim()) {
    throw new Error('The "job" parameter is required.');
  }
  
  if (pageIndex < 0) {
    throw new Error('The "pageIndex" must be non-negative.');
  }

  logger.log('[API] Calling GET /api/getBase64ImageSet with params:', { job, pageIndex });
  const result = await withRetry<Base64ImageSetResponse>(() =>
    api.get('/api/getBase64ImageSet', { params: { job, pageIndex } })
  );
  
  logger.log('[API] getBase64ImageSet response:', {
    job,
    pageIndex,
    imageCount: result.base64_image?.length || 0,
    pathCount: result.image_path?.length || 0
  });
  
  return result;
};

export const fetchImages = async (
  job: string,
  pageIndex: number
): Promise<ImageSetResponse> => {
  logger.log('[API] fetchImages called:', { job, pageIndex });
  
  if (!job?.trim()) {
    throw new Error('The "job" parameter is required.');
  }
  
  if (pageIndex < 0) {
    throw new Error('The "pageIndex" must be non-negative.');
  }

  logger.log('[API] Calling GET /api/getImageSet with params:', { job, pageIndex });
  const result = await withRetry<ImageSetResponse>(() =>
    api.get('/api/getImageSet', { params: { job, pageIndex } })
  );
  
  logger.log('[API] getImageSet response:', {
    job,
    pageIndex,
    imageCount: result.image_name?.length || 0,
    pathCount: result.image_path?.length || 0
  });
  
  return result;
};

export const updateALLPages = async (
  job: string,
  image_per_page: number
): Promise<GetAllPagesResponse> => {
  logger.log('[API] updateALLPages called:', { job, image_per_page });
  
  if (!job?.trim()) {
    throw new Error('The "job" parameter is required to update all pages.');
  }

  if (!image_per_page || image_per_page <= 0) {
    throw new Error('The "image_per_page" parameter is required and must be greater than 0.');
  }

  logger.log('[API] Calling POST /api/setAllPages with body:', { job, image_per_page });
  const result = await withRetry<GetAllPagesResponse>(() =>
    api.post('/api/setAllPages', { job, image_per_page })
  );
  
  return result;
};

export const fetchJobMetadata = async (job: string): Promise<JobMetadataResponse> => {
  logger.log('[API] fetchJobMetadata called for job:', job);
  
  if (!job?.trim()) {
    throw new Error('The "job" parameter is required to fetch job metadata.');
  }

  logger.log('[API] Calling GET /api/getJobMetadata with params:', { job });
  const result = await withRetry<JobMetadataResponse>(() =>
    api.get('/api/getJobMetadata', { params: { job } })
  );
  
  logger.log('[API] fetchJobMetadata response:', {
    job,
    totalPages: result.total_pages,
    datasetsCount: result.dataset_names?.length || 0,
    datasets: result.dataset_names
  });
  
  return result;
};

export const fetchALLPages = async (job: string): Promise<GetAllPagesResponse> => {
  logger.log('[API] fetchALLPages called for job:', job);
  
  if (!job?.trim()) {
    throw new Error('The "job" parameter is required to fetch all pages.');
  }

  logger.log('[API] Calling GET /api/getAllPages with params:', { job });
  const result = await withRetry<GetAllPagesResponse>(() =>
    api.get('/api/getAllPages', { params: { job } })
  );
  
  logger.log('[API] fetchALLPages response:', {
    job,
    totalPages: result.total_pages,
    pagesCount: result.pages?.length || 0
  });
  
  return result;
};

export const savePendingReview = async (
  data: SavePendingReviewPayload
): Promise<{ success: boolean }> => {
  if (!data.images) {
    data.images = [];
  }

  return withRetry<{ success: boolean }>(() =>
    api.post('/api/savePendingReview', data.images)
  );
};

export const getPendingReview = async (
  flatten: boolean = false
): Promise<PendingReviewResponse> => {
  return withRetry<PendingReviewResponse>(() =>
    api.get('/api/getPendingReview', { params: { flatten } })
  );
};

/**
 * Fetch paginated review images (paths only, no base64)
 * Used for progressive loading in review modal
 */
export const fetchPendingReviewImages = async (
  page: number = 0,
  limit: number = 9
): Promise<{
  image_path: string[];
  base64_image: string[];
  total_items: number;
  page: number;
}> => {
  logger.log('[API] fetchPendingReviewImages called:', { page, limit });
  
  if (page < 0) {
    throw new Error('The "page" parameter must be non-negative.');
  }
  
  if (limit <= 0) {
    throw new Error('The "limit" parameter must be greater than 0.');
  }

  // Note: This endpoint may need to be created on the backend
  // For now, we'll use getPendingReview and paginate on frontend
  // Backend optimization can be added later
  const result = await withRetry<PendingReviewResponse>(() =>
    api.get('/api/getPendingReview', { params: { flatten: true } })
  );
  
  const startIndex = page * limit;
  const endIndex = startIndex + limit;
  const paginatedItems = result.items.slice(startIndex, endIndex);
  
  return {
    image_path: paginatedItems.map(item => item.item_image_path),
    base64_image: [], // No base64 images in this response
    total_items: result.items.length,
    page: page
  };
};

export const deleteSelectedImages = async (
  images: ImageToSave[]
): Promise<{ success: boolean; deleted_count: number }> => {
  if (!images || images.length === 0) {
    throw new Error('At least one image is required to delete.');
  }

  return withRetry<{ success: boolean; deleted_count: number }>(() =>
    api.post('/api/deleteSelectedImages', images)
  );
};

// ============================================================================
// Request Cancellation Utilities
// ============================================================================

/**
 * Cancel all pending requests
 */
export const cancelAllRequests = (): void => {
  pendingRequests.forEach(controller => controller.abort());
  pendingRequests.clear();
};

/**
 * Cancel requests matching a pattern
 */
export const cancelRequestsByPattern = (pattern: string): void => {
  pendingRequests.forEach((controller, key) => {
    if (key.includes(pattern)) {
      controller.abort();
      pendingRequests.delete(key);
    }
  });
};

// ============================================================================
// Export API Instance (for custom requests if needed)
// ============================================================================

export default api;

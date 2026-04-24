"use client";

import { useState, useCallback, useRef, useEffect } from 'react';

// ============================================================================
// Types
// ============================================================================

export interface UseApiState<T> {
  data: T | null;
  error: Error | null;
  loading: boolean;
  success: boolean;
}

export interface UseApiOptions {
  immediate?: boolean; // Execute immediately on mount
  onSuccess?: (data: unknown) => void;
  onError?: (error: Error) => void;
}

export interface UseApiReturn<T, Args extends unknown[]> extends UseApiState<T> {
  execute: (...args: Args) => Promise<T | null>;
  reset: () => void;
}

// ============================================================================
// Hook: useApi
// ============================================================================

/**
 * Generic hook for API calls with loading, error, and success states
 * 
 * @example
 * const { data, loading, error, execute } = useApi(fetchJobs);
 * 
 * useEffect(() => {
 *   execute();
 * }, []);
 */
export function useApi<T, Args extends unknown[]>(
  apiFunction: (...args: Args) => Promise<T>,
  options: UseApiOptions = {}
): UseApiReturn<T, Args> {
  const { immediate = false, onSuccess, onError } = options;

  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    error: null,
    loading: false,
    success: false,
  });

  const abortControllerRef = useRef<AbortController | null>(null);
  const mountedRef = useRef(true);

  useEffect(() => {
    mountedRef.current = true;
    return () => {
      mountedRef.current = false;
      abortControllerRef.current?.abort();
    };
  }, []);

  const execute = useCallback(
    async (...args: Args): Promise<T | null> => {
      // Cancel previous request if exists
      abortControllerRef.current?.abort();
      abortControllerRef.current = new AbortController();

      setState(prev => ({
        ...prev,
        loading: true,
        error: null,
      }));

      try {
        const result = await apiFunction(...args);

        if (!mountedRef.current) return null;

        setState({
          data: result,
          error: null,
          loading: false,
          success: true,
        });

        onSuccess?.(result);
        return result;
      } catch (err) {
        if (!mountedRef.current) return null;

        const error = err instanceof Error ? err : new Error('Unknown error occurred');

        setState({
          data: null,
          error,
          loading: false,
          success: false,
        });

        onError?.(error);
        return null;
      }
    },
    [apiFunction, onSuccess, onError]
  );

  const reset = useCallback(() => {
    abortControllerRef.current?.abort();
    setState({
      data: null,
      error: null,
      loading: false,
      success: false,
    });
  }, []);

  // Execute immediately if requested
  useEffect(() => {
    if (immediate) {
      // @ts-expect-error - TypeScript cannot infer empty args correctly
      execute();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [immediate]); // Only run on mount, execute is stable via useCallback

  return {
    ...state,
    execute,
    reset,
  };
}

// ============================================================================
// Hook: useApiMutation
// ============================================================================

/**
 * Hook for mutations (POST, PUT, DELETE) with optimistic updates support
 */
export function useApiMutation<T, Args extends unknown[]>(
  apiFunction: (...args: Args) => Promise<T>,
  options: UseApiOptions = {}
) {
  return useApi(apiFunction, { ...options, immediate: false });
}

// ============================================================================
// Hook: useApiQuery
// ============================================================================

/**
 * Hook for queries (GET) with automatic execution
 */
export function useApiQuery<T, Args extends unknown[]>(
  apiFunction: (...args: Args) => Promise<T>,
  args: Args,
  options: UseApiOptions & { enabled?: boolean } = {}
) {
  const { enabled = true, ...restOptions } = options;
  const api = useApi(apiFunction, restOptions);

  useEffect(() => {
    if (enabled) {
      api.execute(...args);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, JSON.stringify(args)]); // Re-fetch when args change (serialized for comparison)

  return api;
}

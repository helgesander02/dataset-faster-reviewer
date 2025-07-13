"use client";

import { useEffect, useRef } from 'react';

export function useIntersectionObserver(
  callback: () => void,
  loading: boolean,
  rootMargin: string = '1px'
) {
  const observerRef = useRef<IntersectionObserver | null>(null);
  const loadingRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!loadingRef.current) return;

    observerRef.current = new IntersectionObserver(
      (entries) => {
        const [entry] = entries;
        if (entry.isIntersecting && !loading) {
          callback();
        }
      }, { rootMargin }
    );

    if (loadingRef.current) {
      observerRef.current.observe(loadingRef.current);
    }

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [loading, callback, rootMargin]);

  return loadingRef;
}

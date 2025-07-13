"use client";

import { useCallback, useRef } from 'react';

export function usePageObserver() {
  const pageObserversRef = useRef<Map<number, IntersectionObserver>>(new Map());

  const setupPageObserver = useCallback((
    pageIndex: number, element: HTMLElement, 
    onIntersect: (pageIndex: number) => void
  ) => {

    if (pageObserversRef.current.has(pageIndex)) {
      return;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && entry.intersectionRatio > 0.5) {
            onIntersect(pageIndex);
          }
        });
      },
      {
        root: null,
        rootMargin: '0px',
        threshold: 0.5
      }
    );

    observer.observe(element);
    pageObserversRef.current.set(pageIndex, observer);
  }, []);

  const cleanupObservers = useCallback(() => {
    pageObserversRef.current.forEach((observer) => {
      observer.disconnect();
    });
    pageObserversRef.current.clear();
  }, []);

  const registerPageElement = useCallback((
    pageIndex: number, element: HTMLElement | null, 
    onIntersect: (pageIndex: number) => void
  ) => {
    if (element) {
      setupPageObserver(pageIndex, element, onIntersect);
    }
  }, [setupPageObserver]);

  return {
    registerPageElement,
    cleanupObservers
  };
}

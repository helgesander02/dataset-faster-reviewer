"use client";

import React from 'react';
import { LoadingTriggerProps } from '@/types/HomeImageGrid';

export default function LoadingTrigger({ 
  isLoading, 
  loadingMessage = "Loading more images..." 
}: LoadingTriggerProps) {
  return (
    <div className="loading-trigger">
      {isLoading && (
        <div className="loading-state">
          <div className="loading-state-content">
            <div className="loading-message">{loadingMessage}</div>
          </div>
        </div>
      )}
    </div>
  );
}

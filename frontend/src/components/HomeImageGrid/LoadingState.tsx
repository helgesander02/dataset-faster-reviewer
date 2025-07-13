"use client";

import React from 'react';
import { LoadingStateProps } from '@/types/HomeImageGrid';

export default function LoadingState({ 
    message = "Loading images..." 
}: LoadingStateProps) {
    
  return (
    <div className="loading-state">
      <div className="loading-state-content">
        <div className="loading-message">{message}</div>
      </div>
    </div>
  );
}

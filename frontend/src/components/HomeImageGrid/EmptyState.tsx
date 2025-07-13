"use client";

import React from 'react';
import { EmptyStateProps } from '@/types/HomeImageGrid';

export default function EmptyState({ 
  title = "Welcome to Image Verify Viewer", 
  message = "Please select a job from the left sidebar to start reviewing images." 
}: EmptyStateProps) {
  return (
    <div className="empty-state">
      <div className="empty-state-content">
        <h2>{title}</h2>
        <p>{message}</p>
      </div>
    </div>
  );
}

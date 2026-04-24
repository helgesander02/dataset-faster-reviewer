"use client";

import React from 'react';
import { EmptyStateProps } from '@/types/HomeImageGrid';

export default function EmptyState({ 
  title = "Welcome to Image Verify Viewer", 
  message = "Please select a job from the left sidebar to start reviewing images." 
}: EmptyStateProps) {
  return (
    <div className="flex items-center justify-center h-screen bg-gray-50">
      <div className="text-center p-8">
        <h2 className="text-gray-900 text-2xl font-bold mb-2">{title}</h2>
        <p className="text-gray-600 text-base m-0">{message}</p>
      </div>
    </div>
  );
}

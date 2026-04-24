"use client";

import React from 'react';
import { LoadingStateProps } from '@/types/HomeImageGrid';

export default function LoadingState({ 
    message = "Loading images..." 
}: LoadingStateProps) {
    
  return (
    <div className="flex h-[200px] items-center justify-center bg-gray-50">
      <div className="animate-pulse">
        <div className="text-gray-600 text-base font-medium">{message}</div>
      </div>
    </div>
  );
}

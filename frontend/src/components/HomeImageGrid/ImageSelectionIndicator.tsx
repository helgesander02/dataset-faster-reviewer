"use client";

import React from 'react';
import { SelectionIndicatorProps } from '@/types/HomeImageGrid';

export default function SelectionIndicator({ 
  className = "selection-indicator"
}: SelectionIndicatorProps) {
  
  return (
    <div className={className}>
      <svg 
        className="checkmark" 
        viewBox="0 0 24 24" 
        fill="none" 
        stroke="currentColor" 
        strokeWidth="2"
      >
        <polyline points="20,6 9,17 4,12"></polyline>
      </svg>
    </div>
  );
}

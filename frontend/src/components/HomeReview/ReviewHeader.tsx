"use client";

import { X, ImageIcon, Loader2 } from 'lucide-react';
import { ReviewHeaderProps } from '@/types/HomeReview';

export function ReviewHeader({ 
  saving, onClose 
}: ReviewHeaderProps) {

  return (
    <div className="home-review-header">
      <h2 className="home-review-title">
        <ImageIcon size={20} /> Select Picture
        {saving && <Loader2 size={16} className="loading-spin" />}
      </h2>
      <button 
        onClick={onClose}
        className="home-review-close-btn"
        aria-label="closed"
      >
        <X size={20} />
      </button>
    </div>
  );
}

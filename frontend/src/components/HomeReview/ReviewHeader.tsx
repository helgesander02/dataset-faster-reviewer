"use client";

import { X, ImageIcon, Loader2 } from 'lucide-react';
import { ReviewHeaderProps } from '@/types/HomeReview';

export function ReviewHeader({ 
  saving, onClose 
}: ReviewHeaderProps) {

  return (
    <div className="flex items-center justify-between p-6 border-b border-gray-200">
      <h2 className="flex items-center gap-2 text-xl font-semibold text-gray-900 m-0">
        <ImageIcon size={20} /> Select Picture
        {saving && <Loader2 size={16} className="animate-spin" />}
      </h2>
      <button 
        onClick={onClose}
        className="bg-transparent border-none cursor-pointer p-2 rounded text-gray-500 transition-all duration-200 hover:bg-gray-100 hover:text-gray-700"
        aria-label="closed"
      >
        <X size={20} />
      </button>
    </div>
  );
}

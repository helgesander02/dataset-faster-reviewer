"use client";

import { Loader2 } from 'lucide-react';

export function LoadingState() {
  return (
    <div className="home-review-loading">
      <Loader2 size={24} className="loading-spin" />
      <span>loading...</span>
    </div>
  );
}

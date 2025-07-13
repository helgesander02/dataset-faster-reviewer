"use client";

import { ImageIcon } from 'lucide-react';

export function EmptyState() {
  return (
    <div className="home-review-empty">
      <ImageIcon size={48} />
      <p>There are currently no pictures</p>
    </div>
  );
}

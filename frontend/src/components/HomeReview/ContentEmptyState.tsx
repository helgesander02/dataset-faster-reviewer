"use client";

import { ImageIcon } from 'lucide-react';

export function EmptyState() {
  return (
    <div className="flex flex-col items-center gap-4 p-12 text-gray-500 text-center">
      <ImageIcon size={48} />
      <p className="m-0 text-lg">There are currently no pictures</p>
    </div>
  );
}

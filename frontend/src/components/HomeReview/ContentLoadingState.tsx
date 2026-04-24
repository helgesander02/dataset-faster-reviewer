"use client";

import { Loader2 } from 'lucide-react';

export function LoadingState() {
  return (
    <div className="flex flex-col items-center justify-center gap-4 p-12 text-gray-500">
      <Loader2 size={24} className="animate-spin" />
      <span>loading...</span>
    </div>
  );
}

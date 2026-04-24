"use client";

interface ErrorStateProps {
  error: string;
  onRetry: () => void;
}

export function ErrorState({ error, onRetry }: ErrorStateProps) {
  return (
    <div className="flex flex-col items-center gap-4 p-12 text-center">
      <p className="text-red-500 m-0">{error}</p>
      <button 
        onClick={onRetry}
        className="bg-blue-500 text-white border-none py-2 px-4 rounded-md cursor-pointer transition-colors duration-200 hover:bg-blue-600"
      > Reload 
      </button>
    </div>
  );
}

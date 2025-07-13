"use client";

interface ErrorStateProps {
  error: string;
  onRetry: () => void;
}

export function ErrorState({ error, onRetry }: ErrorStateProps) {
  return (
    <div className="home-review-error">
      <p>{error}</p>
      <button 
        onClick={onRetry}
        className="home-review-retry-btn"
      > Reload 
      </button>
    </div>
  );
}

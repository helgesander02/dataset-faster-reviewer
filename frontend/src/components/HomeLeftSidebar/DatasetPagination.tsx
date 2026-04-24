"use client";

import React, { useState, useCallback } from 'react';
import { PaginationProps } from '@/types/HomeLeftSidebar';

export default function Pagination({ 
  currentPagenation, totalDatasets, datasetsPerPage, 
  onPrevious, onNext, onGoToPage 
}: PaginationProps) {

  const [inputPage, setInputPage] = useState<string>('');
  const [isEditing, setIsEditing] = useState(false);

  const isFirstPage = currentPagenation === 0;
  const totalPages = Math.ceil(totalDatasets / datasetsPerPage);
  const isLastPage = currentPagenation >= totalPages - 1;
  const currentPage = currentPagenation + 1; // Display as 1-indexed

  const handlePageClick = useCallback(() => {
    setIsEditing(true);
    setInputPage(currentPage.toString());
  }, [currentPage]);

  const handleInputChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    // Only allow numbers
    if (value === '' || /^\d+$/.test(value)) {
      setInputPage(value);
    }
  }, []);

  const handleInputBlur = useCallback(() => {
    setIsEditing(false);
    if (inputPage && onGoToPage) {
      const page = parseInt(inputPage, 10);
      if (page >= 1 && page <= totalPages) {
        onGoToPage(page - 1); // Convert back to 0-indexed
      }
    }
    setInputPage('');
  }, [inputPage, onGoToPage, totalPages]);

  const handleInputKeyDown = useCallback((e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      handleInputBlur();
    } else if (e.key === 'Escape') {
      setIsEditing(false);
      setInputPage('');
    }
  }, [handleInputBlur]);
  
  return (
    <div className="flex gap-1.5 mt-2 items-center">
      <button 
        className="flex-1 bg-blue-500 text-white py-1.5 px-3 rounded border-none cursor-pointer text-sm font-medium transition-colors duration-200 hover:bg-blue-600 disabled:bg-gray-300 disabled:text-gray-400 disabled:cursor-not-allowed" 
        onClick={onPrevious} 
        disabled={isFirstPage}
      >
        Previous
      </button>
      
      <div className="flex items-center gap-1 px-1 flex-shrink-0">
        {isEditing ? (
          <input
            type="text"
            value={inputPage}
            onChange={handleInputChange}
            onBlur={handleInputBlur}
            onKeyDown={handleInputKeyDown}
            className="w-12 px-1 py-1 text-center border border-blue-500 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
            autoFocus
            maxLength={3}
          />
        ) : (
          <button
            onClick={handlePageClick}
            className="px-1 py-1 text-blue-600 hover:bg-blue-50 rounded transition-colors cursor-pointer font-medium text-sm min-w-[2rem]"
            title="點擊編輯頁碼"
          >
            {currentPage}
          </button>
        )}
        <span className="text-gray-600 text-sm whitespace-nowrap">/ {totalPages}</span>
      </div>

      <button 
        className="flex-1 bg-blue-500 text-white py-1.5 px-3 rounded border-none cursor-pointer text-sm font-medium transition-colors duration-200 hover:bg-blue-600 disabled:bg-gray-300 disabled:text-gray-400 disabled:cursor-not-allowed" 
        onClick={onNext} 
        disabled={isLastPage}
      >
        Next
      </button>
    </div>
  );
}

"use client";

import React from 'react';
import { JobSelectProps } from '@/types/HomeLeftSidebar';

export default function JobSelect({ 
  currentJobList, selectedJob, loading,
  onJobSelect 
}: JobSelectProps) {
  
  return (
    <div className="mb-4">
      <select 
        className="w-full p-2 border border-gray-300 rounded bg-white text-gray-700 text-sm focus:outline-none focus:border-blue-500 focus:ring-[3px] focus:ring-blue-500/10 disabled:bg-gray-50 disabled:text-gray-400 disabled:cursor-not-allowed"
        value={selectedJob}
        onChange={(e) => onJobSelect(e.target.value)}
        disabled={loading}
      >
        <option value="" disabled>Select a Job</option>
        {currentJobList.map((job, index) => (
          <option key={job} value={job}>
            {index + 1}. {job}
          </option>
        ))}
      </select>
    </div>
  );
}

"use client";

import React from 'react';
import { JobSelectProps } from '@/types/HomeLeftSidebar';

export default function JobSelect({ 
  currentJobList, selectedJob, loading,
  onJobSelect 
}: JobSelectProps) {
  
  return (
    <div className="job-select-container">
      <select 
        className="job-select"
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

import { StatusProps } from '@/types/HomeLeftSidebar';

export default function Status({ 
  selectedJob, selectedDataset 
}: StatusProps) {
  
  return (
    <div className="status">
      {selectedJob && selectedDataset ? (
        <div className="ok">
          Selected: {selectedJob} / {selectedDataset}
        </div>
      ) : (
        <div className="warn">
          {selectedJob ? 'Please select a dataset' : 'Please select a job'}
        </div>
      )}
    </div>
  );
}

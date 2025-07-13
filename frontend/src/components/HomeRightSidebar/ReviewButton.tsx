import { Eye } from 'lucide-react';
import { ReviewButtonProps } from '@/types/HomeRightSidebar';

export default function ReviewButton({ 
  loading, onReview 
}: ReviewButtonProps) {
  
  return (
    <div className="action-buttons">
      <button 
        onClick={onReview}
        className="btn btn-review"
        disabled={loading}
      >
        <Eye size={16} />
        <span>Review</span>
      </button>
    </div>
  );
}

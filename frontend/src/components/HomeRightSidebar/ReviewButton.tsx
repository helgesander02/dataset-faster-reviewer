import { Eye } from 'lucide-react';
import { ReviewButtonProps } from '@/types/HomeRightSidebar';

export default function ReviewButton({ 
  loading, onReview 
}: ReviewButtonProps) {
  
  return (
    <div className="flex flex-col gap-2 mb-4">
      <button 
        onClick={onReview}
        className="w-full py-2 px-4 rounded-md flex items-center justify-center gap-2 cursor-pointer transition-colors duration-200 bg-blue-500 text-white hover:bg-blue-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
        disabled={loading}
      >
        <Eye size={16} />
        <span>Review</span>
      </button>
    </div>
  );
}

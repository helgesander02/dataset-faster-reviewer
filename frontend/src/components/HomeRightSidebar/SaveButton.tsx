import { Save, Clock } from 'lucide-react';
import { SaveButtonProps } from '@/types/HomeRightSidebar';

export default function SaveButton({  
  loading, saveSuccess, cachedImages, disabled, 
  onSave 
}: SaveButtonProps) {

  return (
    <div>
      {saveSuccess && (
        <div className="flex items-center gap-1 text-sm text-emerald-600 mb-2">
          <Clock size={12} /> Saved successfully!
        </div>
      )}
      
      <button 
        onClick={onSave}
        disabled={disabled}
        className="w-full py-2 px-4 rounded-md flex items-center justify-center gap-2 cursor-pointer transition-colors duration-200 bg-emerald-500 text-white hover:bg-emerald-600 disabled:bg-gray-300 disabled:cursor-not-allowed"
      >
        {loading ? (
          <>
            <div className="border-2 border-transparent border-b-white rounded-full w-4 h-4 animate-spin"></div>
            <span>Saving...</span>
          </>
        ) : (
          <>
            <Save size={16} />
            <span>Save ({cachedImages.length})</span>
          </>
        )}
      </button>
    </div>
  );
}

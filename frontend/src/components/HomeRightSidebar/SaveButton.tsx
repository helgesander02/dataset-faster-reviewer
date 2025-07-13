import { Save, Clock } from 'lucide-react';
import { SaveButtonProps } from '@/types/HomeRightSidebar';

export default function SaveButton({  
  loading, saveSuccess, cachedImages, disabled, 
  onSave 
}: SaveButtonProps) {

  return (
    <div>
      {saveSuccess && (
        <div className="save-success">
          <Clock size={12} /> Saved successfully!
        </div>
      )}
      
      <button 
        onClick={onSave}
        disabled={disabled}
        className="btn btn-save"
      >
        {loading ? (
          <>
            <div className="loading-spinner"></div>
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

"use client";

import { Save, Loader2 } from 'lucide-react';

interface ReviewActionsProps {
  selectedCount: number;
  totalCount: number;
  saving: boolean;
  onSelectAll: () => void;
  onDeselectAll: () => void;
  onSave: () => void;
}

export function ReviewActions({ 
  selectedCount, 
  totalCount, 
  saving, 
  onSelectAll, 
  onDeselectAll, 
  onSave 
}: ReviewActionsProps) {
  return (
    <div className="home-review-footer">
      <div className="home-review-total">
        Selected: {selectedCount} / {totalCount}
      </div>
      <div className="home-review-actions">
        <button 
          className="home-review-action-btn secondary"
          onClick={onDeselectAll}
          disabled={saving}
        >
          Cancel All
        </button>
        <button 
          className="home-review-action-btn secondary"
          onClick={onSelectAll}
          disabled={saving}
        >
          Select ALL
        </button>
        <button 
          className="home-review-action-btn primary save-btn"
          onClick={onSave}
          disabled={saving}
        >
          {saving ? (
            <>
              <Loader2 size={16} className="loading-spin" />
              Saving...
            </>
          ) : (
            <>
              <Save size={16} />
              Save
            </>
          )}
        </button>
      </div>
    </div>
  );
}

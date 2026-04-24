"use client";

import { Save, Loader2, Trash2 } from 'lucide-react';
import { ReviewActionsProps } from '@/types/HomeReview';

const BUTTON_BASE_CLASS = "flex items-center gap-2 py-2 px-4 rounded-md text-sm font-medium cursor-pointer transition-all duration-200 border-none disabled:opacity-50 disabled:cursor-not-allowed md:flex-1 md:min-w-[120px]";
const SECONDARY_BUTTON_CLASS = `${BUTTON_BASE_CLASS} bg-gray-100 text-gray-700 border border-gray-300 hover:bg-gray-200`;
const DANGER_BUTTON_CLASS = `${BUTTON_BASE_CLASS} bg-red-500 text-white hover:bg-red-600 disabled:bg-gray-400`;
const PRIMARY_BUTTON_CLASS = `${BUTTON_BASE_CLASS} bg-emerald-500 text-white hover:bg-emerald-600 disabled:bg-gray-400`;

export function ReviewActions({ 
  selectedCount, totalCount, saving, deleting,
  onSelectAll, onDeselectAll, onSave, onDelete
}: ReviewActionsProps) {
  const isDisabled = saving || deleting;

  return (
    <div className="border-t border-gray-200 p-4 bg-gray-50 flex justify-between items-center md:flex-col md:gap-4 md:items-stretch">
      <div className="font-semibold text-gray-900">
        Selected: {selectedCount} / {totalCount}
      </div>
      <div className="flex gap-3 md:justify-center md:flex-wrap">
        <button 
          className={SECONDARY_BUTTON_CLASS}
          onClick={onDeselectAll}
          disabled={isDisabled}
        >
          Cancel All
        </button>
        <button 
          className={SECONDARY_BUTTON_CLASS}
          onClick={onSelectAll}
          disabled={isDisabled}
        >
          Select ALL
        </button>
        <button 
          className={DANGER_BUTTON_CLASS}
          onClick={onDelete}
          disabled={isDisabled || selectedCount === 0}
        >
          {deleting ? (
            <>
              <Loader2 size={16} className="animate-spin" />
              Deleting...
            </>
          ) : (
            <>
              <Trash2 size={16} />
              Delete
            </>
          )}
        </button>
        <button 
          className={PRIMARY_BUTTON_CLASS}
          onClick={onSave}
          disabled={isDisabled}
        >
          {saving ? (
            <>
              <Loader2 size={16} className="animate-spin" />
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

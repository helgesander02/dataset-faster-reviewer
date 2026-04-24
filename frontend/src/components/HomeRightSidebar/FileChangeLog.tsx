import { FileText } from 'lucide-react';
import { FileChangeLogProps } from '@/types/HomeRightSidebar';

export default function FileChangeLog({ 
  groupedImages, cachedImages 
}: FileChangeLogProps) {

  return (
    <div className="flex-1 border border-dashed border-gray-300 bg-white rounded-md mb-4 p-2 overflow-auto">
      <div className="text-sm text-gray-700 font-medium mb-2 flex items-center gap-2">
        <FileText size={16} /> File Change Log
      </div>
      
      {cachedImages.length === 0 ? (
        <div className="text-xs text-gray-500 italic">
          No images selected yet
        </div>
      ) : (
        <>
          {Object.entries(groupedImages).map(([key, images]) => (
            <div key={key} className="mb-2 pb-2 border-b border-gray-200 last:border-b-0">
              <div className="text-xs font-medium text-gray-700 mb-1">{key}</div>
              <div className="text-xs text-gray-600">
                {images.length} image{images.length > 1 ? 's' : ''}
              </div>
              <div className="max-h-20 overflow-y-auto">
                {images.map((img, index) => (
                  <div key={index} className="text-xs text-gray-500 whitespace-nowrap text-ellipsis overflow-hidden ml-1">
                    • {img.item_image_name}
                  </div>
                ))}
              </div>
            </div>
          ))}
          <div className="text-xs text-gray-600 pt-2 border-t border-gray-200">
            Total: {cachedImages.length} images
          </div>
        </>
      )}
    </div>
  );
}

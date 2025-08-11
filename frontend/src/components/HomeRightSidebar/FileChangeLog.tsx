import { FileText } from 'lucide-react';
import { FileChangeLogProps } from '@/types/HomeRightSidebar';

export default function FileChangeLog({ 
  groupedImages, cachedImages 
}: FileChangeLogProps) {

  return (
    <div className="change-log">
      <div className="title">
        <FileText size={16} /> File Change Log
      </div>
      
      {cachedImages.length === 0 ? (
        <div className="empty">
          No images selected yet
        </div>
      ) : (
        <>
          {Object.entries(groupedImages).map(([key, images]) => (
            <div key={key} className="group">
              <div className="group-header">{key}</div>
              <div className="group-count">
                {images.length} image{images.length > 1 ? 's' : ''}
              </div>
              <div className="file-list">
                {images.map((img, index) => (
                  <div key={index} className="file-item">
                    • {img.item_image_name}
                  </div>
                ))}
              </div>
            </div>
          ))}
          <div className="total">
            Total: {cachedImages.length} images
          </div>
        </>
      )}
    </div>
  );
}

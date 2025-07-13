package services

import (
	"backend/models"
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

func (dm *DataManager) GetAllPageDetail() map[int]string {
	log.Println("Fetching all page details from ImagesPerPageCache")
	if dm.ImagesPerPageCache.JobName == "" {
		log.Println("ImagesPerPageCache is not initialized")
		return nil
	}

	detail := make(map[int]string)
	for idx, page := range dm.ImagesPerPageCache.Pages {
		detail[idx] = page.DatasetName
	}
	return detail
}

func (dm *DataManager) InitialImagesCache(jobName string, pageNumber int) {
	dm.ImagesPerPageCache = models.NewImagesPerPageCache()
	log.Println("Cleared images cache")

	allDatasets := dm.GetParentDataAllDatasets(jobName)
	for _, datasetName := range allDatasets {
		images := dm.GetParentDataAllImages(jobName, datasetName)
		for i := 0; i < len(images); i += pageNumber {
			end := i + pageNumber
			if end > len(images) {
				end = len(images)
			}
			imageSet := models.ImageItems{
				DatasetName:    datasetName,
				ImageSet:       images[i:end],
				Base64ImageSet: []models.Image{},
			}
			dm.SaveImagesCache(jobName, imageSet)
		}
	}
	log.Printf("Initialized images cache for job %s with page size %d", jobName, pageNumber)
}

func (dm *DataManager) SaveImagesCache(jobName string, item models.ImageItems) {
	dm.ImagesPerPageCache.JobName = jobName
	dm.ImagesPerPageCache.Pages = append(dm.ImagesPerPageCache.Pages, item)
	dm.ImagesPerPageCache.MaxPage = len(dm.ImagesPerPageCache.Pages)
	log.Printf("Saved images cache for dataset: %s, total images: %d, total base64 image: %d", item.DatasetName, len(item.ImageSet), len(item.Base64ImageSet))
}

func (dm *DataManager) ClearImagesCache(jobName string) bool {
	if dm.ImagesPerPageCache.JobName == "" {
		log.Println("ImagesPerPageCache is not initialized")
		return true
	}

	log.Printf("compare 1:%s 2:%s", jobName, dm.ImagesPerPageCache.JobName)
	return jobName != dm.ImagesPerPageCache.JobName
}

func (dm *DataManager) GetImagesCache(jobName string, pageIndex int, pageNumber int) (models.ImageItems, bool) {
	exist := false
	if dm.ClearImagesCache(jobName) {
		dm.InitialImagesCache(jobName, pageNumber)
	}
	if len(dm.ImagesPerPageCache.Pages) <= pageIndex {
		log.Printf("Page index %d out of range for job %s", pageIndex, jobName)
		return models.ImageItems{}, exist
	}
	if len(dm.ImagesPerPageCache.Pages[pageIndex].Base64ImageSet) == 0 {
		log.Printf("No base64 images found for job %s at page index %d", jobName, pageIndex)
		return models.ImageItems{}, exist
	}

	cachedImages := dm.ImagesPerPageCache.Pages[pageIndex]
	if len(cachedImages.Base64ImageSet) > 0 {
		exist = true
	}

	return cachedImages, exist
}

func (dm *DataManager) GetBase64ImagesCache(jobName string, pageIndex int) {
	images := dm.ImagesPerPageCache.Pages[pageIndex].ImageSet
	for _, img := range images {
		base64Img := compressImages(img.Path)
		compressedImages := models.NewImage(img.Name, base64Img)
		dm.ImagesPerPageCache.Pages[pageIndex].Base64ImageSet = append(dm.ImagesPerPageCache.Pages[pageIndex].Base64ImageSet, compressedImages)
	}
	log.Printf("Compressed %d images for job %s at page index %d", len(dm.ImagesPerPageCache.Pages[pageIndex].Base64ImageSet), jobName, pageIndex)
}

func compressImages(imgPath string) string {
	log.Println("Processing image:", imgPath)
	file, err := os.Open(imgPath)
	if err != nil {
		log.Println("Failed to open image:", err)
		return ""
	}
	defer file.Close()

	decodedImg, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Failed to decode image (format: %s): %v\n", format, err)
		return ""
	}

	resizedImg := resize.Resize(150, 0, decodedImg, resize.Lanczos3)

	var buf bytes.Buffer
	opts := &webp.Options{Lossless: false, Quality: 75}
	if err := webp.Encode(&buf, resizedImg, opts); err != nil {
		log.Println("Failed to encode WebP:", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (dm *DataManager) ImageCacheJobExists(jobName string) bool {
	if dm.ImagesPerPageCache.JobName == "" {
		log.Println("ImagesPerPageCache is not initialized")
		return false
	}
	return dm.ImagesPerPageCache.JobName == jobName
}

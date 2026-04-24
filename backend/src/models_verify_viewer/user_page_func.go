package models_verify_viewer

import "log"

func (pages *Pages) JobName() string {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	return pages.jobName
}

func (pages *Pages) SetJobName(jobName string) {
	pages.mu.Lock()
	defer pages.mu.Unlock()

	pages.jobName = jobName
}

func (pages *Pages) GetDatasetNames() []string {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	datasets := make([]string, len(pages.datasets))
	copy(datasets, pages.datasets)
	return datasets
}

func (pages *Pages) PageItems() []PageItem {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	items := make([]PageItem, len(pages.pageItems))
	copy(items, pages.pageItems)
	return items
}

func (pages *Pages) PageItemsReadOnly() []PageItem {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	return pages.pageItems
}

func (pages *Pages) AddPage(datasetName string, imageSet []Image) {
	pages.mu.Lock()
	defer pages.mu.Unlock()

	pageItem := PageItem{
		DatasetName: datasetName,
		ImageSet:    imageSet,
	}
	pages.pageItems = append(pages.pageItems, pageItem)

	// Cache dataset name for this page (allows duplicates)
	pages.datasets = append(pages.datasets, datasetName)
}

func (pages *Pages) Clear() {
	pages.mu.Lock()
	defer pages.mu.Unlock()

	pages.jobName = ""
	pages.datasets = pages.datasets[:0]
	pages.pageItems = pages.pageItems[:0]
}

func (pages *Pages) PageAt(index int) (PageItem, bool) {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	if index < 0 || index >= len(pages.pageItems) {
		log.Printf("Index %d out of range in PageAt (len=%d)", index, len(pages.pageItems))
		return PageItem{}, false
	}
	return pages.pageItems[index], true
}

func (pages *Pages) ImageNamesAt(index int) []string {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	if index < 0 || index >= len(pages.pageItems) {
		return []string{}
	}

	pageItem := pages.pageItems[index]
	imageNames := make([]string, 0, len(pageItem.ImageSet))
	for _, image := range pageItem.ImageSet {
		imageNames = append(imageNames, image.Name)
	}
	return imageNames
}

func (pages *Pages) ImagePathsAt(index int) []string {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	if index < 0 || index >= len(pages.pageItems) {
		return []string{}
	}

	pageItem := pages.pageItems[index]
	imagePaths := make([]string, 0, len(pageItem.ImageSet))
	for _, image := range pageItem.ImageSet {
		imagePaths = append(imagePaths, image.Path)
	}
	return imagePaths
}

func (pages *Pages) Len() int {
	pages.mu.RLock()
	defer pages.mu.RUnlock()

	return len(pages.pageItems)
}

func (pages *Pages) RemoveImages(imagePaths []string) int {
	pages.mu.Lock()
	defer pages.mu.Unlock()

	if len(imagePaths) == 0 {
		return 0
	}

	pathsToRemove := make(map[string]bool, len(imagePaths))
	for _, path := range imagePaths {
		pathsToRemove[path] = true
	}

	totalRemoved := 0
	newPageItems := make([]PageItem, 0, len(pages.pageItems))

	for _, pageItem := range pages.pageItems {
		newImageSet := make([]Image, 0, len(pageItem.ImageSet))
		for _, image := range pageItem.ImageSet {
			if !pathsToRemove[image.Path] {
				newImageSet = append(newImageSet, image)
			} else {
				totalRemoved++
			}
		}

		if len(newImageSet) > 0 {
			pageItem.ImageSet = newImageSet
			newPageItems = append(newPageItems, pageItem)
		}
	}

	pages.pageItems = newPageItems

	if totalRemoved > 0 {
		log.Printf("Removed %d images from page data for job %s", totalRemoved, pages.jobName)
	}

	return totalRemoved
}

package models_verify_viewer

import "log"

// PageCache Function
func (pages *Pages) FillJobName(new_job_name string) {
	pages.JobName = new_job_name
}

func (pages *Pages) FillPageItemSet(new_page_items []PageItems) {
	pages.PageItems = new_page_items
}

func (pages *Pages) AppendPageItems(datasetName string, imageSet []Image) {
	new_page_item := NewPageItems()
	new_page_item.FillDatasetName(datasetName)
	new_page_item.FillImageSet(imageSet)
	pages.PageItems = append(pages.PageItems, new_page_item)
}

func (pages *Pages) ClearPages() {
	pages.JobName = ""
	pages.PageItems = NewPageItemsSet()
}

func (pages Pages) GetPageItemsByIndex(index int) PageItems {
	if index < 0 || index >= pages.GetPageItemsLength() {
		log.Println("Index out of range in GetPageItemByIndex, returning a new PageItems")
		return NewPageItems()
	}
	return pages.PageItems[index]
}

func (pages Pages) GetPageItemAllImageNameByIndex(index int) []string {
	pageItem := pages.GetPageItemsByIndex(index)
	return pageItem.GetAllImageName()
}

func (pages Pages) GetPageItemAllImagePathByIndex(index int) []string {
	pageItem := pages.GetPageItemsByIndex(index)
	return pageItem.GetAllImagePath()
}

func (pages Pages) GetPageItemsLength() int {
	return len(pages.PageItems)
}

// PageItems Function
func (page_item *PageItems) FillDatasetName(new_dataset_name string) {
	page_item.DatasetName = new_dataset_name
}

func (page_item *PageItems) FillImageSet(new_image_set []Image) {
	page_item.ImageSet = new_image_set
}

func (page_item PageItems) GetAllImageName() []string {
	var allImageNames []string
	for _, image := range page_item.ImageSet {
		allImageNames = append(allImageNames, image.Name)
	}
	return allImageNames
}

func (page_item PageItems) GetAllImagePath() []string {
	var allImagePaths []string
	for _, image := range page_item.ImageSet {
		allImagePaths = append(allImagePaths, image.Path)
	}
	return allImagePaths
}

func (page_item PageItems) GetImageSetLength() int {
	return len(page_item.ImageSet)
}

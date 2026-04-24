package models_verify_viewer

import "sync"

type Pages struct {
	jobName   string
	datasets  []string
	pageItems []PageItem
	mu        sync.RWMutex
}

type PageItem struct {
	DatasetName string  `json:"item_dataset_name"`
	ImageSet    []Image `json:"item_image_set"`
}

func NewPages() *Pages {
	return &Pages{
		jobName:   "",
		datasets:  make([]string, 0),
		pageItems: make([]PageItem, 0),
	}
}

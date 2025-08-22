package models_verify_viewer

type Pages struct {
	JobName   string      `json:"job_name"`
	PageItems []PageItems `json:"page_items"`
}

type PageItems struct {
	DatasetName string  `json:"item_dataset_name"`
	ImageSet    []Image `json:"item_image_set"`
}

func NewPages() Pages {
	return Pages{
		JobName:   "",
		PageItems: NewPageItemsSet(),
	}
}

func NewPageItems() PageItems {
	return PageItems{
		DatasetName: "",
		ImageSet:    NewImageSet(),
	}
}

func NewPageItemsSet() []PageItems {
	return []PageItems{}
}

package models_verify_viewer

type PendingReview struct {
	Items []PendingReviewItem `json:"items"`
}

type PendingReviewItem struct {
	JobName     string `json:"item_job_name"`
	DatasetName string `json:"item_dataset_name"`
	ImageName   string `json:"item_image_name"`
	ImagePath   string `json:"item_image_path"`
}

func NewPendingReview() PendingReview {
	return PendingReview{
		Items: NewPendingReviewItemSet(),
	}
}

func NewPendingReviewItem() PendingReviewItem {
	return PendingReviewItem{
		JobName:     "",
		DatasetName: "",
		ImageName:   "",
		ImagePath:   "",
	}
}

func NewPendingReviewItemSet() []PendingReviewItem {
	return []PendingReviewItem{}
}

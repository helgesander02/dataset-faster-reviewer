package models_verify_viewer

type PendingReview struct {
	Items []PendingReviewItem `json:"items"`
}

type PendingReviewItem struct {
	Job       string `json:"job"`
	Dataset   string `json:"dataset"`
	ImageName string `json:"imageName"`
	ImagePath string `json:"imagePath"`
}

func NewPendingReview() PendingReview {
	return PendingReview{
		Items: []PendingReviewItem{},
	}
}

func NewPendingReviewItem(job, dataset, image_name, image_path string) PendingReviewItem {
	return PendingReviewItem{
		Job:       job,
		Dataset:   dataset,
		ImageName: image_name,
		ImagePath: image_path,
	}
}

func NewPendingReviewItems() []PendingReviewItem {
	return []PendingReviewItem{}
}

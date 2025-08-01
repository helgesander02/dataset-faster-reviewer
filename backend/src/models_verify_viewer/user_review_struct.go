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
		Items: NewPendingReviewItems(),
	}
}

func NewPendingReviewItems() []PendingReviewItem {
	return []PendingReviewItem{}
}

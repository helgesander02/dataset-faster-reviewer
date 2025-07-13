package models

type PendingReview struct {
	Items []PendingReviewItem `json:"items"`
}

type PendingReviewItem struct {
	Job       string `json:"job"`
	Dataset   string `json:"dataset"`
	ImageName string `json:"imageName"`
	ImagePath string `json:"imagePath"`
}

type JobCache struct {
	Jobs map[string]Job `json:"jobs"`
}

type ImagesPerPageCache struct {
	JobName string       `json:"job_name"`
	MaxPage int          `json:"max_page"`
	Pages   []ImageItems `json:"pages"`
}

type ImageItems struct {
	DatasetName    string  `json:"dataset_name"`
	ImageSet       []Image `json:"image_set"`
	Base64ImageSet []Image `json:"base64_image_set"`
}

type Parent struct {
	Jobs []string `json:"jobs"`
}

type Job struct {
	Name     string    `json:"name"`
	Datasets []Dataset `json:"datasets"`
}

type Dataset struct {
	Name  string  `json:"name"`
	Image []Image `json:"images"`
	Label []Label `json:"labels"`
}

type Image struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Label struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

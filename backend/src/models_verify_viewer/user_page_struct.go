package models_verify_viewer

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

func NewImagesPerPageCache() ImagesPerPageCache {
	return ImagesPerPageCache{
		JobName: "",
		Pages:   []ImageItems{},
	}
}

func NewImageItems(datasetName string) ImageItems {
	return ImageItems{
		DatasetName: datasetName,
	}
}

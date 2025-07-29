package models

func NewJobCache() JobCache {
	return JobCache{
		Jobs: make(map[string]Job),
	}
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

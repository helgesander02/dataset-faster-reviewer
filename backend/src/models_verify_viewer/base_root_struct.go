package models_verify_viewer

type Job struct {
	Name     string    `json:"job_name"`
	Datasets []Dataset `json:"job_datasets"`
}

func NewJob() Job {
	return Job{
		Name:     "",
		Datasets: NewDatasetSet(),
	}
}

type Dataset struct {
	Name  string  `json:"dataset_name"`
	Image []Image `json:"dataset_images"`
	Label []Label `json:"dataset_labels"`
}

func NewDataset() Dataset {
	return Dataset{
		Name:  "",
		Image: NewImageSet(),
		Label: NewLabelSet(),
	}
}

func NewDatasetSet() []Dataset {
	return []Dataset{}
}

type Image struct {
	Name string `json:"image_name"`
	Path string `json:"image_path"`
}

func NewImage(imageName string, imagePath string) Image {
	return Image{
		Name: imageName,
		Path: imagePath,
	}
}

func NewImageSet() []Image {
	return []Image{}
}

type Label struct {
	Name string `json:"label_name"`
	Path string `json:"label_path"`
}

func NewLabel(labelName string, labelPath string) Label {
	return Label{
		Name: labelName,
		Path: labelPath,
	}
}

func NewLabelSet() []Label {
	return []Label{}
}

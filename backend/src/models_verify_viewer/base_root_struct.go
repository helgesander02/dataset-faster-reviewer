package models_verify_viewer

// Job, Dataset, Image, and Label are structures used to represent
//
// file structures on root level
// |- job-name
//   |- dataset-name
//	   |- image
//	   |- label

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

func NewImage(image_name string, image_path string) Image {
	return Image{
		Name: image_name,
		Path: image_path,
	}
}

func NewImageSet() []Image {
	return []Image{}
}

type Label struct {
	Name string `json:"label_name"`
	Path string `json:"label_path"`
}

func NewLabel(label_name string, label_path string) Label {
	return Label{
		Name: label_name,
		Path: label_path,
	}
}

func NewLabelSet() []Label {
	return []Label{}
}

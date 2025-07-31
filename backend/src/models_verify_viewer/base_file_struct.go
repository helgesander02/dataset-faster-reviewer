package models_verify_viewer

// Job, Dataset, Image, and Label are structures used to represent
//
// file structures on root level
// |- job-name
//   |- dataset-name
//	   |- image
//	   |- label

type Job struct {
	Name     string    `json:"name"`
	Datasets []Dataset `json:"datasets"`
}

func NewJob(job_name string) Job {
	return Job{
		Name:     job_name,
		Datasets: []Dataset{},
	}
}

type Dataset struct {
	Name  string  `json:"name"`
	Image []Image `json:"images"`
	Label []Label `json:"labels"`
}

func NewDataset(dataset_name string) Dataset {
	return Dataset{
		Name:  dataset_name,
		Image: []Image{},
		Label: []Label{},
	}
}

type Image struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func NewImage(image_name string, image_path string) Image {
	return Image{
		Name: image_name,
		Path: image_path,
	}
}

type Label struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func NewLabel(label_name string, label_path string) Label {
	return Label{
		Name: label_name,
		Path: label_path,
	}
}

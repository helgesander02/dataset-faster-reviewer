package models

func NewParentData() Parent {
	return Parent{
		Jobs: []string{},
	}
}

func NewJob(job_name string) Job {
	return Job{
		Name:     job_name,
		Datasets: []Dataset{},
	}
}

func NewDataset(dataset_name string) Dataset {
	return Dataset{
		Name:  dataset_name,
		Image: []Image{},
		Label: []Label{},
	}
}

func NewImage(image_name string, image_path string) Image {
	return Image{
		Name: image_name,
		Path: image_path,
	}
}

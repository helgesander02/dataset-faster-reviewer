package models_verify_viewer

func (job *Job) FillJobName(new_job_name string) {
	job.Name = new_job_name
}

func (job Job) GetDatasetLength() int {
	return len(job.Datasets)
}

func (dataset *Dataset) FillDatasetName(new_dataset_name string) {
	dataset.Name = new_dataset_name
}
func (dataset Dataset) GetImageLength() int {
	return len(dataset.Image)
}

package models_verify_viewer

func (job *Job) FillJobName(jobName string) {
	job.Name = jobName
}

func (job Job) GetDatasetLength() int {
	return len(job.Datasets)
}

func (dataset *Dataset) FillDatasetName(datasetName string) {
	dataset.Name = datasetName
}

func (dataset Dataset) GetImageLength() int {
	return len(dataset.Image)
}

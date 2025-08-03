package models_verify_viewer

func (job_list JobList) GetJobs() []string {
	return job_list.Jobs
}

func (job_list JobList) GetJobListLength() int {
	return len(job_list.Jobs)
}

func (job_list JobList) JobExists(jobName string) bool {
	for _, j := range job_list.Jobs {
		if j == jobName {
			return true
		}
	}
	return false
}

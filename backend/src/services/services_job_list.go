package services

// job function
func (js *JointServices) JobExists(jobName string) bool {
	return js.JobList.JobExists(jobName)
}

func (js *JointServices) GetJobList() []string {
	return js.JobList.GetJobs()
}

package services

func (js *JointServices) JobExists(jobName string) bool {
	return js.JobList.Contains(jobName)
}

func (js *JointServices) GetJobList() []string {
	return js.JobList.Jobs()
}

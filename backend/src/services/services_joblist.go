package services

func (js *JointServices) GetAllJobs() []string {
	return js.JobList.Jobs
}

func (js *JointServices) DataJobExists(jobName string) bool {
	for _, j := range js.JobList.Jobs {
		if j == jobName {
			return true
		}
	}
	return false
}

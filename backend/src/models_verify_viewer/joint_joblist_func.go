package models_verify_viewer

func (jl *JobList) Jobs() []string {
	jl.mu.RLock()
	defer jl.mu.RUnlock()

	result := make([]string, len(jl.jobs))
	copy(result, jl.jobs)
	return result
}

func (jl *JobList) Len() int {
	jl.mu.RLock()
	defer jl.mu.RUnlock()

	return len(jl.jobs)
}

func (jl *JobList) Contains(jobName string) bool {
	jl.mu.RLock()
	defer jl.mu.RUnlock()

	for _, j := range jl.jobs {
		if j == jobName {
			return true
		}
	}
	return false
}

func (jl *JobList) Add(jobName string) {
	jl.mu.Lock()
	defer jl.mu.Unlock()

	for _, j := range jl.jobs {
		if j == jobName {
			return
		}
	}
	jl.jobs = append(jl.jobs, jobName)
}

func (jl *JobList) Replace(jobs []string) {
	jl.mu.Lock()
	defer jl.mu.Unlock()

	jl.jobs = jobs
}

func (jl *JobList) Remove(jobName string) {
	jl.mu.Lock()
	defer jl.mu.Unlock()

	newJobs := make([]string, 0, len(jl.jobs))
	for _, j := range jl.jobs {
		if j != jobName {
			newJobs = append(newJobs, j)
		}
	}
	jl.jobs = newJobs
}

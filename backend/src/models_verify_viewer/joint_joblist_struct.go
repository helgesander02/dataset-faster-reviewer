package models_verify_viewer

// JobList is a global structure used to store the names of all jobs

type JobList struct {
	Jobs []string `json:"jobs"`
}

func NewJobList() JobList {
	return JobList{
		Jobs: []string{},
	}
}

package cache

import "fmt"

type ErrCacheNotFound struct {
	JobName string
}

func (e ErrCacheNotFound) Error() string {
	return fmt.Sprintf("cache not found for job: %s", e.JobName)
}

type ErrInvalidCacheData struct {
	JobName string
}

func (e ErrInvalidCacheData) Error() string {
	return fmt.Sprintf("invalid cache data type for job: %s", e.JobName)
}

type ErrPageIndexOutOfRange struct {
	JobName   string
	PageIndex int
	MaxPages  int
}

func (e ErrPageIndexOutOfRange) Error() string {
	return fmt.Sprintf("page index %d out of range for job %s (max pages: %d)", 
		e.PageIndex, e.JobName, e.MaxPages)
}

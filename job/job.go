package job

import "strconv"

// Job ...
type Job struct {
	ID        string
	Operation func()
}

// Status ...
type Status struct {
	ID     string
	Result string
}

// JobExecutionDetails ...
type JobExecutionDetails struct {
	ID          string
	Status      string
	Detailsteps interface{}
}

// AssignJobs write jobs to the channel
func AssignJobs(jobch chan<- Job, nojobs int) {
	for _, job := range getjobs(nojobs) {
		jobch <- job
	}
	close(jobch)
}

// GetJobs stub, it can be received from any data store
func getjobs(jobcnt int) []Job {
	var jobs []Job

	for i := 0; i < jobcnt; i++ {
		//get the uid and assign it as JobID
		id := strconv.Itoa(i)

		jobs = append(jobs, Job{ID: id})
	}
	return jobs

}

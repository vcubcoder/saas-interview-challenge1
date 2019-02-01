package main

import (
	"strconv"
	"sync"
)

// Job ...
type Job struct {
	ID        string
	Operation string
}

// JobStatus ...
type JobStatus struct {
	ID   string
	Pass bool
	//ErrorMsg string
}

func main() {

	noworkers := 5
	nojobs := 20

	jobch := make(chan Job, 10)
	workstatusch := make(chan JobStatus, 10)

	var wg sync.WaitGroup
	wg.Add(noworkers)

	// create workers
	for i := 0; i < noworkers; i++ {
		go func() {
			worker(jobch, workstatusch)
			wg.Done()
		}()
	}
	wg.Wait()

	//Assign jobs
	// joblist := getJobs()
	for _, job := range getJobs(nojobs) {
		jobch <- job
	}
	close(jobch)

	for i := 0; i < nojobs; i++ {
		res := <-workstatusch
		if !res.Pass {
			//get the data from redis and print the error
		}
	}
}

func getJobs(jobcnt int) []Job {
	var jobs []Job

	for i := 0; i < jobcnt; i++ {
		//get the uid and assign it as JobID
		id := strconv.Itoa(i)

		jobs = append(jobs, Job{ID: id, Operation: "operation" + id})
	}
	return jobs

}
func worker(jobs <-chan Job, status chan<- JobStatus) {
	for job := range jobs {
		//Perform the operation
		// Write the info on to redis

		//write the status
		status <- JobStatus{ID: job.ID, Pass: true}
	}
}

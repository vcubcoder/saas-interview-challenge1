package main

import (
	"fmt"
	"math/rand"
	redis "saas-interview-challenge1/redisprovider"
	"strconv"
	"sync"
)

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

func main() {

	// This program assumes fixed number of workers and fixed number of jobs, thus job jobs channels are closed
	// and workers are not open ended/unbound

	noworkers := 5
	nojobs := 10

	jobch := make(chan Job, 10)
	workstatusch := make(chan Status, 10)
	rclient := redis.NewRedisClient()

	//Assign jobs
	go assignjobs(jobch, nojobs)

	// create worker pool
	var wg sync.WaitGroup
	wg.Add(noworkers)

	for i := 0; i < noworkers; i++ {
		go func() {
			worker(jobch, workstatusch, rclient)
			wg.Done()
		}()
	}
	wg.Wait() // wait for all the workers to finish

	//After all the workers are done get the job results
	for i := 0; i < nojobs; i++ {
		res := <-workstatusch
		jd, _ := rclient.GetJobExecutionDetails(res.ID)
		if jd.Status != "Pass" {
			//get the data from redis and print the error
			fmt.Printf("JobID: %s, Statu: %s, Job steps staus : %v\n", jd.ID, jd.Status, jd.Detailsteps)
		} else {
			fmt.Printf("JobID : %s, Status : %s\n", jd.ID, jd.Status)
		}
	}
}

// Get Jobs stub, it can be received from any data store
func getJobs(jobcnt int) []Job {
	var jobs []Job

	for i := 0; i < jobcnt; i++ {
		//get the uid and assign it as JobID
		id := strconv.Itoa(i)

		jobs = append(jobs, Job{ID: id})
	}
	return jobs

}

// worker receives the job in the channel
// 1. perform the required operations
// 2. write the execution details on to redis
// 3  send the job status on to status channel
func worker(jobs <-chan Job, status chan<- Status, rclient redis.Client) {
	for job := range jobs {
		jobstatus := Status{ID: job.ID}
		var err error
		//Perform the job.operation()
		if (rand.Intn(10) % 2) == 0 {
			jobstatus.Result = "Pass"
			// Write the info on to redis
			err = rclient.PutJobExecutionDetails(job.ID, &redis.JobExecutionDetails{ID: job.ID, Status: "Pass"})
		} else {
			jobstatus.Result = "Fail"

			//Write Job execution detail on to redis
			jd := &redis.JobExecutionDetails{
				ID:     job.ID,
				Status: "Fail",
				Detailsteps: map[string]string{
					"ec2-create": "Passed",
					"sg-create":  "Passed",
					"asg-create": "Failed",
				},
			}
			err = rclient.PutJobExecutionDetails(job.ID, jd)
		}
		if err != nil {
			fmt.Printf("Writing to Redis failed for the job: %s with error message: %s\n", job.ID, err)
		}
		//write to status channel
		status <- jobstatus
	}
}

func assignjobs(jobch chan<- Job, nojobs int) {
	for _, job := range getJobs(nojobs) {
		jobch <- job
	}
	close(jobch)
}

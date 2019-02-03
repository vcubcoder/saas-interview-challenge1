package worker

import (
	"fmt"
	"math/rand"
	"saas-interview-challenge1/job"
	redis "saas-interview-challenge1/redisprovider"
)

// Worker ...
type Worker struct {
	jobch    <-chan job.Job
	resultch chan<- job.Status
	client   redis.Client
}

// Run receives the job in the channel
// 1. perform the required operations
// 2. write the execution details on to redis
// 3  send the job status on to status channel
func (w Worker) Run() {
	for jb := range w.jobch {
		jobstatus := job.Status{ID: jb.ID}
		var err error
		//Perform the job.operation()
		if (rand.Intn(10) % 2) == 0 {
			jobstatus.Result = "Pass"

			//Write Job execution detail on to redis
			err = w.client.PutJobExecutionDetails(jb.ID, &job.JobExecutionDetails{ID: jb.ID, Status: "Pass"})
		} else {
			jobstatus.Result = "Fail"

			//Write Job execution detail on to redis
			jd := &job.JobExecutionDetails{
				ID:     jb.ID,
				Status: "Fail",
				Detailsteps: map[string]string{
					"ec2-create": "Passed",
					"sg-create":  "Passed",
					"asg-create": "Failed",
				},
			}
			err = w.client.PutJobExecutionDetails(jb.ID, jd)
		}
		if err != nil {
			fmt.Printf("Writing to Redis failed for the job: %s with error message: %s\n", jb.ID, err)
		}
		//write to status channel
		w.resultch <- jobstatus
	}
}

// NewWorker ...
func NewWorker(jobch <-chan job.Job, statusch chan<- job.Status, rclient redis.Client) *Worker {
	return &Worker{
		jobch:    jobch,
		resultch: statusch,
		client:   rclient,
	}
}

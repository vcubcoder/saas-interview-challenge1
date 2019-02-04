package worker

import (
	"saas-interview-challenge1/job"
	redis "saas-interview-challenge1/redisprovider"
	"testing"
)

func TestGet(t *testing.T) {

	jobch := make(chan job.Job, 10)
	resultch := make(chan job.Status, 10)
	rclient := redis.NewRedismockClient()

	//Assign jobs
	job.AssignJobs(jobch, 1)
	worker := NewWorker(jobch, resultch, rclient)
	worker.Run()
	jb := <-resultch
	jd, _ := worker.client.GetJobExecutionDetails(jb.ID)

	if jd.ID != "ID1" {
		t.Errorf("Job ID is incorrect, got: %s, want: ID1.", jd.ID)
		t.Fail()
	}

}

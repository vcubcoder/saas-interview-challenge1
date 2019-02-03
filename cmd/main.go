package main

import (
	"fmt"
	"sync"

	"saas-interview-challenge1/job"
	redis "saas-interview-challenge1/redisprovider"
	"saas-interview-challenge1/worker"
)

func main() {

	// This program assumes fixed number of workers and fixed number of jobs, thus job jobs channels are closed

	noworkers := 5
	nojobs := 10

	jobch := make(chan job.Job, 10)
	resultch := make(chan job.Status, 10)
	rclient := redis.NewRedisClient("localhost")

	worker := worker.NewWorker(jobch, resultch, rclient)

	//Assign jobs
	go job.AssignJobs(jobch, nojobs)

	// create worker pool(workers are not open ended) and run them
	var wg sync.WaitGroup
	wg.Add(noworkers)
	for i := 0; i < noworkers; i++ {
		go func() {
			worker.Run()
			wg.Done()
		}()
	}
	wg.Wait() // wait for all the workers to finish

	//After all the workers are done collect the job results
	for i := 0; i < nojobs; i++ {
		res := <-resultch
		jd, _ := rclient.GetJobExecutionDetails(res.ID)
		if jd.Status != "Pass" {
			//get the data from redis and print the error
			fmt.Printf("JobID: %s, Statu: %s, Job steps staus : %v\n", jd.ID, jd.Status, jd.Detailsteps)
		} else {
			fmt.Printf("JobID : %s, Status : %s\n", jd.ID, jd.Status)
		}
	}
}

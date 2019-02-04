package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"saas-interview-challenge1/job"
	redis "saas-interview-challenge1/redisprovider"
	"saas-interview-challenge1/worker"
)

func main() {

	// This program assumes fixed number of workers and fixed number of jobs,
	// jobs channel is closed after assigning the jobs

	redisip := flag.String("redis-host", "", "10.11.192.1 (mandatory param)")
	workers := flag.Int("workers", 4, "number of workers, defaulted to 4")

	flag.Parse()

	if *redisip == "" {
		fmt.Println("Redis Host is not provided")
		flag.PrintDefaults()
		os.Exit(1)
	}

	noworkers := *workers
	nojobs := 10

	fmt.Println("No of Workers :", noworkers)
	fmt.Println("Redis Host", *redisip)

	jobch := make(chan job.Job, 10)
	resultch := make(chan job.Status, 10)
	rclient := redis.NewRedisClient(*redisip)
	worker := worker.NewWorker(jobch, resultch, rclient)

	//Assign jobs
	go job.AssignJobs(jobch, nojobs)

	// create worker pool and run them
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
		jd, err := rclient.GetJobExecutionDetails(res.ID)
		if err != nil {
			fmt.Printf("Pull job execution details failed. job id: %s, error : %s", res.ID, err)
		} else {
			if jd.Status != "Pass" {
				//get the data from redis and print the error
				fmt.Printf("JobID: %s, Statu: %s, Job steps staus : %v\n", jd.ID, jd.Status, jd.Detailsteps)
			} else {
				fmt.Printf("JobID : %s, Status : %s\n", jd.ID, jd.Status)
			}
		}
	}
}

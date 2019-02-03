package redisprovider

import "saas-interview-challenge1/job"

type mockredisclient struct{}

func (mrc *mockredisclient) GetJobExecutionDetails(jobid string) (*job.JobExecutionDetails, error) {
	return &job.JobExecutionDetails{
		ID:     "ID1",
		Status: "Pass",
	}, nil
}

func (mrc *mockredisclient) PutJobExecutionDetails(key string, value interface{}) error {
	return nil
}

// NewRedismockClient ...
func NewRedismockClient() Client {
	return &mockredisclient{}
}

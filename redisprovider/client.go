package redisprovider

import (
	"encoding/json"
	"errors"
	"fmt"
	"saas-interview-challenge1/job"

	"github.com/go-redis/redis"
)

// Client ...
type Client interface {
	PutJobExecutionDetails(string, interface{}) error
	GetJobExecutionDetails(string) (*job.JobExecutionDetails, error)
}

// RedisClient - warpper for redis client interface
type redisClient struct {
	client *redis.Client
}

// GetJobExecutionDetails gets the detailed job status from data store
func (rc *redisClient) GetJobExecutionDetails(jobid string) (*job.JobExecutionDetails, error) {
	var jd job.JobExecutionDetails
	val, err := rc.client.Get(jobid).Bytes()

	if err == redis.Nil {
		return nil, errors.New("jodid: " + jobid + " not exist")
	} else if err != nil {
		return nil, err
	}
	err = json.Unmarshal(val, &jd)
	if err != nil {
		fmt.Println("Redis data Unmarshalling failed")
		return nil, err
	}
	return &jd, nil
}

// PutJobExecutionDetails write the job execution details on to data store
func (rc *redisClient) PutJobExecutionDetails(key string, value interface{}) error {
	val, _ := json.Marshal(value)
	err := rc.client.Set(key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// NewRedisClient ...
func NewRedisClient(addr string) Client {
	internalredis := redis.NewClient(&redis.Options{
		Addr:     addr + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &redisClient{
		client: internalredis,
	}
}

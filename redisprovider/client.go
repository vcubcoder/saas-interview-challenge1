package redisprovider

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// Client ...
type Client interface {
	// Get(string) (error, interface{})
	PutJobExecutionDetails(string, interface{}) error
	GetJobExecutionDetails(string) (*JobExecutionDetails, error)
}

// JobExecutionDetails ...
type JobExecutionDetails struct {
	ID          string
	Status      string
	Detailsteps interface{}
}

// RedisClient - warpper for redis client interface
type redisClient struct {
	client *redis.Client
}

// GetJobExecutionDetails gets the detailed job status from data store
func (rc *redisClient) GetJobExecutionDetails(jobid string) (*JobExecutionDetails, error) {
	var jd JobExecutionDetails
	val, err := rc.client.Get(jobid).Bytes()

	if err == redis.Nil {
		return nil, errors.New("jodid: " + jobid + " not exist")
	} else if err != nil {
		return nil, err
	}
	err = json.Unmarshal(val, &jd)
	return &jd, nil
}

// func (rc *redisClient) Get() (interface{}, error) {
// 	return nil, nil
// }

func (rc *redisClient) PutJobExecutionDetails(key string, value interface{}) error {
	val, _ := json.Marshal(value)
	err := rc.client.Set(key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// NewRedisClient ...
func NewRedisClient() Client {
	internalredis := redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	return &redisClient{
		client: internalredis,
	}
}

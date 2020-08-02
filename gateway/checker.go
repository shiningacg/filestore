package gateway

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type Checker interface {
	Get(token string) (*CheckResult, error)
	Set(result *CheckResult) error
}

type CheckResult struct {
	Token  string
	Name   string
	Size   uint64
	Status bool
}

func (c *CheckResult) Checked() *CheckResult {
	c.Status = true
	return c
}

func (c *CheckResult) Decode(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *CheckResult) Encode() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewRedisChecker(redisAddr string, secret string) (*RedisChecker, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr, Password: secret})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &RedisChecker{
		Client: rdb,
	}, nil
}

type RedisChecker struct {
	*redis.Client
}

func (c *RedisChecker) Get(token string) (*CheckResult, error) {
	result := c.Client.Get(context.TODO(), token)
	b, err := result.Bytes()
	if err != nil {
		return nil, err
	}
	var checkResult = &CheckResult{}
	return checkResult, checkResult.Decode(b)
}

func (c *RedisChecker) Set(result *CheckResult) error {
	b, err := result.Encode()
	if err != nil {
		return err
	}
	res := c.Client.Set(context.TODO(), result.Token, b, time.Hour*24)
	return res.Err()
}

type GrpcChecker struct{}

func (g GrpcChecker) Get(token string) (*CheckResult, error) {
	panic("implement me")
}

type HttpChecker struct{}

func (c *HttpChecker) Get(token string) (*CheckResult, error) {
	panic("implement me")
}

type MockChecker struct{}

func (m MockChecker) Get(token string) (*CheckResult, error) {
	return &CheckResult{
		Token:  token,
		Name:   "mock",
		Size:   1024,
		Status: false,
	}, nil
}

func (m MockChecker) Set(result *CheckResult) error {
	return nil
}

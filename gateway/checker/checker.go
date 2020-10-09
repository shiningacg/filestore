package checker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	svc "github.com/shiningacg/ServiceFile"
	"google.golang.org/grpc"
	"time"
)

// check 用户上传时的权限检测
type Checker interface {
	Get(token string) (*CheckResult, error)
	Set(result *CheckResult) error
}

type CheckResult struct {
	PostToken string
	UUID      string
	Name      string
	Size      uint64
}

func (c *CheckResult) Checked(uuid string) *CheckResult {
	c.UUID = uuid
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
	res := c.Client.Set(context.TODO(), result.PostToken, b, time.Hour*24)
	return res.Err()
}

func NewGrpcChecker(addr string, secret string) (*GrpcChecker, error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.New("无法建立与checker的连接")
	}
	client := svc.NewFileClient(conn)
	return &GrpcChecker{FileClient: client}, nil
}

type GrpcChecker struct {
	svc.FileClient
}

func (g *GrpcChecker) Set(result *CheckResult) error {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*5)
	defer cf()
	_, err := g.FileClient.UploadCallback(ctx, &svc.UploadCallbackRequest{
		Uuid:  result.UUID,
		Node:  result.Name,
		Token: result.PostToken,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *GrpcChecker) Get(token string) (*CheckResult, error) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*5)
	defer cf()
	info, err := g.FileClient.File(ctx, &svc.FileRequest{Fid: token})
	if err != nil {
		return nil, err
	}
	return &CheckResult{
		PostToken: token,
		UUID:      "",
		Name:      info.Name,
		Size:      info.Size,
	}, nil
}

type HttpChecker struct{}

func (c *HttpChecker) Get(token string) (*CheckResult, error) {
	panic("implement me")
}

type MockChecker struct{}

func (m MockChecker) Get(token string) (*CheckResult, error) {
	return &CheckResult{
		PostToken: token,
		UUID:      "test",
		Name:      "mock",
		Size:      12,
	}, nil
}

func (m MockChecker) Set(result *CheckResult) error {
	return nil
}

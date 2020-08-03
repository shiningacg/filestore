package gateway

import (
	"fmt"
	"testing"
)

func testNewRedisChecker() Checker {
	checker, err := NewRedisChecker("127.0.0.1:6379", "")
	if err != nil {
		panic(err)
	}
	return checker
}

func TestSetChecker(t *testing.T) {
	checker := testNewRedisChecker()
	var checkResult = &CheckResult{
		PostToken: "aaa",
		Name:      "a.test",
		Size:      1024,
		Status:    false,
	}
	err := checker.Set(checkResult)
	if err != nil {
		panic(err)
	}
	res, err := checker.Get(checkResult.PostToken)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res)
}

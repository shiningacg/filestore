package gateway

import (
	"strconv"
	"strings"
)

func ParseRange(rg string) []int {
	var rang = []int{0, 0}
	// 如果不需要控制长度
	if rg == "" {
		return rang
	}
	// 如果需要控制长度，那么就必须正确
	temp := strings.Split(rg[6:], "-")
	if len(temp) != 2 {
		return nil
	}
	head, err := strconv.ParseInt(temp[0], 10, 64)
	if err == nil {
		rang[0] = int(head)
	}
	foot, err := strconv.ParseInt(temp[1], 10, 64)
	if err == nil {
		rang[1] = int(foot)
	}
	return rang
}

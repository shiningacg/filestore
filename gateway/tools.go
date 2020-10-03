package gateway

import (
	"strconv"
	"strings"
)

func ParseRange(rg string) []int {
	var rang = []int{0, 0}
	if rg == "" {
		return rang
	}
	temp := strings.Split(rg, "-")
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

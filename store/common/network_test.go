package common

import (
	"container/ring"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDefaultNetwork_Stats(t *testing.T) {
	n := NewDefaultNetwork(context.Background())
	for {
		time.Sleep(time.Second * 5)
		fmt.Println(n.Stats())
	}
}
func TestNewDefaultNetwork(t *testing.T) {
	n := &DefaultNetwork{
		Ring: ring.New(61),
		len:  60,
	}
	n.next().Download = 111
	n.next().Download = 222
	fmt.Println(n.get(-1))
	fmt.Println(n.Value)
}

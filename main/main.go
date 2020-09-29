package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	n := NewDefaultNetwork(context.Background())
	for {
		time.Sleep(time.Second * 5)
		fmt.Println(n.Stats())
	}
}

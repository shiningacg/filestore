package gateway

import (
	"fmt"
	"testing"
)

func TestParseRange(t *testing.T) {
	range1 := "3-12"
	fmt.Println(ParseRange(range1))
}

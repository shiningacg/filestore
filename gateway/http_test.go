package gateway

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestUUID(t *testing.T) {
	fmt.Println(uuid.New().String())
}

func TestGetAction(t *testing.T) {
	fmt.Println(getAction("/get/aaa"))
}

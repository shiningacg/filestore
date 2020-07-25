package remote

import (
	"fmt"
	"github.com/shiningacg/filestore"
	"io/ioutil"
	"testing"
)

// 通过
func TestHttpAdder_Find(t *testing.T) {
	adder := NewHttpAdder("127.0.0.1:8888")
	var bf = &filestore.BaseFileStruct{}
	bf.SetUUID("aaaa")
	rf := adder.Find(bf)
	bt, _ := ioutil.ReadAll(rf)
	fmt.Println(string(bt))
}

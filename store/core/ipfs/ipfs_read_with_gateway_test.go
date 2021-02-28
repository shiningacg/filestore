package ipfs

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestHttpBodyReader_Read(t *testing.T) {
	st, err := NewCore(&log.Logger{})
	if err != nil {
		panic(err)
	}
	reader := &HttpBodyReader{
		Node:     st.ipfs.Get("QmWCXym1Y3mSGxzD7zeMJSAt8DJrYyQwm5KbwyfVzWEkQS"),
		index:    0,
		Response: nil,
	}
	f, err := os.Create("xx.mp4")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(f, reader)
	if err != nil {
		panic(err)
	}
}

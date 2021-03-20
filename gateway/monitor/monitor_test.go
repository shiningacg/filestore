package monitor

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"
)

//func TestCrontab(t *testing.T)  {
//	// 每天的17：27分执行一次
//	spec := "0 27 17 * * *"
//	c := cron.New()
//	c.AddFunc(spec, func() {
//		fmt.Println("hi")
//	})
//	c.Start()
//
//	select {}
//}

var _monitor = NewMonitor()

func TestCopy(t *testing.T) {
	f, _ := os.Create("test")
	data := []byte("aaa")
	buffer := bufio.NewReader(bytes.NewBuffer(data))
	reader := _monitor.In(buffer, "aaa", 0)
	n, err := io.Copy(f, reader)
	if err != nil {
		panic(err)
	}
	if int(n) != len(data) {
		panic("not equal")
	}
	if _monitor.Stats().In.Total != uint64(n) {
		panic("not equal")
	}
}

func TestCopyWithLimit(t *testing.T) {
	f, err := os.Open("/Users/shlande/Pictures/06.webp")
	if err != nil {
		panic(err)
	}
	stats, _ := f.Stat()
	w, err := os.Create("./06.webp")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	reader := _monitor.In(f, "test", stats.Size()-1)
	defer reader.Close()
	_, err = io.Copy(w, reader)
	if err != ErrReachMaxSize {
		panic(err)
	}

}

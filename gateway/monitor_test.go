package gateway

import (
	"context"
	store "github.com/shiningacg/filestore"
	"reflect"
	"testing"
	"time"
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

func TestAddRecord(t *testing.T) {
	record1 := &Record{
		RequestID: "aaa",
		Ip:        "11",
		FileID:    "dsf",
		Bandwidth: 1231,
		StartTime: uint64(time.Now().Unix()),
		EndTime:   0,
	}
	record2 := &Record{
		RequestID: "aaa",
		Ip:        "",
		FileID:    "",
		Bandwidth: 251351,
		StartTime: 0,
		EndTime:   0,
	}
	record3 := &Record{
		RequestID: "aaa",
		Ip:        "",
		FileID:    "",
		Bandwidth: 322,
		StartTime: 0,
		EndTime:   uint64(time.Now().Unix()) + 1,
	}
	monitor := NewMonitor(context.Background())
	go monitor.Run()
	monitor.AddRecord(record1)
	// 第一次添加后，测试visit和bandwidth是否符合
	time.Sleep(time.Millisecond * 10)
	bd := record1.Bandwidth
	if monitor.bandwidth != bd {
		t.Fatal("bandwidth err :", monitor.bandwidth, bd)
	}

	// 第二次添加后测试record是否被更新，总流量是否被更新
	monitor.AddRecord(record2)
	time.Sleep(time.Millisecond * 10)
	bd = record1.Bandwidth + record2.Bandwidth
	if monitor.bandwidth != bd {
		t.Fatal("bandwidth err :", monitor.bandwidth, bd)
	}
	wantRecord := &Record{
		RequestID: record1.RequestID,
		Ip:        record1.Ip,
		FileID:    record1.FileID,
		Bandwidth: bd,
		StartTime: record1.StartTime,
		EndTime:   0,
	}
	gotRecord := monitor.records[0]
	if !reflect.DeepEqual(wantRecord, gotRecord) {
		t.Fatalf("record err : %#v %#v", wantRecord, gotRecord)
	}

	// 测试后添加以后是否正确获得时间
	monitor.AddRecord(record3)
	time.Sleep(time.Millisecond * 10)
	wantRecord.EndTime = record3.EndTime
	wantRecord.Bandwidth += record3.Bandwidth
	if !reflect.DeepEqual(wantRecord, gotRecord) {
		t.Fatalf("record err : %#v %#v", wantRecord, gotRecord)
	}
	// 测试最后测试一下输出是否正确
	outputBandwidth := record1.Bandwidth + record2.Bandwidth + record3.Bandwidth
	wantGatway := &store.Gateway{
		Visit:         1,
		DayVisit:      1,
		HourVisit:     1,
		Bandwidth:     outputBandwidth,
		DayBandwidth:  outputBandwidth,
		HourBandwidth: outputBandwidth,
	}
	if gotGateway := monitor.Bandwidth(); !reflect.DeepEqual(*wantGatway, *gotGateway) {
		t.Fatal("gateway err ：", wantGatway, gotGateway)
	}
}

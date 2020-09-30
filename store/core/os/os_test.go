package os

// TODO：重写测试代码
//func testOpenStore() *Store {
//	log.OpenLog(&log.Config{})
//	logger := log.Default()
//	return NewCore(&StoreConfig{
//		GatewayAddr: ":8887",
//		StorePath:   ".",
//	}, checker.MockChecker{}, logger)
//}
//
//func TestNewOStore(t *testing.T) {
//	store := testOpenStore()
//	go func() {
//		for {
//			time.Sleep(time.Second * 2)
//			fmt.Println(store.Gateway())
//		}
//	}()
//	err := store.gateway.Run(context.TODO())
//	if err != nil {
//		panic(err)
//	}
//}
//
//func TestAPI_Add(t *testing.T) {
//	var bs = &fs.BaseFileStruct{}
//	f, _ := os.Open("./aaa")
//	bs.SetName("aaa")
//	bs.SetSize(100)
//	store := testOpenStore()
//	err := store.Add(fs.NewReadableFile(bs, f))
//	if err != nil {
//		panic(err)
//	}
//}

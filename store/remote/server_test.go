package remote

//func TestNewStoreServer(t *testing.T) {
//	log.OpenLog(&log.Config{})
//	g := gateway.DesignMyginGateway(":8888", checker.MockChecker{})
//	s := mock.NewCore()
//
//	NewStoreGRPCServer("127.0.0.1:5060", MockAdder{}, store)
//	for {
//		fmt.Println(g.BandWidth())
//		time.Sleep(time.Second * 10)
//	}
//}
//
//// TODO： 重写！！！
//func TestNewStoreServerWithRedisChecker(t *testing.T) {
//	log.OpenLog(&log.Config{})
//	checker, err := checker.NewRedisChecker("127.0.0.1:6379", "")
//	if err != nil {
//		panic(err)
//	}
//	g := gateway.NewMyginGateway(":8888", checker)
//	store := mock2.NewFileStore(g)
//	g.SetStore(store)
//	NewStoreGRPCServer("127.0.0.1:6666", MockAdder{}, store)
//	for {
//		fmt.Println(g.BandWidth())
//		time.Sleep(time.Second * 10)
//	}
//}

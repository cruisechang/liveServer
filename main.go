package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cruisechang/liveServer/builtinevent"
	"github.com/cruisechang/liveServer/command"
	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/dbConf"
	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/liveServer/method"
	pb "github.com/cruisechang/liveServer/protobuf"
	"github.com/cruisechang/liveServer/rpc"

	nx "github.com/cruisechang/nex"
	nxBuiltinEvent "github.com/cruisechang/nex/builtinEvent"
	nxLog "github.com/cruisechang/nex/log"
)

//加新遊戲
//fake create room
func main() {


	nex, err := nx.NewNex(getConfigFilePosition("nexConfig.json"))
	if err != nil {
		exit(1, fmt.Errorf("NewNex error:%s", err.Error()))
	}
	logger := nex.GetLogger()


	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic,fmt.Sprintf("main panic=%v",r))
		}
	}()

	//configurer，所有程式共用同一個
	conf, err := config.NewConfigurer("config.json")
	if err != nil {
		logger.LogFile(nxLog.LevelError, err.Error())
		exit(2, fmt.Errorf("loadConfig error:%s", err.Error()))
	}

	rCtrl := roomCtrl.NewController(conf)

	dbCtrl := control.NewDBController(conf.DBAPIServer())

	//query limitations
	limitations, err := method.QueryLimitations(dbCtrl, "GET", dbConf.DBPathLimitations)
	if err != nil {
		exit(100, err)
	}
	for _, v := range limitations {
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("limitation from db=%+v", v))
	}

	transferredLimit, err := method.TransferLimitations(limitations)
	if err != nil {
		logger.LogFile(nxLog.LevelError, err.Error())
		exit(101, err)
	}

	//query dealers
	dealers, err := method.QueryDealers(dbCtrl, "GET", dbConf.DBPathDealers)
	if err != nil {
		exit(102, err)
	}
	for _, v := range dealers {
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("dealer from db=%+v", v))
	}

	//query room
	roomData, err := method.QueryRooms(dbCtrl, "GET", dbConf.DBPathRooms)
	if err != nil {
		exit(200, err)
	}

	//query halls
	hallData, err := method.QueryHalls(dbCtrl, "GET", dbConf.DBPathHalls)
	if err != nil {
		exit(250, err)
	}

	//room manager
	rm := nex.GetRoomManager()
	rooms, err := method.CreateRooms(rm, rCtrl, roomData, transferredLimit, dealers)
	if err != nil {
		logger.LogFile(nxLog.LevelError, err.Error())
		exit(251, err)
	}

	hm := nex.GetHallManager()

	method.CreateHalls(hm, hallData)
	halls := hm.GetHalls()
	for _, h := range halls {
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("hall complete %+v", h))
	}

	method.SetHallRoom(halls, rooms)

	//init road map data
	rmc:=control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl,nex.GetLogger())
	err=rmc.InitRoadMapData(rooms)
	if err!=nil{
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("roadMap error=%s", err.Error()))
		exit(260, err)
	}
	//post路單server
	//取所有room
	//hr,_:=rCtrl.GetHistoryResult(rooms[0])
	//getHistoryResultPostData(rooms[0].HallID(), rooms[0].ID(), rooms[0].Type(), hr)


	//processor
	lp, err := command.NewLoginProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(300, err)
	}

	cr, err := command.NewClientReadyProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(301, err)
	}

	hi, err := command.NewHallInfoProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(302, err)
	}

	ui, err := command.NewUserInfoProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(303, err)
	}

	eh, err := command.NewEnterHallProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(304, err)
	}

	ri, err := command.NewRoomInfoProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(305, err)
	}

	er, err := command.NewEnterRoomProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(306, err)
	}

	lr, err := command.NewLeaveRoomProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(307, err)
	}

	lh, err := command.NewLeaveHallProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(308, err)
	}

	lg, err := command.NewLogoutProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(309, err)
	}

	hb, err := command.NewHeartbeatProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(310, err)
	}

	st, err := command.NewServerTimeProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(311, err)
	}

	//type0
	bt0, err := command.NewBetType0Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(320, err)
	}
	hrt0, err := command.NewHistoryResultType0Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(321, err)
	}

	//typ1
	bt1, err := command.NewBetType1Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(330, err)
	}
	hrt1, err := command.NewHistoryResultType1Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(331, err)
	}
	//typ1
	bt2, err := command.NewBetType2Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(340, err)
	}
	hrt2, err := command.NewHistoryResultType2Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(341, err)
	}

	//type6
	bt6, err := command.NewBetType6Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(350, err)
	}
	hrt6, err := command.NewHistoryResultType6Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(351, err)
	}

	//type˙
	bt7, err := command.NewBetType7Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(360, err)
	}
	hrt7, err := command.NewHistoryResultType7Processor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(361, err)
	}

	//err	pro := NewBasicProcessor(nex, conf)

	err = nex.RegisterCommandProcessor(conf.CmdLogin(), lp)
	if err != nil {
		exit(400, err)
	}

	err = nex.RegisterCommandProcessor(conf.CmdClientReady(), cr)
	if err != nil {
		exit(401, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdHallInfo(), hi)
	if err != nil {
		exit(402, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdUserInfo(), ui)
	if err != nil {
		exit(403, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdEnterHall(), eh)
	if err != nil {
		exit(404, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdRoomInfo(), ri)
	if err != nil {
		exit(405, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdEnterRoom(), er)
	if err != nil {
		exit(406, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdLeaveRoom(), lr)
	if err != nil {
		exit(407, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdLeaveHall(), lh)
	if err != nil {
		exit(408, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdHeartbeat(), hb)
	if err != nil {
		exit(409, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdLogout(), lg)
	if err != nil {
		exit(410, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdServerTime(), st)
	if err != nil {
		exit(411, err)
	}

	//type0
	err = nex.RegisterCommandProcessor(conf.CmdHistoryResultType0(), hrt0)
	if err != nil {
		exit(412, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdBetType0(), bt0)
	if err != nil {
		exit(413, err)
	}

	//type1
	err = nex.RegisterCommandProcessor(conf.CmdHistoryResultType1(), hrt1)
	if err != nil {
		exit(414, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdBetType1(), bt1)
	if err != nil {
		exit(415, err)
	}

	//type2
	err = nex.RegisterCommandProcessor(conf.CmdHistoryResultType2(), hrt2)
	if err != nil {
		exit(416, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdBetType2(), bt2)
	if err != nil {
		exit(417, err)
	}

	//type6
	err = nex.RegisterCommandProcessor(conf.CmdHistoryResultType6(), hrt6)
	if err != nil {
		exit(418, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdBetType6(), bt6)
	if err != nil {
		exit(419, err)
	}

	//type7
	err = nex.RegisterCommandProcessor(conf.CmdHistoryResultType7(), hrt7)
	if err != nil {
		exit(420, err)
	}
	err = nex.RegisterCommandProcessor(conf.CmdBetType7(), bt7)
	if err != nil {
		exit(421, err)
	}

	//builtin event
	ule, err := builtinevent.NewUserLostEventProcessor(command.NewBasicProcessor(nex, conf, dbCtrl,rmc))
	if err != nil {
		exit(422, err)
	}
	err = nex.RegisterBuiltinEventProcessor(nxBuiltinEvent.EventUserLost, ule)
	if err != nil {
		exit(423, err)
	}

	//r00 := fake.CreateRoom(nex, rCtrl, fake.CreateRoomField(1, "百家樂1", conf.RoomType0()))
	//r0 := fake.CreateRoom(nex, rCtrl, fake.CreateRoomField(1000, "龍虎1", conf.RoomType1()))
	//r10 := fake.CreateRoom(nex, rCtrl, fake.CreateRoomField(6000, "骰寶1", conf.RoomType6()))
	//r30 := fake.CreateRoom(nex, rCtrl, fake.CreateRoomField(2000, "牛牛1", conf.RoomType2()))
	//r40 := fake.CreateRoom(nex, rCtrl, fake.CreateRoomField(7000, "輪盤1", conf.RoomType7()))

	/*
		if rCtrl.NeedInitBootRound(r00) {
			rCtrl.SetRound(r00, 180817100100028)
		}
		if rCtrl.NeedInitBootRound(r0) {
			rCtrl.SetRound(r0, 180817200100048)
		}
		if rCtrl.NeedInitBootRound(r10) {
			rCtrl.SetRound(r10, 180817203100003)
		}
		if rCtrl.NeedInitBootRound(r30) {
			rCtrl.SetRound(r30, 180817203100003)
		}
		if rCtrl.NeedInitBootRound(r40) {
			rCtrl.SetRound(r40, 180817203100003)
		}

		h0.AddRoom(r00)
		h0.AddRoom(r0)
		h0.AddRoom(r10)
		h0.AddRoom(r30)
		h0.AddRoom(r40)
	*/

	//grpc server
	//addr, port := nex.GetConfig().RPCServerAddress()
	//err=nex.StartGRPCServer(pb.RegisterRPCServer, rpc.NewRPCServer(nex, conf, roomCtrl.NewController(conf), control.NewRateController(conf), control.NewDBController(conf.DBAPIServer())))
	err=nex.StartGRPCServer(pb.RegisterRPCServer, rpc.NewRPCServer(nex, conf, rCtrl, control.NewRateController(conf), dbCtrl,rmc))
	if err!=nil{
		exit(500, err)
	}
	nex.Start()

	//go profileStatus()
	//profileUtil.ProfileCPU("cpuprofile.prof")
}

func exit(id int, err error) {
	fmt.Println(err)
	log.Println(err)
	os.Exit(id)
}

//func profileStatus() {
//	http://localhost:6060/debug/pprof/  to see data
//	log.Println(http.ListenAndServe("localhost:6060", nil))
//}

func getConfigFilePosition(fileName string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}

	var buf bytes.Buffer
	buf.WriteString(dir)
	buf.WriteString("/")
	buf.WriteString(fileName)

	return buf.String()
}

//func poolUpdate(requester Requester, conf *config.Configurer, duration time.Duration) chan bool {
//	ticker := time.NewTicker(duration)
//	stop := make(chan bool, 1)
//
//	bp := &config.BackstageGetPoolParam{}
//
//	code := 0
//	err := errors.New("")
//	reason := ""
//	poolRes := config.BackstageGetPoolResData{}
//	go func() {
//		for {
//			select {
//			case <-ticker.C:
//				poolRes, code, reason, err = requester.GetPoolRequest(config.PathGetPool, bp)
//				if code == config.CodeSuccess && err == nil {
//					conf.SetJackpot(poolRes.Jackpot)
//					conf.SetReturn(poolRes.Return)
//					conf.SetJackpotReserved(poolRes.JackpotReserved)
//				}
//			case <-stop:
//				return
//			}
//		}
//	}()
//
//	return stop
//}

//func poolSettingUpdate(requester Requester, conf Configurer, duration time.Duration) chan bool {
//	ticker := time.NewTicker(duration)
//	stop := make(chan bool, 1)
//
//	bp := &config.BackstageGetPoolSettingParam{}
//
//	code := 0
//	err := errors.New("")
//	reason := ""
//	poolRes := config.BackstageGetPoolSettingResData{}
//	go func() {
//		for {
//			select {
//			case <-ticker.C:
//				poolRes, code, reason, err = requester.GetPoolSettingRequest(config.PathGetPoolSetting, bp)
//				if code == config.CodeSuccess && err == nil {
//					conf.SetJackpotRatio(poolRes.Jackpot)
//					conf.SetReturnRatio(poolRes.Return)
//					conf.SetJackpotReservedRatio(poolRes.JackpotReserved)
//					conf.SetPartnerKillRatio(poolRes.Kill)
//				}
//			case <-stop:
//				return
//			}
//		}
//	}()
//
//	return stop
//}



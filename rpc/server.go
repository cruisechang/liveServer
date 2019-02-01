package rpc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	pb "github.com/cruisechang/liveServer/protobuf"

	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cruisechang/liveServer/config/dbConf"

	"github.com/cruisechang/nex"
	nxHTTP "github.com/cruisechang/nex/http"
	nxLog "github.com/cruisechang/nex/log"
	"golang.org/x/net/context"
)

func NewRPCServer(nex nex.Nex, conf config.Configurer, controller *roomCtrl.Controller, rateController *control.RateController, dbController *control.DBController, roadMapCtrl *control.RoadMapController) *rpcServer {
	r := &rpcServer{
		nex:         nex,
		configure:   conf,
		roomCtrl:    controller,
		rateCtrl:    rateController,
		dbCtrl:      dbController,
		roadMapCtrl: roadMapCtrl,
	}

	nexConf := nex.GetConfig()
	address, port := nexConf.HttpClientAddress()

	cl, _ := nxHTTP.NewClient(address, port, nexConf.HttpClientTCPConnectTimeoutSecond(), nexConf.HttpClientHandshakeTimeoutSecond(), nexConf.HttpClientRequestTimeoutSecond())

	cl.SetHost("http", address, port)
	r.httpClient = cl

	return r
}

// server is used to implement helloworld.GreeterServer.
type rpcServer struct {
	nex         nex.Nex
	configure   config.Configurer
	roomCtrl    *roomCtrl.Controller
	rateCtrl    *control.RateController
	dbCtrl      *control.DBController
	roadMapCtrl *control.RoadMapController
	httpClient  nxHTTP.Client
}

func (s *rpcServer) DealerLogin(ctx context.Context, in *pb.DealerLoginData) (*pb.DealerLoginRes, error) {
	logPrefix := "rpcServer DealerLogin"
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	dealerID := int(in.Dealer)
	dID := strconv.FormatInt(in.Dealer, 10)

	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s roomID not found =%d", logPrefix, in.RoomID))
		return nil, fmt.Errorf("%s roomID not found", logPrefix)
	}

	//check password  不需要

	//check dealer active
	path := fmt.Sprintf("/dealers/" + dID)
	resp, err := s.dbCtrl.Do("GET", path, nil)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db response error:%s, dealerID=%d", logPrefix, err.Error(), dealerID))
		return nil, fmt.Errorf("%s db error", logPrefix)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db http status code wrong=%d ,error:%s,dealerID=%d", logPrefix, resp.StatusCode, err.Error(), dealerID))
		return nil, fmt.Errorf("%s db http status code wrong=%d", logPrefix, resp.StatusCode)
	}

	if resp.Body == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db response body==nil,error:%s,dealerID=%d", logPrefix, err.Error(), dealerID))
		return nil, fmt.Errorf("%s db response body==nil,error:%s,dealerID=%d", logPrefix, err.Error(), dealerID)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	//db got
	got := &dbConf.DealerData{}
	err = json.Unmarshal(body, got)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s unmarshal db got body error:%s", logPrefix, err.Error()))
		return nil, fmt.Errorf("%s unmarshal db got body error=%s", logPrefix, err.Error())
	}

	if got.Code != config.CodeSuccess {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db got code=%d,dealerID=%d", logPrefix, got.Code, dealerID))
		return nil, fmt.Errorf("%s db got code code=%d ,dealerID=%d", logPrefix, got.Code, dealerID)
	}
	if got.Count != 1 {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db got error count=%d,dealer=%d", logPrefix, got.Count, dealerID))
		return nil, fmt.Errorf("%s db got error count=%d, dealerID=%d", logPrefix, got.Count, dealerID)
	}

	if len(got.Data) != 1 {
		if got.Data[0].Active != 1 {
			logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db got dealer inactive  dealer=%d", logPrefix, dealerID))
			return &pb.DealerLoginRes{Success: 0}, nil
		}
	}

	//DB API patch room
	path = fmt.Sprintf("/rooms/" + strconv.FormatInt(int64(r.ID()), 10) + "/dealerID")

	bo, err := json.Marshal(struct {
		DealerID int `json:"dealerID"`
	}{dealerID})
	if err != nil {
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s json marshal body erro=%s", logPrefix, err.Error()))
		return nil, fmt.Errorf("%s db json marshal error=%s, dealerID=%d", logPrefix, err.Error(), dealerID)
	}
	resp2, err := s.dbCtrl.Do("PATCH", path, bytes.NewBuffer(bo))

	if err != nil {
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s db patch err=%s", logPrefix, err.Error()))
		return nil, fmt.Errorf("%s db patch err=%s", logPrefix, err.Error())
	}
	defer resp2.Body.Close()

	//set room
	r.SetStatus(s.configure.RoomStatusChangeDealer())

	//send client
	receivers := getHallReceivers(s.nex.GetUsers(), r.HallID())

	//send to client
	if len(receivers) > 0 {
		resData := []config.ChangeDealerResData{
			{
				HallID:      1,
				RoomID:      int(in.RoomID),
				DealerID:    dealerID,
				Name:        "",
				PortraitURL: "",
			},
		}

		b, _ := json.Marshal(resData)
		sendDataStr := base64.StdEncoding.EncodeToString(b)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdChangeDealer(), sendDataStr, nil, receivers)
		s.nex.GetLogger().LogFile(nxLog.LevelInfo, fmt.Sprintf("%s toClient=%s", logPrefix, string(b)))

		//roomStatus
		data := s.roomCtrl.GetBootRoundBetMinBetMax(r)
		b, _ = getRoomStatusResponse(r, s.configure.RoomStatusChangeDealer(), data.Boot, data.Round, getStatusStart())
		sendDataStr = base64.StdEncoding.EncodeToString(b)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoomStatus(), sendDataStr, nil, receivers)
		s.nex.GetLogger().LogFile(nxLog.LevelInfo, fmt.Sprintf("%s roomStatus=%s", logPrefix, string(b)))

	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete", logPrefix))

	return &pb.DealerLoginRes{
		Success: 1,
	}, nil
}

func (s *rpcServer) RoomLogin(ctx context.Context, in *pb.RoomLoginData) (*pb.Empty, error) {
	logPrefix := "rpcServer RoomLogin"
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	//in.RoomID
	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s roomID not found :%d", logPrefix, in.RoomID))
		return nil, errors.New("RoomLogin roomID not found")
	}
	r.SetActive(1)

	//send client
	receivers := getHallReceivers(s.nex.GetUsers(), r.HallID())

	//send to client
	if len(receivers) > 0 {
		resData := []config.RoomActiveResData{
			{
				RoomID: int(in.RoomID),
				Active: s.configure.Active(),
			},
		}

		b, _ := json.Marshal(resData)
		sendDataStr := base64.StdEncoding.EncodeToString(b)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoomActive(), sendDataStr, nil, receivers)
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s roomActive =%s", logPrefix, string(b)))
	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete", logPrefix))

	return &pb.Empty{}, nil
}
func (s *rpcServer) OnlineNotify(ctx context.Context, in *pb.OnlineNotifyData) (*pb.Empty, error) {
	logPrefix := "rpcServer OnlineNotify"
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()
	//in.RoomID

	return &pb.Empty{}, nil
}

//GetRoomInfo gets info of target room
func (s *rpcServer) GetRoomInfo(ctx context.Context, in *pb.GetRoomInfoData) (*pb.GetRoomInfoRes, error) {
	logPrefix := "rpcServer GetRoomInfo"
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v\n", logPrefix, r))
		}
	}()

	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s roomID not found :%d", logPrefix, int(in.RoomID)))
		return nil, errors.New("rpcServer GetRoomInfo roomID not found")
	}

	brData := s.roomCtrl.GetBootRoundBetMinBetMax(r)
	if brData.Err != nil {
		//error
		return nil, brData.Err
	}

	//Get typeData by room type
	return &pb.GetRoomInfoRes{
		Boot:                int64(brData.Boot),
		Round:               brData.Round,
		RoomID:              int64(r.ID()),
		RoomName:            r.Name(),
		BankerPlayerMin:     int64(brData.BetMin),
		BankerPlayerPairMax: int64(brData.BetMax),
		TieMin:              int64(brData.BetMin),
		TieMax:              int64(brData.BetMax),
		BankerPlayerPairMin: int64(brData.BetMin),
		BankerPlayerMax:     int64(brData.BetMax),
		Online:              999,
		BetCountDown:        20,
	}, nil
}

func (s *rpcServer) ChangeBoot(ctx context.Context, in *pb.ChangeBootData) (*pb.ChangeBootRes, error) {
	logPrefix := "rpcServer ChangeBoot"
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v\n", logPrefix, r))
		}
	}()

	//in.RoomID
	r, ok := s.nex.GetRoomManager().GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s roomID not found :%d", logPrefix, in.RoomID))
		return &pb.ChangeBootRes{}, fmt.Errorf("%s roomID not found", logPrefix)
	}

	//更改room 狀態
	r.SetStatus(s.configure.RoomStatusChangeBoot())
	newBoot, newRound := s.roomCtrl.HandleNewBoot(r)
	s.roomCtrl.SetHistoryResultEmpty(r)

	//send client
	receivers := getRoomReceivers(s.nex.GetUsers(), r.ID())

	//send to client
	if len(receivers) > 0 {

		//roomStatus
		b, _ := getRoomStatusResponse(r, s.configure.RoomStatusChangeBoot(), newBoot, newRound, getStatusStart())
		sendDataStr := base64.StdEncoding.EncodeToString(b)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoomStatus(), sendDataStr, nil, receivers)
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s roomStatus =%s", logPrefix, string(b)))

	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete", logPrefix))

	return &pb.ChangeBootRes{
		Boot: int64(newBoot),
	}, nil
}

func (s *rpcServer) Waiting(ctx context.Context, in *pb.WaitingData) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer Waiting panic:%v", r))
		}
	}()

	//in.RoomID
	//in.Round
	//in.Boot
	//r, ok := s.nex.GetRoomManager().GetRoom(int(in.RoomID))
	//if !ok {
	//	return &pb.Empty{
	//	}, errors.New("roomID not found")
	//}
	//
	//r.SetStatus(s.configure.RoomStatusStopBet())

	return &pb.Empty{}, nil
}

func (s *rpcServer) BeginBetting(ctx context.Context, in *pb.BeginBettingData) (*pb.BeginBettingRes, error) {
	logPrefix := "rpcServer BeginBetting"
	logger := s.nex.GetLogger()
	conf := s.configure

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	//in.RoomID
	//in.Round
	//in.Boot
	//改room status
	r, ok := s.nex.GetRoomManager().GetRoom(int(in.RoomID))
	if !ok {
		return &pb.BeginBettingRes{}, errors.New("roomID not found")
	}

	statusStart := getStatusStart()
	r.SetStatus(s.configure.RoomStatusBeginBet())
	r.SetStatusStart(statusStart)

	boot, round := s.roomCtrl.HandleNewRound(r)

	s.roomCtrl.InitRoomBet(r)

	receivers := getHallReceivers(s.nex.GetUsers(), r.HallID())

	//send client
	if len(receivers) > 0 {

		b, _ := getRoomStatusResponse(r, s.configure.RoomStatusBeginBet(), boot, round, statusStart)
		sendDataStr := base64.StdEncoding.EncodeToString(b)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoomStatus(), sendDataStr, nil, receivers)
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf(" roomStatus =%s", string(b)))
	}

	//DB API patch room boot/round/status
	path := fmt.Sprintf("/rooms/" + strconv.FormatInt(int64(r.ID()), 10) + "/newRound")

	bo, err := json.Marshal(dbConf.RoomNewRoundData{Boot: boot, RoundID: round, Status: conf.RoomStatusBeginBet()})
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json marshal body erro=%s", logPrefix, err.Error()))
		//return nil, fmt.Errorf("%s db json marshal error=%s, roomID=%d", logPrefix, err.Error(), r.ID())
	}
	resp, err := s.dbCtrl.Do("PATCH", path, bytes.NewBuffer(bo))

	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db patch err=%s", logPrefix, err.Error()))
	} else {
		defer resp.Body.Close()
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete", logPrefix))
	}

	return &pb.BeginBettingRes{
		Round: round,
	}, nil
}

func (s *rpcServer) EndBetting(ctx context.Context, in *pb.EndBettingData) (*pb.Empty, error) {

	logPrefix := "rpcServer EndBetting"
	logger := s.nex.GetLogger()
	conf := s.configure

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	//in.RoomID
	r, ok := s.nex.GetRoomManager().GetRoom(int(in.RoomID))
	if !ok {
		return &pb.Empty{}, errors.New("roomID not found")
	}

	statusStart := getStatusStart()
	r.SetStatus(s.configure.RoomStatusEndBet())
	r.SetStatusStart(statusStart)

	receivers := getHallReceivers(s.nex.GetUsers(), r.HallID())

	//logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer endBetting complete  user len=%d, connID=%s \n", len(receivers), receivers[0]))

	if len(receivers) > 0 {

		brData := s.roomCtrl.GetBootRoundBetMinBetMax(r)

		b, _ := getRoomStatusResponse(r, s.configure.RoomStatusEndBet(), brData.Boot, brData.Round, statusStart)
		sendDataStr := base64.StdEncoding.EncodeToString(b)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoomStatus(), sendDataStr, nil, receivers)
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s roomStatus =%s", logPrefix, string(b)))
	}

	//DB API patch room status
	path := fmt.Sprintf("/rooms/" + strconv.FormatInt(int64(r.ID()), 10) + "/status")

	bo, err := json.Marshal(dbConf.Status{Status: conf.RoomStatusEndBet()})
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json marshal body erro=%s", logPrefix, err.Error()))
		//return nil, fmt.Errorf("%s db json marshal error=%s, roomID=%d", logPrefix, err.Error(), r.ID())
	}
	resp, err := s.dbCtrl.Do("PATCH", path, bytes.NewBuffer(bo))

	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db patch err=%s", logPrefix, err.Error()))
	} else {
		defer resp.Body.Close()
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete", logPrefix))
	}

	return &pb.Empty{}, nil
}

func (s *rpcServer) CancelRound(ctx context.Context, in *pb.CancelRoundData) (*pb.Empty, error) {
	logPrefix := "rpcServer CancelRound"
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	//取room manager
	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s roomID not found :%d", logPrefix, int(in.RoomID)))
		return nil, fmt.Errorf("%s rpcServer CancelRound roomID not found", logPrefix)
	}

	//check round
	brData := s.roomCtrl.GetBootRoundBetMinBetMax(r)
	if brData.Round != in.Round {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s round not match roomID=%d, get=%d, want=%d", logPrefix, int(in.RoomID), in.Round, brData.Round))
		return nil, fmt.Errorf("%s rpcServer CancelRound round not match", logPrefix)
	}

	//清除下注
	s.roomCtrl.InitRoomBet(r)

	//round record
	//

	//send client
	receivers := getRoomReceivers(s.nex.GetUsers(), r.ID())
	if len(receivers) > 0 {

		resData := []config.CancelRoundResData{
			{
				HallID: r.HallID(),
				RoomID: r.ID(),
				Boot:   brData.Boot,
				Round:  brData.Round,
			},
		}

		b, _ := json.Marshal(resData)
		sendDataStr := base64.StdEncoding.EncodeToString(b)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdCancelRound(), sendDataStr, nil, receivers)
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s toClient=%s", logPrefix, string(b)))
	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete", logPrefix))

	return &pb.Empty{}, nil
}

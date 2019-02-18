package rpc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	pb "github.com/cruisechang/liveServer/protobuf"

	"bytes"
	"strconv"

	"github.com/cruisechang/nex/entity"
	nxLog "github.com/cruisechang/nex/log"
)

func newRoundResultHandler(rpc *rpcServer, logPrefix string, roomID int32, round int64, in interface{}) *roundResultHandler {
	return &roundResultHandler{
		rpc:       rpc,
		logPrefix: logPrefix,
		roomID:    roomID,
		round:     round,
		in:        in,
	}
}

type roundResultHandler struct {
	logPrefix string
	rpc       *rpcServer
	roomID    int32
	round     int64
	in        interface{}
	brData    *roomCtrl.BootRoundData
	room      entity.Room
}

//1
//post round record  to db
//send round result to client

//2
//玩家投注

//3
//room add result
//room get new result

//4
//get road map request data
//post road map server
//set road map data from result
//send road map data to client
func (h *roundResultHandler) handle() (*pb.Empty, error) {
	logger := h.rpc.nex.GetLogger()
	conf := h.rpc.configure
	rm := h.rpc.nex.GetRoomManager()
	roadMapCtrl := h.rpc.roadMapCtrl

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", h.logPrefix, r))
		}
	}()

	r, ok := rm.GetRoom(int(h.roomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s room not found roomID=%d", h.logPrefix, h.roomID))
		return nil, errors.New("room not found")
	}
	h.room = r
	h.brData = h.rpc.roomCtrl.GetBootRoundBetMinBetMax(h.room)

	//送給hall所有人牌局結果
	hallReceivers := getHallReceivers(h.rpc.nex.GetUsers(), h.room.HallID())

	/////1
	//post round result to db api server
	h.postRoundResultToDBAPI(h.logPrefix, logger, h.room, h.in, h.brData, getRoundRecordData)

	//send to client
	if len(hallReceivers) > 0 {
		cmdStr, _ := h.getRoundResultCommandString(conf, h.room.Type())
		resultToClient, err := getRoundResultResponse(h.in, h.brData)

		if err != nil {
			logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s getRoundResultResponse error :%s", h.logPrefix, err.Error()))
		} else {
			sendDataStr := base64.StdEncoding.EncodeToString(resultToClient)
			h.rpc.sendCommand(config.CodeSuccess, 0, cmdStr, sendDataStr, nil, hallReceivers)
			logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s send result  =%s", h.logPrefix, string(resultToClient)))
		}

		//roomStatus
		b, err := getRoomStatusResponse(h.room, conf.RoomStatusRoundResult(), h.brData.Boot, h.round, getStatusStart())
		if err != nil {
			logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s getRoomStatusResponse  error=%s", h.logPrefix, err.Error()))
		} else {
			sendDataStr := base64.StdEncoding.EncodeToString(b)
			h.rpc.sendCommand(config.CodeSuccess, 0, conf.CmdRoomStatus(), sendDataStr, nil, hallReceivers)
			logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s send roomStatus =%s", h.logPrefix, string(b)))
		}
	}

	/////2
	//取玩家投注資料
	//存投注資料
	bd, err := h.rpc.roomCtrl.GetBet(h.room)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s roomCtrl.GetBet error =%s", h.logPrefix, err.Error()))
	} else {
		//投注人數
		ln, err := h.countBetLength(bd)
		if err != nil {
			logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s bet data assertion error =%s", h.logPrefix, err.Error()))

		} else if ln > 0 {
			logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s has bet", h.logPrefix))

			//有人投注  計算所有投注的人結果
			userRes := []config.UserResultResData{
				{
					RoomID: int(h.roomID),
					Round:  h.brData.Round,
					Result: []config.UserResultResDataResult{},
				},
			}

			//轉型 為了range
			//betData, _ := bd.(map[int]interface{})

			//算錢, patch user credit, post bet
			//var pid int64
			//for uid, v := range betData {
			//	if err=h.handleBet(uid, userRes, h.room, v, h.brData);err!=nil{
			//		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s %s", h.logPrefix, err.Error()))
			//	}
			//}
			if err, errMsg := h.handleBet(bd, userRes, h.room, h.brData); err != nil {
				errStrings := strings.Join(errMsg, " # ")
				logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s %s err strings=%s", h.logPrefix, err.Error(), errStrings))
			}

			//user result
			if len(userRes[0].Result) > 0 {
				if roomReceivers := getRoomReceivers(h.rpc.nex.GetUsers(), h.room.ID()); len(roomReceivers) > 0 {
					b, _ := json.Marshal(userRes)
					sendDataStr := base64.StdEncoding.EncodeToString(b)
					h.rpc.sendCommand(config.CodeSuccess, 0, conf.CmdUserResult(), sendDataStr, nil, roomReceivers)
					logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s userResult =%s", h.logPrefix, string(b)))
				}
			}
		}
	}

	/////3
	//room add result
	err = h.addHistoryResult(h.logPrefix, h.room, h.in)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s addHistoryResult error=%s", h.logPrefix, err.Error()))
	}
	hrs, err := h.rpc.roomCtrl.GetHistoryResult(h.room)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s get history result error=%s ", h.logPrefix, err.Error()))
	}

	//////4 road map
	postData, err := roadMapCtrl.GetRoadMapRequestData(h.room.HallID(), h.room.ID(), h.room.Type(), hrs)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s get history result error=%s ", h.logPrefix, err.Error()))
	}

	//post road map api server
	body, err := roadMapCtrl.RequestRoadMap(postData)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s RequestRoadMap() error=%s", h.logPrefix, err.Error()))
		return &pb.Empty{}, err
	}
	resData := &config.RoadMapResponse{}

	err = json.Unmarshal(body, resData)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s RequestRoadMap() RoadMapResponse Unmarshal() error=%s", h.logPrefix, err.Error()))
	}

	if resData.Code != config.CodeSuccess {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s RequestRoadMap RoadMapResponse resData code error code=%d", h.logPrefix, resData.Code))
	}

	err = roadMapCtrl.SetRoadMapDataFromResult(h.room.ID(), h.room.Type(), resData.Result)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s SetRoadMapDataFromResult() error=%s", h.logPrefix, err.Error()))
	}
	//send to client
	rmToClient, err := h.getRoomMapData(h.room.Type())
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s getRoomMapData() error=%s", h.logPrefix, err.Error()))
	}
	b, err := json.Marshal(rmToClient)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s getRoomMapData() json marshal roadmap error=%s", h.logPrefix, err.Error()))
	}

	sendData := base64.StdEncoding.EncodeToString(b)
	rmCmd, _ := h.getRoadMapCommand(h.room.Type())

	h.rpc.sendCommand(config.CodeSuccess, 0, rmCmd, sendData, nil, hallReceivers)
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("roadMapType0 resData=%s ", string(b)))

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete", h.logPrefix))

	return &pb.Empty{}, nil
}
func (h *roundResultHandler) handleBet(betData interface{}, userRes []config.UserResultResData, room entity.Room, brData *roomCtrl.BootRoundData) (error, []string) {

	var reErr error
	errMsg := []string{}

	if bd, ok := betData.(map[int]*roomConf.BetType0Data); ok {
		for uid, v := range bd {
			if err := h.handleBetProcess(uid, userRes, room, v, brData); err != nil {
				errMsg = append(errMsg, err.Error())
			}
		}

		if len(errMsg) > 0 {
			reErr = fmt.Errorf(" handlerBet")
		}
		return reErr, errMsg
	}
	if bd, ok := betData.(map[int]*roomConf.BetType1Data); ok {
		for uid, v := range bd {
			if err := h.handleBetProcess(uid, userRes, room, v, brData); err != nil {
				errMsg = append(errMsg, err.Error())
			}
		}
		if len(errMsg) > 0 {
			reErr = fmt.Errorf(" handlerBet")
		}
		return reErr, errMsg
	}
	if bd, ok := betData.(map[int]*roomConf.BetType2Data); ok {
		for uid, v := range bd {
			if err := h.handleBetProcess(uid, userRes, room, v, brData); err != nil {
				errMsg = append(errMsg, err.Error())
			}
		}
		if len(errMsg) > 0 {
			reErr = fmt.Errorf(" handlerBet")
		}
		return reErr, errMsg
	}
	if bd, ok := betData.(map[int]*roomConf.BetType6Data); ok {
		for uid, v := range bd {
			if err := h.handleBetProcess(uid, userRes, room, v, brData); err != nil {
				errMsg = append(errMsg, err.Error())
			}
		}
		if len(errMsg) > 0 {
			reErr = fmt.Errorf(" handlerBet")
		}
		return reErr, errMsg
	}
	if bd, ok := betData.(map[int]*roomConf.BetType7Data); ok {
		for uid, v := range bd {
			if err := h.handleBetProcess(uid, userRes, room, v, brData); err != nil {
				errMsg = append(errMsg, err.Error())
			}
		}
		if len(errMsg) > 0 {
			reErr = fmt.Errorf(" handlerBet")
		}
		return reErr, errMsg
	}
	errMsg = append(errMsg, "*roomConf.BetType data not found")
	return fmt.Errorf(" handlerBet"), errMsg

}

func (h *roundResultHandler) handleBetProcess(uid int, userRes []config.UserResultResData, room entity.Room, v interface{}, brData *roomCtrl.BootRoundData) error {

	if user, ok := h.rpc.nex.GetUser(uid); ok {

		oriCredit := user.Credit()
		betCredit, activeCredit, prizeCredit, resultCredit, balanceCredit, winLose, err := h.rpc.rateCtrl.Count(room.Type(), user.Credit(), v, h.in)

		userRes[0].Result = append(userRes[0].Result, getUserResResult(int64(uid), winLose, resultCredit, balanceCredit))

		//set user data
		user.SetCredit(balanceCredit)
		DBUserID, err := user.GetInt64Variable(h.rpc.configure.UserVarDBUserID())
		if err != nil {
			h.rpc.nex.GetLogger().LogFile(nxLog.LevelError, fmt.Sprintf("%s post bet get user db userID error=%s\n", h.logPrefix, err.Error()))
			return fmt.Errorf(" handlerBet %s", err.Error())
		}

		//patch credit
		b, _ := json.Marshal(struct{ Credit float32 }{balanceCredit})
		path := fmt.Sprintf("/users/" + strconv.FormatInt(DBUserID, 10) + "/credit")
		//h.rpc.httpDo("PATCH", path, bytes.NewBuffer(b))
		res, err := h.rpc.dbCtrl.Do("PATCH", path, bytes.NewBuffer(b))
		if err != nil {
			return fmt.Errorf(" handleBet %s", err.Error())
		}
		defer res.Body.Close()

		//post bet data
		rrb, err := json.Marshal(v)
		if err != nil {
			h.rpc.nex.GetLogger().LogFile(nxLog.LevelError, fmt.Sprintf("%s post bet json marshal betData error=%s\n", h.logPrefix, err.Error()))
			return fmt.Errorf(" handleBet %s", err.Error())
		}
		pid, _ := user.GetInt64Variable(h.rpc.configure.UserVarPartnerID())

		//bet record
		br := &config.BetPostParam{
			pid,
			DBUserID,
			room.ID(),
			room.Type(),
			brData.Round,
			0,
			betCredit,
			activeCredit,
			prizeCredit,
			resultCredit,
			balanceCredit,
			oriCredit,
			string(rrb),
			1,
		}

		//
		bb, err := json.Marshal(br)
		if err != nil {
			h.rpc.nex.GetLogger().LogFile(nxLog.LevelError, fmt.Sprintf("%s post bet json marshal  error=%s\n", h.logPrefix, err.Error()))
			return fmt.Errorf(" handleBet %s", err.Error())
		}
		res2, err2 := h.rpc.dbCtrl.Do("POST", "/bets", bytes.NewBuffer(bb))
		if err2 != nil {
			return fmt.Errorf("  handleBet %s", err.Error())
		}
		defer res2.Body.Close()

		h.rpc.nex.GetLogger().LogFile(nxLog.LevelInfo, fmt.Sprintf("%s handleBet post bets complete dbUserID=%d ,param=%s", h.logPrefix, DBUserID, bb))

		return nil
	}

	return fmt.Errorf(" handleBet h.rpc.nex.GetUser() error ,userID=%d", uid)
}

func (h *roundResultHandler) countBetLength(bd interface{}) (int, error) {

	if betData, ok := bd.(map[int]*roomConf.BetType0Data); ok {
		return len(betData), nil
	}
	if betData, ok := bd.(map[int]*roomConf.BetType1Data); ok {
		return len(betData), nil
	}
	if betData, ok := bd.(map[int]*roomConf.BetType2Data); ok {
		return len(betData), nil
	}
	if betData, ok := bd.(map[int]*roomConf.BetType6Data); ok {
		return len(betData), nil
	}
	if betData, ok := bd.(map[int]*roomConf.BetType7Data); ok {
		return len(betData), nil
	}

	return 0, errors.New(" not match type")

}

func (h *roundResultHandler) addHistoryResult(logPrefix string, room entity.Room, data interface{}) error {

	rc := h.rpc.roomCtrl
	if in, ok := data.(*pb.RoundResultType0Data); ok {
		return rc.AddHistoryResultType(room, []int32{in.Result, in.BankerPair, in.PlayerPair})
	}

	if in, ok := data.(*pb.RoundResultType1Data); ok {
		return rc.AddHistoryResultType(room, in.Result)
	}

	if in, ok := data.(*pb.RoundResultType2Data); ok {
		return rc.AddHistoryResultType(room, []int32{in.Owner0.Result, in.Owner0.Pattern, in.Owner1.Result, in.Owner1.Pattern, in.Owner2.Result, in.Owner2.Pattern, in.Owner3.Result, in.Owner3.Pattern})
	}

	if in, ok := data.(*pb.RoundResultType6Data); ok {
		convertedDice := []int{}
		for _, dv := range in.Dice {
			convertedDice = append(convertedDice, int(dv))
		}

		//room result 增加
		hrt6 := &roomConf.HistoryResultType6{
			HallID:   room.HallID(),
			RoomID:   room.ID(),
			Dice:     convertedDice,
			Sum:      int(in.Sum),
			BigSmall: int(in.BigSmall),
			OddEven:  int(in.OddEven),
		}
		return rc.AddHistoryResultType(room, hrt6)
	}

	if in, ok := data.(*pb.RoundResultType7Data); ok {
		return rc.AddHistoryResultType(room, in.Result)
	}

	return fmt.Errorf("%s addHistoryFunc data assertion error=", logPrefix)
}

//post round result to db api
func (h *roundResultHandler) postRoundResultToDBAPI(logPrefix string, logger nxLog.Logger, room entity.Room, rpcIn interface{}, brData *roomCtrl.BootRoundData, getRoundRecordDataFunc func(in interface{}, brData *roomCtrl.BootRoundData) ([]byte, error)) {

	//round record
	rr, err := getRoundRecordDataFunc(rpcIn, brData)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s postRoundRecord getRoundRecordData error=%s\n", logPrefix, err.Error()))
	}

	rp := config.RoundPostParam{
		HallID:   room.HallID(),
		RoomID:   room.ID(),
		RoomType: room.Type(),
		Brief:    "",
		Record:   string(rr),
		Status:   1,
	}

	b, err := json.Marshal(rp)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s postRoundResult json marshal  error=%s\n", logPrefix, err.Error()))
	}
	resp, err := h.rpc.dbCtrl.Do("POST", "/rounds", bytes.NewBuffer(b))
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s post db   error=%s\n", logPrefix, err.Error()))
	}
	defer resp.Body.Close()

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s postRoundRecord complete body=%s\n", logPrefix, string(b)))
}

func (h *roundResultHandler) getRoundResultCommandString(conf config.Configurer, roomType int) (string, error) {

	if roomType == conf.RoomType0() {
		return conf.CmdRoundResult0(), nil
	}
	if roomType == conf.RoomType1() {
		return conf.CmdRoundResult1(), nil
	}
	if roomType == conf.RoomType2() {
		return conf.CmdRoundResult2(), nil
	}
	if roomType == conf.RoomType6() {
		return conf.CmdRoundResult6(), nil
	}
	if roomType == conf.RoomType7() {
		return conf.CmdRoundResult7(), nil
	}

	return "", errors.New(" not match type")
}

func (h *roundResultHandler) getRoomMapData(roomType int) (interface{}, error) {

	switch roomType {
	case h.rpc.configure.RoomType0():
		toClient := []*config.RoadMapType0ResData{}
		if got, ok := h.rpc.roadMapCtrl.RoadMapDataType0(h.room.ID()); ok {
			toClient = append(toClient, got)
		}
		return toClient, nil
	case h.rpc.configure.RoomType1():
		toClient := []*config.RoadMapType1ResData{}
		if got, ok := h.rpc.roadMapCtrl.RoadMapDataType1(h.room.ID()); ok {
			toClient = append(toClient, got)
		}
		return toClient, nil
	case h.rpc.configure.RoomType2():
		toClient := []*config.RoadMapType2ResData{}
		if got, ok := h.rpc.roadMapCtrl.RoadMapDataType2(h.room.ID()); ok {
			toClient = append(toClient, got)
		}
		return toClient, nil
	case h.rpc.configure.RoomType6():
		toClient := []*config.RoadMapType6ResData{}
		if got, ok := h.rpc.roadMapCtrl.RoadMapDataType6(h.room.ID()); ok {
			toClient = append(toClient, got)
		}
		return toClient, nil
	case h.rpc.configure.RoomType7():
		toClient := []*config.RoadMapType7ResData{}
		if got, ok := h.rpc.roadMapCtrl.RoadMapDataType7(h.room.ID()); ok {
			toClient = append(toClient, got)
		}
		return toClient, nil
	default:
		return nil, fmt.Errorf("room type not found=%d", roomType)
	}
}

func (h *roundResultHandler) getRoadMapCommand(roomType int) (string, error) {
	switch roomType {
	case h.rpc.configure.RoomType0():
		return h.rpc.configure.CmdRoadMapType0(), nil
	case h.rpc.configure.RoomType1():
		return h.rpc.configure.CmdRoadMapType1(), nil
	case h.rpc.configure.RoomType2():
		return h.rpc.configure.CmdRoadMapType2(), nil
	case h.rpc.configure.RoomType6():
		return h.rpc.configure.CmdRoadMapType6(), nil
	case h.rpc.configure.RoomType7():
		return h.rpc.configure.CmdRoadMapType7(), nil
	default:
		return "", fmt.Errorf("room type not found=%d", roomType)
	}
}

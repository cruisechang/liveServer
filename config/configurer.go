package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	active     int = 1
	inactive   int = 0
	unAssigned int = -1

	//room type
	roomTypeBaccarat    int = 0
	roomTypeDragonTiger int = 1
	roomTypeNiuniu      int = 2
	roomTypeFantan      int = 3
	roomTypeSangong     int = 4
	roomTypeThisBar     int = 5
	roomTypeSicbo       int = 6
	roomTypeRolette     int = 7

	roomStatusBeginBet     int = 1
	roomStatusEndBet       int = 2
	roomStatusRoundResult  int = 4
	roomStatusChangeBoot   int = 5
	roomStatusChangeDealer int = 6

	//room hall
	cmdHallInfo    string = "hallInfo"
	cmdEnterHall   string = "enterHall"
	cmdLeaveHall   string = "leaveHall"
	cmdRoomInfo    string = "roomInfo"
	cmdEnterRoom   string = "enterRoom"
	cmdLeaveRoom   string = "leaveRoom"
	cmdRoomStatus  string = "roomStatus"
	cmdCancelRound string = "cancelRound"

	cmdUserResult string = "userResult"

	//user behavior
	cmdLogin     string = "login"
	cmdLogout    string = "logout"
	cmdHeartbeat string = "heartbeat"
	cmdUserInfo  string = "userInfo"

	//support
	cmdUserCount    string = "userCount"
	cmdBroadcast    string = "broadcast"
	cmdBanner       string = "banner"
	cmdChangeDealer string = "changeDealer"
	cmdRoomActive   string = "roomActive"
	cmdHallActive   string = "hallActive"
	cmdServerTime   string = "serverTime"
	cmdClientReady  string = "clientReady"

	//
	cmdBetType0 = "betType0"
	cmdBetType1 = "betType1"
	cmdBetType2 = "betType2"
	cmdBetType3 = "betType3"
	cmdBetType4 = "betType4"
	cmdBetType5 = "betType5"
	cmdBetType6 = "betType6"
	cmdBetType7 = "betType7"
	//
	cmdRoundProcess0 = "roundProcess0"
	cmdRoundProcess1 = "roundProcess1"
	cmdRoundProcess2 = "roundProcess2"
	cmdRoundProcess3 = "roundProcess3"
	cmdRoundProcess4 = "roundProcess4"
	cmdRoundProcess5 = "roundProcess5"
	cmdRoundProcess6 = "roundProcess6"

	cmdRoundResultType0 = "roundResultType0"
	cmdRoundResultType1 = "roundResultType1"
	cmdRoundResultType2 = "roundResultType2"
	cmdRoundResultType3 = "roundResultType3"
	cmdRoundResultType4 = "roundResultType4"
	cmdRoundResultType5 = "roundResultType5"
	cmdRoundResultType6 = "roundResultType6"
	cmdRoundResultType7 = "roundResultType7"

	cmdHistoryResultType0 = "historyResultType0"
	cmdHistoryResultType1 = "historyResultType1"
	cmdHistoryResultType2 = "historyResultType2"
	cmdHistoryResultType3 = "historyResultType3"
	cmdHistoryResultType4 = "historyResultType4"
	cmdHistoryResultType5 = "historyResultType5"
	cmdHistoryResultType6 = "historyResultType6"
	cmdHistoryResultType7 = "historyResultType7"

	cmdRoadMapType0 = "roadMapType0"
	cmdRoadMapType1 = "roadMapType1"
	cmdRoadMapType2 = "roadMapType2"
	cmdRoadMapType6 = "roadMapType6"
	cmdRoadMapType7 = "roadMapType7"

	cmdRerollDice = "rerollDice"
	cmdRethrow    = "rethrow"

	//pattern 牌型 編號
	patternType25Face  = 11 //5公
	patternType2NiuNiu = 10 //牛牛
	patternType2Niu9   = 9  //牛9
	patternType2Niu8   = 8
	patternType2Niu7   = 7
	patternType2Niu6   = 6
	patternType2Niu5   = 5
	patternType2Niu4   = 4
	patternType2Niu3   = 3
	patternType2Niu2   = 2
	patternType2Niu1   = 1
	patternType2None   = 0
)

//Configurer loads config.json
//載入外部相關的的server
//設定命令、房間類型、pattern等資料
type Configurer interface {
	Version() string
	DBAPIServer() string
	RoadMapAPIHost() string

	ClientReady() int
	ClientNotReady() int

	//未指定的id
	UnAssigned() int
	Active() int
	Inactive() int

	//hall / room
	HallID0() int
	HallID1() int

	RoomStatusBeginBet() int
	RoomStatusEndBet() int
	RoomStatusRoundResult() int
	RoomStatusChangeDealer() int
	RoomStatusChangeBoot() int

	//command
	CmdLogin() string
	CmdLogout() string

	CmdHallInfo() string
	CmdEnterHall() string
	CmdLeaveHall() string
	CmdRoomInfo() string
	CmdEnterRoom() string
	CmdLeaveRoom() string
	CmdRoomStatus() string
	CmdCancelRound() string

	CmdRoundProcess0() string
	CmdRoundProcess1() string
	CmdRoundProcess2() string
	CmdRoundProcess3() string
	CmdRoundProcess4() string
	CmdRoundProcess5() string
	CmdRoundProcess6() string

	CmdRoundResult0() string
	CmdRoundResult1() string
	CmdRoundResult2() string
	CmdRoundResult3() string
	CmdRoundResult4() string
	CmdRoundResult5() string
	CmdRoundResult6() string
	CmdRoundResult7() string

	CmdHistoryResultType0() string
	CmdHistoryResultType1() string
	CmdHistoryResultType2() string
	CmdHistoryResultType3() string
	CmdHistoryResultType4() string
	CmdHistoryResultType5() string
	CmdHistoryResultType6() string
	CmdHistoryResultType7() string

	CmdBetType0() string
	CmdBetType1() string
	CmdBetType2() string
	CmdBetType3() string
	CmdBetType4() string
	CmdBetType5() string
	CmdBetType6() string
	CmdBetType7() string

	CmdRoadMapType0() string
	CmdRoadMapType1() string
	CmdRoadMapType2() string
	CmdRoadMapType6() string
	CmdRoadMapType7() string

	CmdRerollDice() string
	CmdRethrow() string

	CmdUserResult() string
	//user行為
	CmdHeartbeat() string

	//支援性質
	CmdUserCount() string
	CmdClientReady() string
	CmdUserInfo() string
	CmdChangeDealer() string
	CmdHallActive() string
	CmdRoomActive() string
	CmdBroadcast() string
	CmdBanner() string
	CmdServerTime() string

	RoomType0() int
	RoomType1() int
	RoomType2() int
	RoomType3() int
	RoomType4() int
	RoomType5() int
	RoomType6() int
	RoomType7() int

	//pattern
	PatternType25Face() int
	PatternType2NiuNiu() int
	PatternType2Niu9() int
	PatternType2Niu8() int
	PatternType2Niu7() int
	PatternType2Niu6() int
	PatternType2Niu5() int
	PatternType2Niu4() int
	PatternType2Niu3() int
	PatternType2Niu2() int
	PatternType2Niu1() int
	PatternType2None() int

	UserVarDBUserID() string
	UserVarPartnerID() string
	UserVarClientReady() string
}

//Config config main struct
type configurer struct {
	data *configData

	//value
	clientReady    int
	clientNotReady int

	//hall / room
	hallID0 int
	hallID1 int

	// variable

	roomVarTypeData      string
	roomVarHLSURL        string
	roomVarDealer        string
	roomVarBetCountDown  string
	roomVarUserBet       string
	roomVarResultHistory string

	userVarDBUserID    string
	userVarPartnerID   string
	userVarClientReady string
}

type configData struct {
	Version          string
	DBAPIServer      string
	RoadMapAPIServer string
}

//NewConfig make a new config struct
func NewConfigurer(configFileName string) (Configurer, error) {

	defer func() {
		if r := recover(); r != nil {
			//p.GetLogger().LogFile(nxLog.LevelPanic, fmt.Sprintf("CountLinkNum panic:%v\n", r))
			log.Printf("NewConfigurer panic=%v", r)
		}
	}()

	cf := &configurer{

		data: &configData{},

		//value
		clientReady:    1,
		clientNotReady: 0,

		hallID0: 100,
		hallID1: 200,


		//variable
		roomVarTypeData:      "rvtd",
		roomVarHLSURL:        "rvuls",
		roomVarDealer:        "rvde",
		roomVarBetCountDown:  "rvbcd",
		roomVarUserBet:       "rvub",
		roomVarResultHistory: "rvrh",

		userVarDBUserID:    "uid",  //userID  from db
		userVarPartnerID:   "pid",  //partnerID
		userVarClientReady: "uvcr", //收到clientReady 1= ready, 0= not ready
	}

	//config
	path := cf.getConfigFilePosition(configFileName)

	_, err := cf.loadConfig(path, cf.data)

	if err != nil {
		return nil, err
	}

	return cf, err
}

func (c *configurer) getConfigFilePosition(fileName string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}

	var buf bytes.Buffer
	buf.WriteString(dir)
	buf.WriteString("/")
	buf.WriteString(fileName)

	return buf.String()
}

func (c *configurer) loadConfig(filePath string, container interface{}) (interface{}, error) {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	//con := &configData{}
	//unmarshal to struct
	if err := json.Unmarshal(b, container); err != nil {
		return nil, err
	}

	return container, nil
}
func (c *configurer) loadFileToStruct(fileName string, container [][][]int) ([][][]int, error) {
	path := c.getConfigFilePosition(fileName)
	_, err := c.loadConfig(path, &container)

	if err != nil {
		return nil, err
	}
	return container, nil
}

func (c *configurer) Version() string        { return c.data.Version }
func (c *configurer) DBAPIServer() string    { return c.data.DBAPIServer }
func (c *configurer) RoadMapAPIHost() string { return c.data.RoadMapAPIServer }

func (c *configurer) ClientReady() int    { return c.clientReady }
func (c *configurer) ClientNotReady() int { return c.clientNotReady }

//未指定的id
func (c *configurer) UnAssigned() int { return unAssigned }
func (c *configurer) Active() int     { return active }
func (c *configurer) Inactive() int   { return inactive }

//hall / room
func (c *configurer) HallID0() int { return c.hallID0 }
func (c *configurer) HallID1() int { return c.hallID1 }

//room status
func (c *configurer) RoomStatusBeginBet() int     { return roomStatusBeginBet }
func (c *configurer) RoomStatusEndBet() int       { return roomStatusEndBet }
func (c *configurer) RoomStatusRoundResult() int  { return roomStatusRoundResult }
func (c *configurer) RoomStatusChangeBoot() int   { return roomStatusChangeBoot }
func (c *configurer) RoomStatusChangeDealer() int { return roomStatusChangeDealer }

//command
func (c *configurer) CmdLogin() string  { return cmdLogin }
func (c *configurer) CmdLogout() string { return cmdLogout }

func (c *configurer) CmdHallInfo() string    { return cmdHallInfo }
func (c *configurer) CmdEnterHall() string   { return cmdEnterHall }
func (c *configurer) CmdLeaveHall() string   { return cmdLeaveHall }
func (c *configurer) CmdRoomInfo() string    { return cmdRoomInfo }
func (c *configurer) CmdEnterRoom() string   { return cmdEnterRoom }
func (c *configurer) CmdLeaveRoom() string   { return cmdLeaveRoom }
func (c *configurer) CmdRoomStatus() string  { return cmdRoomStatus }
func (c *configurer) CmdCancelRound() string { return cmdCancelRound }

func (c *configurer) CmdRoundProcess0() string { return cmdRoundProcess0 }
func (c *configurer) CmdRoundProcess1() string { return cmdRoundProcess1 }
func (c *configurer) CmdRoundProcess2() string { return cmdRoundProcess2 }
func (c *configurer) CmdRoundProcess3() string { return cmdRoundProcess3 }
func (c *configurer) CmdRoundProcess4() string { return cmdRoundProcess4 }
func (c *configurer) CmdRoundProcess5() string { return cmdRoundProcess5 }
func (c *configurer) CmdRoundProcess6() string { return cmdRoundProcess6 }

func (c *configurer) CmdRoundResult0() string { return cmdRoundResultType0 }
func (c *configurer) CmdRoundResult1() string { return cmdRoundResultType1 }
func (c *configurer) CmdRoundResult2() string { return cmdRoundResultType2 }
func (c *configurer) CmdRoundResult3() string { return cmdRoundResultType3 }
func (c *configurer) CmdRoundResult4() string { return cmdRoundResultType4 }
func (c *configurer) CmdRoundResult5() string { return cmdRoundResultType5 }
func (c *configurer) CmdRoundResult6() string { return cmdRoundResultType6 }
func (c *configurer) CmdRoundResult7() string { return cmdRoundResultType7 }

func (c *configurer) CmdHistoryResultType0() string { return cmdHistoryResultType0 }
func (c *configurer) CmdHistoryResultType1() string { return cmdHistoryResultType1 }
func (c *configurer) CmdHistoryResultType2() string { return cmdHistoryResultType2 }
func (c *configurer) CmdHistoryResultType3() string { return cmdHistoryResultType3 }
func (c *configurer) CmdHistoryResultType4() string { return cmdHistoryResultType4 }
func (c *configurer) CmdHistoryResultType5() string { return cmdHistoryResultType5 }
func (c *configurer) CmdHistoryResultType6() string { return cmdHistoryResultType6 }
func (c *configurer) CmdHistoryResultType7() string { return cmdHistoryResultType7 }

func (c *configurer) CmdBetType0() string { return cmdBetType0 }
func (c *configurer) CmdBetType1() string { return cmdBetType1 }
func (c *configurer) CmdBetType2() string { return cmdBetType2 }
func (c *configurer) CmdBetType3() string { return cmdBetType3 }
func (c *configurer) CmdBetType4() string { return cmdBetType4 }
func (c *configurer) CmdBetType5() string { return cmdBetType5 }
func (c *configurer) CmdBetType6() string { return cmdBetType6 }
func (c *configurer) CmdBetType7() string { return cmdBetType7 }

func (c *configurer) CmdRoadMapType0() string { return cmdRoadMapType0 }
func (c *configurer) CmdRoadMapType1() string { return cmdRoadMapType1 }
func (c *configurer) CmdRoadMapType2() string { return cmdRoadMapType2 }
func (c *configurer) CmdRoadMapType6() string { return cmdRoadMapType6 }
func (c *configurer) CmdRoadMapType7() string { return cmdRoadMapType7 }

func (c *configurer) CmdRerollDice() string { return cmdRerollDice }
func (c *configurer) CmdRethrow() string    { return cmdRethrow }

func (c *configurer) CmdUserResult() string { return cmdUserResult }

//user行為
func (c *configurer) CmdHeartbeat() string { return cmdHeartbeat }

//支援性質
func (c *configurer) CmdUserCount() string    { return cmdUserCount }
func (c *configurer) CmdClientReady() string  { return cmdClientReady }
func (c *configurer) CmdUserInfo() string     { return cmdUserInfo }
func (c *configurer) CmdChangeDealer() string { return cmdChangeDealer }
func (c *configurer) CmdRoomActive() string   { return cmdRoomActive }
func (c *configurer) CmdHallActive() string   { return cmdHallActive }
func (c *configurer) CmdBroadcast() string    { return cmdBroadcast }
func (c *configurer) CmdBanner() string       { return cmdBanner }
func (c *configurer) CmdServerTime() string   { return cmdServerTime }

func (c *configurer) RoomType0() int { return roomTypeBaccarat }
func (c *configurer) RoomType1() int { return roomTypeDragonTiger }
func (c *configurer) RoomType2() int { return roomTypeNiuniu }
func (c *configurer) RoomType3() int { return roomTypeFantan }
func (c *configurer) RoomType4() int { return roomTypeSangong }
func (c *configurer) RoomType5() int { return roomTypeThisBar }
func (c *configurer) RoomType6() int { return roomTypeSicbo }
func (c *configurer) RoomType7() int { return roomTypeRolette }

//pattern type2
func (c *configurer) PatternType25Face() int  { return patternType25Face }
func (c *configurer) PatternType2NiuNiu() int { return patternType2NiuNiu }
func (c *configurer) PatternType2Niu9() int   { return patternType2Niu9 }
func (c *configurer) PatternType2Niu8() int   { return patternType2Niu8 }
func (c *configurer) PatternType2Niu7() int   { return patternType2Niu7 }
func (c *configurer) PatternType2Niu6() int   { return patternType2Niu6 }
func (c *configurer) PatternType2Niu5() int   { return patternType2Niu5 }
func (c *configurer) PatternType2Niu4() int   { return patternType2Niu4 }
func (c *configurer) PatternType2Niu3() int   { return patternType2Niu3 }
func (c *configurer) PatternType2Niu2() int   { return patternType2Niu2 }
func (c *configurer) PatternType2Niu1() int   { return patternType2Niu1 }
func (c *configurer) PatternType2None() int   { return patternType2None }

//user variable
func (c *configurer) UserVarDBUserID() string {
	return c.userVarDBUserID
}
func (c *configurer) UserVarPartnerID() string {
	return c.userVarPartnerID
}
func (c *configurer) UserVarClientReady() string {
	return c.userVarClientReady
}

func (c *configurer) RoomVarTypeData() string {
	return c.roomVarTypeData
}

func (c *configurer) RoomVarHLSURL() string {
	return c.roomVarHLSURL
}

func (c *configurer) RoomVarDealer() string {
	return c.roomVarDealer
}
func (c *configurer) RoomVarBetCountDown() string {
	return c.roomVarBetCountDown
}
func (c *configurer) RoomVarUserBet() string {
	return c.roomVarUserBet
}
func (c *configurer) RoomVarResultHistory() string {
	return c.roomVarResultHistory
}

package config

import (
	"github.com/cruisechang/liveServer/config/roomConf"
)

const (
	//code for client
	CodeSuccess = 0 //Both for client and backstage
	CodePanic   = 1

	CodeUserNil = 50

	CodeLoginFailed           = 100
	CodeLoginRequestError     = 101
	CodeLoginPasswordError    = 102
	CodeLoginUserAleradyIn    = 103
	CodeLoginParamsError      = 104
	CodeLoginAccessTokenError = 105
	CodeLoginUseDeactive      = 106

	CodeLogutFailed = 200

	CodeGameInfoFailed = 300

	CodeUserInfoFailed = 400

	CodeBetFailed = 500

	CodeHeartbeatFailed = 600

	CodeJsonFailed              = 1000
	CodeJsonUnmarshalJsonFailed = 1001 //解析json錯誤
	CodeMarshalJsonFailed       = 1002 //編碼成json錯誤

	CodeBase64DecodeFailed = 1100
	CodeBase64EncodeFailed = 1101

	CodeReceivedDataError = 1150

	//hall
	CodeHallInactive = 1200
	CodeHallIDError  = 1201
	CodeHallNotFound = 1202

	//in processor
	CodeProcessorUserNil = 1250

	//room
	CodeRoomInactive = 1300
	CodeRoomIDError  = 1301
	CodeRoomNotFound = 1302

	//user
	CodeCreditNotEnough = 1400

	//DB
	CodeDBError               = 1500
	CodeDBHTTPStatusCodeNotOK = 1501
	CodeDBBodyNil             = 1502
	CodeDBQueryFail           = 1503
)

/**
* client
 */

//for heartbeat,clientReady, or emptyRes
//enterHall
type EmptyResData struct {
}

//hall and room
type HallInfoCmdData struct {
	HallID int `json:"hallID"`
}
type HallInfoResData struct {
	HallID int    `json:"hallID"`
	Name   string `json:"name"`
	Active int    `json:"active"`
}
type EnterHallCmdData struct {
	HallID int `json:"hallID"`
}
type LeaveHallCmdData struct {
}
type RoomInfoCmdData struct {
	HallID int `json:"hallID"`
	RoomID int `json:"roomID"`
}

type RoomInfoResData struct {
	RoomID       int             `json:"roomID"`
	Name         string          `json:"name"`
	Active       int             `json:"active"`
	HLSURL       string          `json:"hlsURL"`
	Type         int             `json:"type"`
	Dealer       roomConf.Dealer `json:"dealer"`
	BetCoundDown int             `json:"betCoundDown"`

	TypeData0 *roomConf.TypeData0 `json:"type0"`
	TypeData1 *roomConf.TypeData1 `json:"type1"`
	TypeData2 *roomConf.TypeData2 `json:"type2"`
	TypeData3 *roomConf.TypeData3 `json:"type3"`
	TypeData4 *roomConf.TypeData4 `json:"type4"`
	TypeData5 *roomConf.TypeData5 `json:"type5"`
	TypeData6 *roomConf.TypeData6 `json:"type6"`
	TypeData7 *roomConf.TypeData7 `json:"type7"`
}

type RoomStatusResData struct {
	RoomID      int   `json:"roomID"`
	Status      int   `json:"status"`
	StatusStart int64 `json:"statusStart"` //876778778     // long int  狀態開始時間  // 从January 1, 1970, 00:00:00开的毫秒数
	Boot        int   `json:"boot"`
	Round       int64 `json:"round"`
}

type EnterRoomCmdData struct {
	HallID int `json:"hallID"`
	RoomID int `json:"roomID"`
}

type LeaveRoomCmdData struct {
}
type CancelRoundResData struct {
	HallID int   `json:"hallID"`
	RoomID int   `json:"roomID"`
	Boot   int   `json:"boot"`
	Round  int64 `json:"round"`
}

//
//user behavior
type LoginCmdData struct {
	SessionID string `json:"sessionID"`
}
type LoginCmdResData struct {
	UserID int `json:"userID"`
}
type LogoutCmdData struct {
	Type int `json:"type"`
}
type UserInfoCmdData struct {
	UserID int `json:"userID"`
}
type UserInfoResData struct {
	UserID  int     `json:"userID"`
	Account string  `json:"account"`
	Credit  float32 `json:"credit"`
	Name    string  `json:"name"`
}
type HeartbeatCmdData struct {
}

type ClientReadyCmdData struct {
}

type UserResultResData struct {
	RoomID int                       `json:"roomID"`
	Round  int64                     `json:"round"`
	Result []UserResultResDataResult `json:"result"`
}

type UserResultResDataResult struct {
	UserID        int64   `json:"userID"`
	Result        int     `json:"result"`
	ResultCredit  float32 `json:"resultCredit"`
	BalanceCredit float32 `json:"balanceCredit"`
}

//支援
type RoomActiveResData struct {
	RoomID int `json:"roomID"`
	Active int `json:"active"`
}
type ChangeDealerResData struct {
	HallID      int    `json:"hallID"`
	RoomID      int    `json:"roomID"`
	DealerID    int    `json:"dealerID"`
	Name        string `json:"name"`
	PortraitURL string `json:"portraitURL"`
}
type ServerTimeResData struct {
	ServerTime int `json:"serverTime"`
}

//for all type
type HistoryResultTypeCmdData struct {
	HallID int `json:"hallID"`
	RoomID int `json:"roomID"`
}

//baccarat type0
type BetType0CmdData struct {
	RoomID      int `json:"roomID"`
	Banker      int `json:"banker"`
	Player      int `json:"player"`
	Tie         int `json:"tie"`
	BankerPair  int `json:"bankerPair"`
	PlayerPair  int `json:"playerPair"`
	Big         int `json:"big"`
	Small       int `json:"small"`
	AnyPair     int `json:"anyPair"`
	PerfectPair int `json:"perfectPair"`
	SuperSix    int `json:"superSix"`
	Commission  int `json:"commission"`
}
type BetType0ResData struct {
	UserID      int `json:"userID"`
	RoomID      int `json:"roomID"`
	Banker      int `json:"banker"`
	Player      int `json:"player"`
	Tie         int `json:"tie"`
	BankerPair  int `json:"bankerPair"`
	PlayerPair  int `json:"playerPair"`
	Big         int `json:"big"`
	Small       int `json:"small"`
	AnyPair     int `json:"anyPair"`
	PerfectPair int `json:"perfectPair"`
	SuperSix    int `json:"superSix"`
}

type RoundProcess0ResData struct {
	RoomID int `json:"roomID"`
	Owner  int `json:"owner"`
	Poker  int `json:"poker"`
	Index  int `json:"index"`
}
type RoundResultType0ResData struct {
	RoomID      int   `json:"roomID"`
	Round       int64 `json:"round"`
	Result      int   `json:"result"`
	BankerPair  int   `json:"bankerPair"`
	PlayerPair  int   `json:"playerPair"`
	BigSmall    int   `json:"bigSmall"`
	AnyPair     int   `json:"anyPair"`
	PerfectPair int   `json:"perfectPair"`
	SuperSix    int   `json:"superSix"`
	BankerPoint int   `json:"bankerPoint"`
	PlayerPoint int   `json:"playerPoint"`
}
type HistoryResultType0ResData struct {
	HallID int                         `json:"hallID"`
	RoomID int                         `json:"roomID"`
	Result roomConf.HistoryResultType0 `json:"result"`
}

//dragonTiger type1
type BetType1CmdData struct {
	RoomID      int `json:"roomID"`
	Dragon      int `json:"dragon"`
	Tiger       int `json:"tiger"`
	Tie         int `json:"tie"`
	DragonOdd   int `json:"dragonOdd"`
	DragonEven  int `json:"dragonEven"`
	DragonRed   int `json:"dragonRed"`
	DragonBlack int `json:"dragonBlack"`
	TigerOdd    int `json:"tigerOdd"`
	TigerEven   int `json:"tigerEven"`
	TigerRed    int `json:"tigerRed"`
	TigerBlack  int `json:"tigerBlack"`
}
type BetType1ResData struct {
	UserID      int `json:"userID"`
	RoomID      int `json:"roomID"`
	Dragon      int `json:"dragon"`
	Tiger       int `json:"tiger"`
	Tie         int `json:"tie"`
	DragonOdd   int `json:"dragonOdd"`
	DragonEven  int `json:"dragonEven"`
	DragonRed   int `json:"dragonRed"`
	DragonBlack int `json:"dragonBlack"`
	TigerOdd    int `json:"tigerOdd"`
	TigerEven   int `json:"tigerEven"`
	TigerRed    int `json:"tigerRed"`
	TigerBlack  int `json:"tigerBlack"`
}
type RoundProcess1ResData struct {
	RoomID int `json:"roomID"`
	Owner  int `json:"owner"`
	Poker  int `json:"poker"`
	Index  int `json:"index"`
}

type RoundResultType1ResData struct {
	RoomID         int   `json:"roomID"`
	Round          int64 `json:"round"`
	Result         int   `json:"result"`
	DragonOddEven  int   `json:"dragonOddEven"`
	DragonRedBlack int   `json:"dragonRedBlack"`
	TigerOddEven   int   `json:"tigerOddEven"`
	TigerRedBlack  int   `json:"tigerRedBlack"`
}

type HistoryResultType1ResData struct {
	HallID int                         `json:"hallID"`
	RoomID int                         `json:"roomID"`
	Result roomConf.HistoryResultType1 `json:"result"`
}

//niuniu type2
type BetType2CmdData struct {
	RoomID int                  `json:"roomID"`
	Owner1 BetType2CmdDataOwner `json:"owner1"`
	Owner2 BetType2CmdDataOwner `json:"owner2"`
	Owner3 BetType2CmdDataOwner `json:"owner3"`
}
type BetType2CmdDataOwner struct {
	WinEqual  int `json:"winEqual"`
	WinTimes  int `json:"winTimes"`
	LoseTimes int `json:"loseTimes"`
	LoseEqual int `json:"loseEqual"`
}
type BetType2ResData struct {
	UserID int                  `json:"userID"`
	RoomID int                  `json:"roomID"`
	Owner1 BetType2CmdDataOwner `json:"owner1"`
	Owner2 BetType2CmdDataOwner `json:"owner2"`
	Owner3 BetType2CmdDataOwner `json:"owner3"`
}
type RoundProcess2ResData struct {
	RoomID int `json:"roomID"`
	Owner  int `json:"owner"`
	Poker  int `json:"poker"`
	Index  int `json:"index"`
}

type RoundResultType2ResData struct {
	RoomID int                   `json:"roomID"`
	Round  int64                 `json:"round"`
	Head   int                   `json:"head"`
	Owner0 RoundResultType2Owner `json:"owner0"`
	Owner1 RoundResultType2Owner `json:"owner1"`
	Owner2 RoundResultType2Owner `json:"owner2"`
	Owner3 RoundResultType2Owner `json:"owner3"`
}
type RoundResultType2Owner struct {
	Result  int   `json:"result"`
	Pattern int   `json:"pattern"`
	Poker   []int `json:"poker"`
}

type HistoryResultType2CmdData struct {
	HallID int `json:"hallID"`
	RoomID int `json:"roomID"`
}

type HistoryResultType2ResData struct {
	HallID int                         `json:"hallID"`
	RoomID int                         `json:"roomID"`
	Result roomConf.HistoryResultType2 `json:"result"`
}

//sicbo type 6
type BetType6CmdData struct {
	RoomID int   `json:"roomID"`
	Small  int   `json:"small"`
	Big    int   `json:"big"`
	Odd    int   `json:"odd"`
	Even   int   `json:"even"`
	Sum    []int `json:"sum"`
	Dice   []int `json:"dice"`
	Triple []int `json:"triple"`
	Pair   []int `json:"pair"`
	Paigow []int `json:"paigow"`
}
type BetType6ResData struct {
	UserID int   `json:"userID"`
	RoomID int   `json:"roomID"`
	Small  int   `json:"small"`
	Big    int   `json:"big"`
	Odd    int   `json:"odd"`
	Even   int   `json:"even"`
	Sum    []int `json:"sum"`
	Dice   []int `json:"dice"`
	Triple []int `json:"triple"`
	Pair   []int `json:"pair"`
	Paigow []int `json:"paigow"`
}

type RoundResultType6ResData struct {
	RoomID   int     `json:"roomID"`
	BigSmall int     `json:"bigSmall"`
	OddEven  int     `json:"oddEven"`
	Sum      int     `json:"sum"`
	Dice     []int   `json:"dice"`
	Triple   int     `json:"triple"`
	Pair     int     `json:"pair"`
	Paigow   [][]int `json:"paigow"`
}
type HistoryResultType6CmdData struct {
	HallID int `json:"hallID"`
	RoomID int `json:"roomID"`
}

type HistoryResultType6ResData struct {
	HallID   int   `json:"hallID"`
	RoomID   int   `json:"roomID"`
	Dice     []int `json:"dice"`
	Sum      int   `json:"sum"`
	BigSmall int   `json:"bigSmall"`
	OddEven  int   `json:"oddEven"`
}

type RerollDiceResData struct {
	HallID int   `json:"hallID"`
	RoomID int   `json:"roomID"`
	Round  int64 `json:"round"`
}

//roulette type7
type BetType7CmdData struct {
	RoomID int   `json:"roomID"`
	One    []int `json:"one"`    //下注 0-36單一數字 37 index   [0,0,0...],
	Two    []int `json:"two"`    //下注兩個號碼   60 index [0,0,0,0….],
	Three  []int `json:"three"`  //下注三個號碼   14 index
	Four   []int `json:"four"`   //四個號碼  23 index
	Six    []int `json:"six"`    //六個號碼  11 index
	Column []int `json:"column"` // 一列,    3 index
	Dozen  []int `json:"dozen"`  //12個一組   3 index
	Big    int   `json:"big"`    //  大 1:1
	Small  int   `json:"small"`  //小
	Odd    int   `json:"odd"`    //  單
	Even   int   `json:"even"`   //雙
	Red    int   `json:"red"`    //紅
	Black  int   `json:"black"`  // 黑
}
type BetType7ResData struct {
	UserID int   `json:"userID"`
	RoomID int   `json:"roomID"`
	One    []int `json:"one"`    //下注 0-36單一數字 37 index   [0,0,0...],
	Two    []int `json:"two"`    //下注兩個號碼   60 index [0,0,0,0….],
	Three  []int `json:"three"`  //下注三個號碼   14 index
	Four   []int `json:"four"`   //四個號碼  23 index
	Six    []int `json:"six"`    //六個號碼  11 index
	Column []int `json:"column"` // 一列,    3 index
	Dozen  []int `json:"dozen"`  //12個一組   3 index
	Big    int   `json:"big"`    //  大 1:1
	Small  int   `json:"small"`  //小
	Odd    int   `json:"odd"`    //  單
	Even   int   `json:"even"`   //雙
	Red    int   `json:"red"`    //紅
	Black  int   `json:"black"`  // 黑
}

type HistoryResultType7ResData struct {
	HallID int                         `json:"hallID"`
	RoomID int                         `json:"roomID"`
	Result roomConf.HistoryResultType7 `json:"result"`
}
type RethrowResData struct {
	HallID int   `json:"hallID"`
	RoomID int   `json:"roomID"`
	Round  int64 `json:"round"`
}
type RoundResultType7ResData struct {
	RoomID   int `json:"roomID"`
	Result   int `json:"result"`
	BigSmall int `json:"bigSmall"`
	OddEven  int `json:"oddEven"`
	RedBlack int `json:"redBlack"`
	Dozen    int `json:"dozen"`
	Column   int `json:"column"`
}

//roadmap 回傳給client
type RoadMapType0ResData struct {
	HallID        int    `json:"hallID"`
	RoomID        int    `json:"roomID"`
	PicBall       string `json:"picBall"`
	PicBig        string `json:"picBig"`
	Pic3          string `json:"pic3"`
	PicBallBanker string `json:"picBallBanker"` //珠盤莊問路
	PicBallPlayer string `json:"picBallPlayer"` //珠盤閒問路
	PicBigBanker  string `json:"picBigBanker"`  //大路莊問路
	PicBigPlayer  string `json:"picBigPlayer"`  //大路閒問路
	Pic3Banker    string `json:"pic3Banker"`    //3小路莊問路
	Pic3Player    string `json:"pic3Player"`    //3小路閒問路
	Sum           int    `json:"sum"`           //總局數
	Tie           int    `json:"tie"`           //和局數
	Banker        int    `json:"banker"`        //莊贏局數
	Player        int    `json:"player"`        //閒贏局數
}
type RoadMapType1ResData struct {
	HallID        int    `json:"hallID"`
	RoomID        int    `json:"roomID"`
	PicBall       string `json:"picBall"`       //珠盤路圖
	PicBig        string `json:"picBig"`        //大路圖
	Pic3          string `json:"pic3"`          //三小路圖
	PicBallDragon string `json:"picBallDragon"` //珠盤龍問路
	PicBallTiger  string `json:"picBallTiger"`  //xxxx””,         //珠盤虎問路
	PicBigDragon  string `json:"picBigDragon"`  //xxxx””,      //大路龍問路
	PicBigTiger   string `json:"picBigTiger"`   //xxxx””,          //大路虎問路
	Pic3Dragon    string `json:"pic3Dragon"`    //xxxx””,         //3小路龍問路
	Pic3Tiger     string `json:"pic3Tiger"`     //xxxx””,             //3小路虎問路
	Sum           int    `json:"sum"`           //總局數
	Tie           int    `json:"tie"`           //和局數
	Dragon        int    `json:"dragon"`        //龍贏局數
	Tiger         int    `json:"tiger"`         //虎贏局數
}

type RoadMapType2ResData struct {
	HallID int    `json:"hallID"`
	RoomID int    `json:"roomID"`
	Pic    string `json:"pic"`    //路單圖
	Sum    int    `json:"sum"`    //總局數
	Owner1 int    `json:"owner1"` //閒1贏局數
	Owner2 int    `json:"owner2"` //閒2贏局數
	Owner3 int    `json:"owner3"` //閒3贏局數
}

type RoadMapType6ResData struct {
	HallID      int    `json:"hallID"`
	RoomID      int    `json:"roomID"`
	PicBigSmall string `json:"picBigSmall"` //大小路單圖
	PicOddEven  string `json:"picOddEven"`  //單雙路單圖
	PicSum      string `json:"picSum"`      //和值(三個骰子總和)路單圖
	PicDice     string `json:"picDice"`     //三個骰子的點數
	Big         int    `json:"big"`         //點數大局數
	Small       int    `json:"small"`       //點數小局數
	Odd         int    `json:"odd"`         //單數局數
	Even        int    `json:"even"`        //雙數局數
	Triple      int    `json:"triple"`      //圍骰局數
}

type RoadMapType7ResData struct {
	HallID       int     `json:"hallID"`
	RoomID       int     `json:"roomID"`
	Pic          string  `json:"pic"`          //路單圖
	Hot          []int   `json:"hot"`          //熱門點數點
	Cold         []int   `json:"cold"`         //冷門點數
	RedPercent   float32 `json:"redPercent"`   //出現紅色的百分比，
	BlackPercent float32 `json:"blackPercent"` //出出現黑色的百分比，
}

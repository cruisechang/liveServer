package config

const (
	LimitationID0   = 0
	LimitationID100 = 100
	LimitationID200 = 200
	LimitationID600 = 600
	LimitationID700 = 700
)

//從db limitation 轉換過
type TransferredLimitation struct {
	LimitationID int         `json:"limitationID"`
	Limitation   interface{} `json:"limitation"`
}
type Limitation0 struct {
	BetLimit         []int `json:"betLimit"` //限紅最小, 最小
	BankerLimit      []int `json:"bankerLimit"`
	PlayerLimit      []int `json:"playerLimit"`
	TieLimit         []int `json:"tieLimit"`
	BankerPairLimit  []int `json:"bankerPairLimit"`
	PlayerPairLimit  []int `json:"playerPairLimit"`
	AnyPairLimit     []int `json:"anyPairLimit"`
	PerfectPairLimit []int `json:"perfectPairLimit"`
	SuperSixLimit    []int `json:"superSixLimit"`
	BigSmallLimit    []int `json:"bigSmallLimit"`
}
type Limitation100 struct {
	BetLimit      []int `json:"betLimit"` //限紅最小, 最小
	DragonLimit   []int `json:"dragonLimit"`
	TigerLimit    []int `json:"tigerLimit"`
	TieLimit      []int `json:"tieLimit"`
	OddEvenLimit  []int `json:"oddEvenLimit"`
	RedBlackLimit []int `json:"redBlackLimit"`
}
type Limitation200 struct {
	BetLimit       []int `json:"betLimit"` //限紅最小, 最小
	WinTimesLimit  []int `json:"winTimesLimit"`
	WinEqualLimit  []int `json:"winEqualLimit"`
	LoseTimesLimit []int `json:"loseTimesLimit"`
	LoseEqualLimit []int `json:"loseEqualLimit"`
}
type Limitation600 struct {
	BetLimit         []int `json:"betLimit"`         //限紅最小, 最小
	BigSmallLimit    []int `json:"bigSmallLimit"`    //限紅最小, 最小
	OddEvenLimit     []int `json:"oddEvenLimit"`     //單雙
	Sum0417Limit     []int `json:"sum0417Limit"`     //總和 4 ,17
	Sum0516Limit     []int `json:"sum0516Limit"`     //總和 5 ,16
	Sum0615Limit     []int `json:"sum0615Limit"`     //總和 6 ,15
	Sum0714Limit     []int `json:"sum0714Limit"`     //總和 7 ,14
	Sum0813Limit     []int `json:"sum0813Limit"`     //總和 8 ,13
	Sum09101112Limit []int `json:"sum09101112Limit"` //總和 其他

	DiceLimit      []int `json:"diceLimit"`      //單一點數
	PairLimit      []int `json:"pairLimit"`      //對子
	PaigowLimit    []int `json:"paigowLimit"`    //牌九
	TripleLimit    []int `json:"tripleLimit"`    //圍骰
	TripleAllLimit []int `json:"tripleAllLimit"` //全圍
}
type Limitation700 struct {
	BetLimit      []int `json:"betLimit"`      //限紅最小, 最小
	OneLimit      []int `json:"oneLimit"`      //
	TwoLimit      []int `json:"twoLimit"`      //
	ThreeLimit    []int `json:"threeLimit"`    //
	FourLimit     []int `json:"fourLimit"`     //
	SixLimit      []int `json:"sixLimit"`      //
	ColumnLimit   []int `json:"columnLimit"`   //
	DozenLimit    []int `json:"dozenLimit"`    //
	BigSmallLimit []int `json:"bigSmallLimit"` //

	OddEvenLimit  []int `json:"oddEvenLimit"`  //
	RedBlackLimit []int `json:"redBlackLimit"` //
}

//rpc used
//round post
type RoundPostParam struct {
	HallID   int    `json:"hallID"`
	RoomID   int    `json:"roomID"`
	RoomType int    `json:"roomType"`
	Brief    string `json:"brief"`
	Record   string `json:"record"` //json string
	Status   int    `json:"status"`
}

type RoundRecordType0 struct {
	Boot        int   `json:"boot"`
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

type RoundRecordType1 struct {
	Boot           int   `json:"boot"`
	Round          int64 `json:"round"`
	Result         int   `json:"result"`
	DragonOddEven  int   `json:"dragonOddEven"`
	DragonRedBlack int   `json:"dragonRedBlack"`
	TigerOddEven   int   `json:"tigerOddEven"`
	TigerRedBlack  int   `json:"tigerRedBlack"`
}

type RoundRecordType2 struct {
	Boot   int                   `json:"boot"`
	Round  int64                 `json:"round"`
	Head   int                   `json:"head"`
	Owner0 RoundResultType2Owner `json:"owner0"`
	Owner1 RoundResultType2Owner `json:"owner1"`
	Owner2 RoundResultType2Owner `json:"owner2"`
	Owner3 RoundResultType2Owner `json:"owner3"`
}

//type RoundRecordType2Owner struct {
//	Result  int   `json:"result"`
//	Pattern int   `json:"pattern"`
//	Poker   []int `json:"poker"`
//}

type RoundRecordType6 struct {
	Round    int64   `json:"round"`
	Dice     []int   `json:"dice"`
	Sum      int     `json:"sum"`
	BigSmall int     `json:"bigSmall"`
	OddEven  int     `json:"oddEven"`
	Triple   int     `json:"triple"`
	Pair     int     `json:"pair"`
	Paigow   [][]int `json:"paigow"`
}

type RoundRecordType7 struct {
	Round    int64 `json:"round"`
	Result   int   `json:"result"`
	BigSmall int   `json:"bigSmall"`
	OddEven  int   `json:"oddEven"`
	RedBlack int   `json:"redBlack"`
	Dozen    int   `json:"dozen"`
	Column   int   `json:"column"`
}

//bet
type BetPostParam struct {
	PartnerID      int64   `json:"partnerID"`
	UserID         int64   `json:"userID"`
	RoomID         int     `json:"roomID"`
	RoomType       int     `json:"roomType"`
	Round          int64   `json:"round"`
	SeatID         int     `json:"seatID"`
	BetCredit      int     `json:"betCredit"`
	ActiveCredit   int     `json:"activeCredit"`
	PrizeCredit    float32 `json:"prizeCredit"`
	ResultCredit   float32 `json:"resultCredit"`
	BalanceCredit  float32 `json:"balanceCredit"`
	OriginalCredit float32 `json:"originalCredit"`
	Record         string  `json:"record"`
	Status         int     `json:"status"`
}

//banner for configure
type Banner struct {
	PicURL   string `json:"picURL"`
	LinkURL  string `json:"linkURL"`
	Platform uint   `json:"platform"`
}

//broadcast for configure
type Broadcast struct {
	Content     string `json:"content"`
	Internal    int    `json:"internal"`
	RepeatTimes int    `json:"repeatTimes"`
}

//路單
//RoadMapRequest posts data to road map server
type RoadMapRequest struct {
	HallID   int    `json:"hallID"`
	RoomID   int    `json:"roomID"`
	RoomType int    `json:"roomType"`
	Result   string `json:"result"`
}
//RoadMapResponse receive from road map server
type RoadMapResponse struct {
	Code     int    `json:"code"`
	RoomType int    `json:"roomType"`
	Result   string `json:"result"`
}
type HistoryResultType0 [][]int32

//type 1
type HistoryResultType1 []int32

//type 2
type HistoryResultType2 [][]int32

//type6
type HistoryResultType6 struct {
	Dice     []int `json:"dice"`
	Sum      int   `json:"sum"`
	BigSmall int   `json:"bigSmall"`
	OddEven  int   `json:"oddEven"`
}

//type7
type HistoryResultType7 []int32

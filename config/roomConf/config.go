package roomConf

//room
//type RoomField struct {
//	RoomID int    `json:"roomID"`
//	Name   string `json:"name"`
//	Type   int    `json:"type"`
//	Active int    `json:"active"`
//	HLSURL string `json:"hlsURL"`
//	Dealer       RoomFieldDealer
//	BetCoundDown int `json:"betCoundDown"`
//
//	TypeData      interface{} //type data for different type room,  bet limit , boot, round....
//	ResultHistory interface{} //  result history
//}
//
//type RoomFieldDealer struct {
//	DealerID    int    `json:"dealerID"`
//	Name        string `json:"name"`
//	PortraitURL string `json:"protraitURL"`
//}

type Dealer struct {
	DealerID    int    `json:"dealerID"`
	Name        string `json:"name"`
	PortraitURL string `json:"portraitURL"`
}
type TypeData0 struct {
	Boot             int   `json:"boot"`     //靴id 只有牌類有用
	Round            int64 `json:"round"`    //局id
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

type TypeData1 struct {
	Boot          int   `json:"boot"`     //靴id 只有牌類有用
	Round         int64 `json:"round"`    //局id
	BetLimit      []int `json:"betLimit"` //限紅最小, 最小
	DragonLimit   []int `json:"dragonLimit"`
	TigerLimit    []int `json:"tigerLimit"`
	TieLimit      []int `json:"tieLimit"`
	OddEvenLimit  []int `json:"oddEvenLimit"`
	RedBlackLimit []int `json:"redBlackLimit"`
}
type TypeData2 struct {
	Boot           int   `json:"boot"`     //靴id 只有牌類有用
	Round          int64 `json:"round"`    //局id
	BetLimit       []int `json:"betLimit"` //限紅最小, 最小
	WinTimesLimit  []int `json:"winTimesLimit"`
	WinEqualLimit  []int `json:"winEqualLimit"`
	LoseTimesLimit []int `json:"loseTimesLimit"`
	LoseEqualLimit []int `json:"loseEqualLimit"`
}
type TypeData3 struct {
	Boot     int   `json:"boot"`     //靴id 只有牌類有用
	Round    int64 `json:"round"`    //局id
	BetLimit []int `json:"betLimit"` //限紅最小, 最小
}
type TypeData4 struct {
	Boot     int   `json:"boot"`     //靴id 只有牌類有用
	Round    int64 `json:"round"`    //局id
	BetLimit []int `json:"betLimit"` //限紅最小, 最小
}
type TypeData5 struct {
	Boot     int   `json:"boot"`     //靴id 只有牌類有用
	Round    int64 `json:"round"`    //局id
	BetLimit []int `json:"betLimit"` //限紅最小, 最小
}
type TypeData6 struct {
	Boot     int   `json:"boot"`     //靴id 只有牌類有用
	Round    int64 `json:"round"`    //局id
	BetLimit []int `json:"betLimit"` //限紅最小, 最小

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
type TypeData7 struct {
	Boot     int   `json:"boot"`     //靴id 只有牌類有用
	Round    int64 `json:"round"`    //局id
	BetLimit []int `json:"betLimit"` //限紅最小, 最小

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

//type 0
type BetType0Data struct {
	Banker      int
	Player      int
	Tie         int
	BankerPair  int
	PlayerPair  int
	Big         int
	Small       int
	AnyPair     int
	PerfectPair int
	SuperSix    int
	Commission  int
}
type HistoryResultType0 [][]int32

//type 1
type BetType1Data struct {
	Dragon      int
	Tiger       int
	Tie         int
	DragonBlack int
	DragonRed   int
	DragonOdd   int
	DragonEven  int
	TigerBlack  int
	TigerRed    int
	TigerOdd    int
	TigerEven   int
}
type HistoryResultType1 []int32

//type2
type BetType2Data struct {
	Owner1 BetType2DataOwner
	Owner2 BetType2DataOwner
	Owner3 BetType2DataOwner
}
type BetType2DataOwner struct {
	WinEqual  int
	WinTimes  int
	LoseEqual int
	LoseTimes int
}
type HistoryResultType2 [][]int32

//type 6
type BetType6Data struct {
	Small  int
	Big    int
	Odd    int
	Even   int
	Sum    []int
	Dice   []int //點數
	Triple []int
	Pair   []int
	Paigow []int
}
type HistoryResultType6 struct {
	HallID   int   `json:"hallID"`
	RoomID   int   `json:"roomID"`
	Dice     []int `json:"dice"`
	Sum      int   `json:"sum"`
	BigSmall int   `json:"bigSmall"`
	OddEven  int   `json:"oddEven"`
}

//type7
type BetType7Data struct {
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
type HistoryResultType7 []int32

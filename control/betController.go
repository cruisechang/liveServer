package control

import (
	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
)

const (
	type6DiceLen   = 6
	type6PairLen   = 6
	type6TripleLen = 7
	type6SumLen    = 14
	type6PaigowLen = 15

	type7OneLen    = 37
	type7TwoLen    = 60
	type7ThreeLen  = 14
	type7FourLen   = 23
	type7SixLen    = 11
	type7ColumnLen = 3
	type7DozenLen  = 3
)

type BetController struct {
	conf config.Configurer
}

func NewBetController(conf config.Configurer) *BetController {
	return &BetController{
		conf: conf,
	}
}
func CmdDataToRoomDataType0(cmdBetData *config.BetType0CmdData) (*roomConf.BetType0Data) {

	return &roomConf.BetType0Data{
		Banker:      cmdBetData.Banker,
		Player:      cmdBetData.Player,
		Tie:         cmdBetData.Tie,
		BankerPair:  cmdBetData.BankerPair,
		PlayerPair:  cmdBetData.PlayerPair,
		Big:         cmdBetData.Big,
		Small:       cmdBetData.Small,
		AnyPair:     cmdBetData.AnyPair,
		PerfectPair: cmdBetData.PerfectPair,
		SuperSix:    cmdBetData.SuperSix,
	}
}

func CmdDataToRoomDataType1(cmdBet *config.BetType1CmdData) *roomConf.BetType1Data {

	return &roomConf.BetType1Data{
		Dragon:      cmdBet.Dragon,
		Tiger:       cmdBet.Tiger,
		Tie:         cmdBet.Tie,
		DragonBlack: cmdBet.DragonBlack,
		DragonRed:   cmdBet.DragonRed,
		DragonOdd:   cmdBet.DragonOdd,
		DragonEven:  cmdBet.DragonEven,
		TigerBlack:  cmdBet.TigerBlack,
		TigerRed:    cmdBet.TigerRed,
		TigerOdd:    cmdBet.TigerOdd,
		TigerEven:   cmdBet.TigerEven,
	}
}

func CmdDataToRoomDataType2(cmdBet *config.BetType2CmdData) *roomConf.BetType2Data {

	return &roomConf.BetType2Data{

		Owner1: roomConf.BetType2DataOwner{LoseEqual: cmdBet.Owner1.LoseEqual, LoseTimes: cmdBet.Owner1.LoseTimes, WinEqual: cmdBet.Owner1.WinEqual, WinTimes: cmdBet.Owner1.WinTimes},
		Owner2: roomConf.BetType2DataOwner{LoseEqual: cmdBet.Owner2.LoseEqual, LoseTimes: cmdBet.Owner2.LoseTimes, WinEqual: cmdBet.Owner2.WinEqual, WinTimes: cmdBet.Owner2.WinTimes},
		Owner3: roomConf.BetType2DataOwner{LoseEqual: cmdBet.Owner3.LoseEqual, LoseTimes: cmdBet.Owner3.LoseTimes, WinEqual: cmdBet.Owner3.WinEqual, WinTimes: cmdBet.Owner3.WinTimes},
	}

}

func CmdDataToRoomDataType6(cmdBet *config.BetType6CmdData) *roomConf.BetType6Data {

	return &roomConf.BetType6Data{
		Small:  cmdBet.Small,
		Big:    cmdBet.Big,
		Odd:    cmdBet.Odd,
		Even:   cmdBet.Even,
		Sum:    cmdBet.Sum,
		Dice:   cmdBet.Dice,
		Triple: cmdBet.Triple,
		Pair:   cmdBet.Pair,
		Paigow: cmdBet.Paigow,
	}
}

func CmdDataToRoomDataType7(cmdBet *config.BetType7CmdData) *roomConf.BetType7Data {

	return &roomConf.BetType7Data{
		One:    cmdBet.One,
		Two:    cmdBet.Two,
		Three:  cmdBet.Three,
		Four:   cmdBet.Four,
		Six:    cmdBet.Six,
		Column: cmdBet.Column,
		Dozen:  cmdBet.Dozen,
		Big:    cmdBet.Big,
		Small:  cmdBet.Small,
		Odd:    cmdBet.Odd,
		Even:   cmdBet.Even,
		Red:    cmdBet.Red,
		Black:  cmdBet.Black,
	}
}

//CheckBetType0 checks if bet is correct
func CheckBetType0(bet *roomConf.BetType0Data) bool {

	if bet.Banker < 0 {
		return false
	}
	if bet.Player < 0 {
		return false
	}
	if bet.Tie < 0 {
		return false
	}
	if bet.BankerPair < 0 {
		return false
	}
	if bet.PlayerPair < 0 {
		return false
	}
	if bet.Big < 0 {
		return false
	}
	if bet.Small < 0 {
		return false
	}
	if bet.AnyPair < 0 {
		return false
	}
	if bet.PerfectPair < 0 {
		return false
	}
	if bet.SuperSix < 0 {
		return false
	}
	return true
}
func CheckBetType1(bet *roomConf.BetType1Data) bool {

	if bet.Dragon < 0 {
		return false
	}
	if bet.Tiger < 0 {
		return false
	}
	if bet.Tie < 0 {
		return false
	}
	if bet.DragonBlack < 0 {
		return false
	}
	if bet.DragonRed < 0 {
		return false
	}
	if bet.DragonOdd < 0 {
		return false
	}
	if bet.DragonEven < 0 {
		return false
	}
	if bet.TigerBlack < 0 {
		return false
	}
	if bet.TigerRed < 0 {
		return false
	}
	if bet.TigerOdd < 0 {
		return false
	}
	if bet.TigerEven < 0 {
		return false
	}
	return true
}
func CheckBetType2(bet *roomConf.BetType2Data) bool {

	if bet.Owner1.LoseTimes < 0 {
		return false
	}
	if bet.Owner1.LoseEqual < 0 {
		return false
	}
	if bet.Owner1.WinTimes < 0 {
		return false
	}
	if bet.Owner1.WinEqual < 0 {
		return false
	}

	if bet.Owner2.LoseTimes < 0 {
		return false
	}
	if bet.Owner2.LoseEqual < 0 {
		return false
	}
	if bet.Owner2.WinTimes < 0 {
		return false
	}
	if bet.Owner2.WinEqual < 0 {
		return false
	}

	if bet.Owner3.LoseTimes < 0 {
		return false
	}
	if bet.Owner3.LoseEqual < 0 {
		return false
	}
	if bet.Owner3.WinTimes < 0 {
		return false
	}
	if bet.Owner3.WinEqual < 0 {
		return false
	}

	return true
}
func CheckBetType6(bet *roomConf.BetType6Data) bool {

	if bet.Big < 0 {
		return false
	}
	if bet.Small < 0 {
		return false
	}
	if bet.Odd < 0 {
		return false
	}
	if bet.Even < 0 {
		return false
	}
	// check length
	if len(bet.Dice) != type6DiceLen {
		return false
	}
	if len(bet.Pair) != type6PairLen {
		return false
	}
	if len(bet.Sum) != type6SumLen {
		return false
	}
	if len(bet.Triple) != type6TripleLen {
		return false
	}
	if len(bet.Paigow) != type6PaigowLen {
		return false
	}
	//

	for _, v := range bet.Dice {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Pair {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Triple {
		if v < 0 {
			return false
		}
	}

	for _, v := range bet.Sum {
		if v < 0 {
			return false
		}
	}

	for _, v := range bet.Paigow {
		if v < 0 {
			return false
		}
	}
	return true
}
func CheckBetType7(bet *roomConf.BetType7Data) bool {

	if bet.Big < 0 {
		return false
	}
	if bet.Small < 0 {
		return false
	}
	if bet.Odd < 0 {
		return false
	}
	if bet.Even < 0 {
		return false
	}
	if bet.Red < 0 {
		return false
	}
	if bet.Black < 0 {
		return false
	}

	// check length
	if len(bet.One) != type7OneLen {
		return false
	}

	if len(bet.Two) != type7TwoLen {
		return false
	}
	if len(bet.Three) != type7ThreeLen {
		return false
	}
	if len(bet.Four) != type7FourLen {
		return false
	}
	if len(bet.Six) != type7SixLen {
		return false
	}
	if len(bet.Column) != type7ColumnLen {
		return false
	}
	if len(bet.Dozen) != type7DozenLen {
		return false
	}

	//check number
	for _, v := range bet.One {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Two {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Three {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Four {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Six {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Column {
		if v < 0 {
			return false
		}
	}
	for _, v := range bet.Dozen {
		if v < 0 {
			return false
		}
	}



	return true
}

//CountBetSumType0 counts bet sum
func CountBetSumType0(bet *roomConf.BetType0Data) (sum int, err error) {

	return bet.Banker +
		bet.Player +
		bet.Tie +
		bet.BankerPair +
		bet.PlayerPair +
		bet.Big +
		bet.Small +
		bet.AnyPair +
		bet.PerfectPair +
		bet.SuperSix, nil
}

func CountBetSumType1(bet *roomConf.BetType1Data) (sum int, err error) {

	return bet.Dragon +
		bet.Tiger +
		bet.Tie +
		bet.DragonBlack +
		bet.DragonRed +
		bet.DragonOdd +
		bet.DragonEven +
		bet.TigerBlack +
		bet.TigerRed +
		bet.TigerOdd +
		bet.TigerEven, nil
}

func CountBetSumType2(bet *roomConf.BetType2Data) (sum int, err error) {
	owner1 := bet.Owner1
	owner2 := bet.Owner2
	owner3 := bet.Owner3
	sum += owner1.LoseEqual + owner1.WinEqual
	sum += owner2.LoseEqual + owner2.WinEqual
	sum += owner3.LoseEqual + owner3.WinEqual

	//翻倍要夠五倍才行
	sum += (owner1.LoseTimes + owner1.WinTimes +
		owner2.LoseTimes + owner2.WinTimes +
		owner3.LoseTimes + owner3.WinTimes) * 5
	return sum, nil
}
func CountBetSumType6(bet *roomConf.BetType6Data) (sum int, err error) {

	s := bet.Big + bet.Small + bet.Odd + bet.Even

	//count dice
	for _, v := range bet.Dice {
		s += v
	}
	//count pair
	for _, v := range bet.Pair {
		s += v
	}
	for _, v := range bet.Triple {
		s += v
	}

	//count sum
	for _, v := range bet.Sum {
		s += v
	}

	//count paigow
	for _, v := range bet.Paigow {
		s += v
	}

	return s, nil
}

//
func CountBetSumType7(bet *roomConf.BetType7Data) (sum int, err error) {

	for _, v := range bet.One {
		sum += v
	}
	for _, v := range bet.Two {
		sum += v
	}
	for _, v := range bet.Three {
		sum += v
	}
	for _, v := range bet.Four {
		sum += v
	}
	for _, v := range bet.Six {
		sum += v
	}
	for _, v := range bet.Column {
		sum += v
	}
	for _, v := range bet.Dozen {
		sum += v
	}

	sum += bet.Big
	sum += bet.Small
	sum += bet.Odd
	sum += bet.Even
	sum += bet.Red
	sum += bet.Black

	return sum, nil
}

//Add up
//AddBetUpType0 add bet and original bet up
//要以bet為準，因為oriBet可能是空
func AddBetUpType0(bet *roomConf.BetType0Data, oriBet *roomConf.BetType0Data) (*roomConf.BetType0Data, error) {

	return &roomConf.BetType0Data{
		Banker:      bet.Banker + oriBet.Banker,
		Player:      bet.Player + oriBet.Player,
		Tie:         bet.Tie + oriBet.Tie,
		BankerPair:  bet.BankerPair + oriBet.BankerPair,
		PlayerPair:  bet.PlayerPair + oriBet.PlayerPair,
		Big:         bet.Big + oriBet.Big,
		Small:       bet.Small + oriBet.Small,
		AnyPair:     bet.AnyPair + oriBet.AnyPair,
		PerfectPair: bet.PerfectPair + oriBet.PerfectPair,
		SuperSix:    bet.SuperSix + oriBet.SuperSix,
		Commission:  bet.Commission,
	}, nil
}
func AddBetUpType1(bet *roomConf.BetType1Data, oriBet *roomConf.BetType1Data) (*roomConf.BetType1Data, error) {

	return &roomConf.BetType1Data{
		Dragon:      bet.Dragon + oriBet.Dragon,
		Tiger:       bet.Tiger + oriBet.Tiger,
		Tie:         bet.Tie + oriBet.Tie,
		DragonBlack: bet.DragonBlack + oriBet.DragonBlack,
		DragonRed:   bet.DragonRed + oriBet.DragonRed,
		DragonOdd:   bet.DragonOdd + oriBet.DragonOdd,
		DragonEven:  bet.DragonEven + oriBet.DragonEven,
		TigerBlack:  bet.TigerBlack + oriBet.TigerBlack,
		TigerRed:    bet.TigerRed + oriBet.TigerRed,
		TigerOdd:    bet.TigerOdd + oriBet.TigerOdd,
		TigerEven:   bet.TigerEven + oriBet.TigerEven,
	}, nil
}
func AddBetUpType2(bet *roomConf.BetType2Data, oriBet *roomConf.BetType2Data) (*roomConf.BetType2Data, error) {

	return &roomConf.BetType2Data{
		Owner1: roomConf.BetType2DataOwner{LoseEqual: bet.Owner1.LoseEqual + oriBet.Owner1.LoseEqual, LoseTimes: bet.Owner1.LoseTimes + oriBet.Owner1.LoseTimes, WinEqual: bet.Owner1.WinEqual + oriBet.Owner1.WinEqual, WinTimes: bet.Owner1.WinTimes + oriBet.Owner1.LoseTimes,},
		Owner2: roomConf.BetType2DataOwner{LoseEqual: bet.Owner1.LoseEqual + oriBet.Owner1.LoseEqual, LoseTimes: bet.Owner1.LoseTimes + oriBet.Owner1.LoseTimes, WinEqual: bet.Owner1.WinEqual + oriBet.Owner1.WinEqual, WinTimes: bet.Owner1.WinTimes + oriBet.Owner1.LoseTimes,},
		Owner3: roomConf.BetType2DataOwner{LoseEqual: bet.Owner1.LoseEqual + oriBet.Owner1.LoseEqual, LoseTimes: bet.Owner1.LoseTimes + oriBet.Owner1.LoseTimes, WinEqual: bet.Owner1.WinEqual + oriBet.Owner1.WinEqual, WinTimes: bet.Owner1.WinTimes + oriBet.Owner1.LoseTimes,},
	}, nil
}

func AddBetUpType6(bet *roomConf.BetType6Data, oriBet *roomConf.BetType6Data) (*roomConf.BetType6Data, error) {

	dice := make([]int, type6DiceLen)
	pair := make([]int, type6PairLen)
	triple := make([]int, type6TripleLen)
	sum := make([]int, type6SumLen)
	paigow := make([]int, type6PaigowLen)

	//count dice
	for i, v := range bet.Dice {
		dice[i] = v + oriBet.Dice[i]
	}
	//count pair
	for i, v := range bet.Pair {
		pair[i] = v + oriBet.Pair[i]
	}
	for i, v := range bet.Triple {
		triple[i] = v + oriBet.Triple[i]
	}

	//count sum
	for i, v := range bet.Sum {
		sum[i] = v + oriBet.Sum[i]
	}

	//count paigow
	for i, v := range bet.Paigow {
		paigow[i] = v + oriBet.Paigow[i]
	}

	return &roomConf.BetType6Data{
		Small:  bet.Small + oriBet.Small,
		Big:    bet.Big + oriBet.Big,
		Odd:    bet.Odd + oriBet.Odd,
		Even:   bet.Even + oriBet.Even,
		Dice:   dice,
		Pair:   pair,
		Triple: triple,
		Sum:    sum,
		Paigow: paigow,
	}, nil
}
func AddBetUpType7(bet *roomConf.BetType7Data, oriBet *roomConf.BetType7Data) (*roomConf.BetType7Data, error) {

	one := make([]int, type7OneLen)
	two := make([]int, type7TwoLen)
	three := make([]int, type7ThreeLen)
	four := make([]int, type7FourLen)
	six := make([]int, type7SixLen)
	column := make([]int, type7ColumnLen)
	dozen := make([]int, type7DozenLen)

	for i, v := range bet.One {
		one[i] = v + oriBet.One[i]
	}
	for i, v := range bet.Two {
		two[i] = v + oriBet.Two[i]
	}
	for i, v := range bet.Three {
		three[i] = v + oriBet.Three[i]
	}
	for i, v := range bet.Four {
		four[i] = v + oriBet.Four[i]
	}
	for i, v := range bet.Six {
		six[i] = v + oriBet.Six[i]
	}
	for i, v := range bet.Column {
		column[i] = v + oriBet.Column[i]
	}
	for i, v := range bet.Dozen {
		dozen[i] = v + oriBet.Dozen[i]
	}

	return &roomConf.BetType7Data{
		Big:    bet.Big + oriBet.Big,
		Small:  bet.Small + oriBet.Small,
		Odd:    bet.Odd + oriBet.Odd,
		Even:   bet.Even + oriBet.Even,
		Red:    bet.Red + oriBet.Red,
		Black:  bet.Black + oriBet.Black,
		One:    one,
		Two:    two,
		Three:  three,
		Four:   four,
		Six:    six,
		Column: column,
		Dozen:  dozen,
	}, nil
}



package control

import (
	"sort"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	pb "github.com/cruisechang/liveServer/protobuf"
	"github.com/juju/errors"
)

var (
	//type0
	pRate0BankerComission    float32 = 0.95
	pRate0BankerNoCommission float32 = 0.5
	pRate0Player             float32 = 1
	pRate0Tie                float32 = 8
	pRate0BankerPair         float32 = 11
	pRate0PlayerPair         float32 = 11
	pRate0Big                float32 = 0.5
	pRate0Small              float32 = 1.5
	pRate0AnyPair            float32 = 5
	pRate0PerfectPair        float32 = 20
	pRate0Super6             float32 = 12

	//type1
	pRate1Dragon      float32 = 1
	pRate1Tiger       float32 = 1
	pRate1Tie         float32 = 8
	pRate1DragonOdd   float32 = 0.75
	pRate1DragonEven  float32 = 1.05
	pRate1TigerOdd    float32 = 0.75
	pRate1TigerEven   float32 = 1.05
	pRate1DragonRed   float32 = 0.9
	pRate1DragonBlack float32 = 0.9
	pRate1TigerRed    float32 = 0.9
	pRate1TigerBlack  float32 = 0.9

	//type2
	//壓贏，若我贏，看自己牌型決定賠率
	//壓贏，若我輸，看莊家牌型決定賠率
	//壓輸，若我贏，看自己牌型決定賠率
	//壓輸，若我輸，看莊家牌型決定賠率
	//都是看贏家牌型就對了

	//壓中時，贏家牌型，決定倍率
	pRate2HitWinEqual float32 = 0.95 //閒贏 平倍
	pRate2Hit5Face    float32 = 4.75 // 5公 翻倍
	pRate2HitNiuNiu   float32 = 2.85 // 牛牛 翻倍
	pRate2HitNiu789   float32 = 1.90 // 牛789 翻倍
	pRate2HitOther    float32 = 0.95 // 牛123456 無牛 翻倍

	//押錯時，贏家牌型，決定輸的倍數
	pRate2NoHit5Face  float32 = 5 //押翻倍 錯時 贏家的牌型 5公
	pRate2NoHitNiuNiu float32 = 3 //押錯時 贏家的牌型  牛牛
	pRate2NoHitNiu789 float32 = 2 //押錯時 贏家的牌型 牛789
	pRate2NoHitOther  float32 = 1 //押錯時 贏家的牌型 其他

	//type6
	pRate6BigSmall    float32 = 1
	pRate6OddEven     float32 = 1
	pRate6Sum0417     float32 = 50 //sum 4 or 17
	pRate6Sum0516     float32 = 18 //sum 5 or 16
	pRate6Sum0615     float32 = 14 //sum 6,15
	pRate6Sum0714     float32 = 12 //sum 7,14
	pRate6Sum0813     float32 = 8  //sum 8,13
	pRate6Sum09101112 float32 = 6  //sum 9,10,11,12
	pRate6OneDice     float32 = 1  //
	//pRate6TwoDice     float32 = 2   //
	//pRate6ThreeDice   float32 = 3   //
	pRate6Pair      float32 = 8   //
	pRate6Paigow    float32 = 5   //
	pRate6Triple    float32 = 150 //3顆骰子點數都一樣
	pRate6TripleAll float32 = 24  //

	//type7
	pRate7One      float32 = 35
	pRate7Two      float32 = 17
	pRate7Three    float32 = 11
	pRate7Four     float32 = 8
	pRate7Six      float32 = 5
	pRate7Column   float32 = 2
	pRate7Dozen    float32 = 2
	pRate7BigSmall float32 = 1
	pRate7OddEven  float32 = 1
	pRate7RedBlack float32 = 1

	type7TwoValue = [][]int{{0, 1}, {0, 2}, {0, 3},
		{1, 2}, {2, 3},
		{1, 4}, {2, 5}, {3, 6},
		{4, 5}, {5, 6},
		{4, 7}, {5, 8}, {6, 9},
		{7, 8}, {8, 9},
		{7, 10}, {8, 11}, {9, 12},
		{10, 11}, {11, 12},
		{10, 13}, {11, 14}, {12, 15},
		{13, 14}, {14, 15},
		{13, 16}, {14, 17}, {15, 18},
		{16, 17}, {17, 18},
		{16, 19}, {17, 20}, {18, 21},
		{19, 20}, {20, 21},
		{19, 22}, {20, 23}, {21, 24},
		{22, 23}, {23, 24},
		{22, 25}, {23, 26}, {24, 27},
		{25, 26}, {26, 27},
		{25, 28}, {26, 29}, {27, 30},
		{28, 29}, {29, 30},
		{28, 31}, {29, 32}, {30, 33},
		{31, 32}, {32, 33},
		{31, 34}, {32, 35}, {33, 36},
	}
	type7ThreeValue = [][]int{{0, 1, 2}, {0, 2, 3},
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
		{13, 14, 15},
		{16, 17, 18},
		{19, 20, 21},
		{22, 23, 24},
		{25, 26, 27},
		{28, 29, 30},
		{31, 32, 33},
		{34, 35, 36},
	}

	type7FourValue = [][]int{{0, 1, 2, 3},
		{1, 2, 4, 5}, {2, 3, 5, 6},
		{4, 5, 7, 8}, {5, 6, 8, 9},
		{7, 8, 10, 11}, {8, 9, 11, 12},
		{10, 11, 13, 14}, {11, 12, 14, 15},
		{13, 14, 16, 17}, {14, 15, 17, 18},
		{16, 17, 19, 20}, {17, 18, 20, 21},
		{19, 20, 22, 23}, {20, 21, 23, 24},
		{22, 23, 25, 26}, {23, 24, 26, 27},
		{25, 26, 28, 29}, {26, 27, 29, 30},
		{28, 29, 31, 32}, {29, 30, 32, 33},
		{31, 32, 34, 35}, {32, 33, 35, 36},
	}
	type7SixValue = [][]int{
		{1, 2, 3, 4, 5, 6},
		{4, 5, 6, 7, 8, 9},
		{7, 8, 9, 10, 11, 12},
		{10, 11, 12, 13, 14, 15},
		{13, 14, 15, 16, 17, 18},
		{16, 17, 18, 19, 20, 21},
		{19, 20, 21, 22, 23, 24},
		{22, 23, 24, 25, 26, 27},
		{25, 26, 27, 28, 29, 30},
		{28, 29, 30, 31, 32, 33},
		{31, 32, 33, 34, 35, 36},
	}
	typeColumnValue = [][]int{
		{1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 34},
		{2, 5, 8, 11, 14, 17, 20, 23, 26, 29, 32, 35},
		{3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 33, 36},
	}
	typeDozenValue = [][]int{
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
		{25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36},
	}
)

type RateController struct {
	conf config.Configurer
}

func NewRateController(conf config.Configurer) *RateController {
	return &RateController{
		conf: conf,
	}
}

//winLose 0=lose,1=win, 2=tie, 3=no bet 見protocol-game userResult
//resultCredit 輸贏金額
//newCredit = balanceCredit
func (c *RateController) Count(roomType int, oriCredit float32, betData, resultData interface{}) (betCredit, activeCredit int, prizeCredit, resultCredit, balanceCredit float32, winLose int, err error) {
	err = errors.New("count error")

	switch roomType {
	case 0:

		if b, ok := betData.(*roomConf.BetType0Data); ok {
			if r, ok := resultData.(*pb.RoundResultType0Data); ok {
				betCredit, activeCredit, prizeCredit, balanceCredit = c.countType0(oriCredit, b, r)
				resultCredit = balanceCredit - oriCredit
				winLose = 2
				if resultCredit > 0 {
					winLose = 1
				} else if resultCredit < 0 {
					winLose = 0
				}
				err = nil
			}
		}

	case 1:
		if b, ok := betData.(*roomConf.BetType1Data); ok {
			if r, ok := resultData.(*pb.RoundResultType1Data); ok {
				betCredit, activeCredit, prizeCredit, balanceCredit = c.countType1(oriCredit, b, r)
				resultCredit = balanceCredit - oriCredit
				winLose = 2
				if resultCredit > 0 {
					winLose = 1
				} else if resultCredit < 0 {
					winLose = 0
				}
				err = nil
			}
		}

	case 2:
		if b, ok := betData.(*roomConf.BetType2Data); ok {
			if r, ok := resultData.(*pb.RoundResultType2Data); ok {
				betCredit, activeCredit, prizeCredit, balanceCredit = c.countType2(oriCredit, b, r)
				resultCredit = balanceCredit - oriCredit
				winLose = 2
				if resultCredit > 0 {
					winLose = 1
				} else if resultCredit < 0 {
					winLose = 0
				}
				err = nil
			}
		}

	case 6:
		if b, ok := betData.(*roomConf.BetType6Data); ok {
			if r, ok := resultData.(*pb.RoundResultType6Data); ok {
				betCredit, activeCredit, prizeCredit, balanceCredit = c.countType6(oriCredit, b, r)
				resultCredit = balanceCredit - oriCredit
				winLose = 2
				if resultCredit > 0 {
					winLose = 1
				} else if resultCredit < 0 {
					winLose = 0
				}
				err = nil
			}
		}
	case 7:
		if b, ok := betData.(*roomConf.BetType7Data); ok {
			if r, ok := resultData.(*pb.RoundResultType7Data); ok {
				betCredit, activeCredit, prizeCredit, balanceCredit = c.countType7(oriCredit, b, r)
				resultCredit = balanceCredit - oriCredit
				winLose = 2
				if resultCredit > 0 {
					winLose = 1
				} else if resultCredit < 0 {
					winLose = 0
				}
				err = nil
			}
		}
	default:
		return
	}

	return
}

//type0
func (c *RateController) countType0(credit float32, betData *roomConf.BetType0Data, resultData *pb.RoundResultType0Data) (betCredit, activeCredit int, payout, balanceCredit float32) {

	if resultData.Result == 0 {

		if betData.Banker > 0 {
			//免傭
			if betData.Commission == 0 {
				if resultData.BankerPoint == 6 {
					payout += float32(betData.Banker)*pRate0BankerNoCommission + float32(betData.Banker)
				} else {
					payout += float32(betData.Banker + betData.Banker)
				}
			} else {
				payout += float32(betData.Banker)*pRate0BankerComission + float32(betData.Banker)
			}
		}

	} else if resultData.Result == 1 && betData.Player > 0 {
		payout += float32(betData.Player + betData.Player)

	} else if resultData.Result == 2 {

		//退回莊閒押金
		payout += float32(betData.Banker + betData.Player)

		if betData.Tie > 0 {
			payout += float32(betData.Tie)*pRate0Tie + float32(betData.Tie)
		}
	}

	if resultData.BankerPair > 0 && betData.BankerPair > 0 {
		payout += float32(betData.BankerPair)*pRate0BankerPair + float32(betData.BankerPair)
	}

	if resultData.PlayerPair > 0 && betData.PlayerPair > 0 {
		payout += float32(betData.PlayerPair)*pRate0PlayerPair + float32(betData.PlayerPair)
	}

	if resultData.BigSmall == 0 {
		if betData.Small > 0 {
			payout += float32(betData.Small)*pRate0Small + float32(betData.Small)
		}
	} else if resultData.BigSmall == 1 {
		if betData.Big > 0 {
			payout += float32(betData.Big)*pRate0Big + float32(betData.Big)
		}
	}

	if betData.SuperSix > 0 && resultData.SuperSix > 0 {
		payout += float32(betData.SuperSix)*pRate0Super6 + float32(betData.SuperSix)
	}

	if betData.AnyPair > 0 && resultData.AnyPair > 0 {
		payout += float32(betData.AnyPair)*pRate0AnyPair + float32(betData.AnyPair)
	}

	if betData.PerfectPair > 0 && resultData.PerfectPair > 0 {
		payout += float32(betData.PerfectPair)*pRate0PerfectPair + float32(betData.PerfectPair)
	}

	betCredit = c.countBetSum0(betData)
	activeCredit = c.countActiveBet0(betData)
	balanceCredit = credit - float32(betCredit) + payout

	return betCredit, activeCredit, payout, balanceCredit

}

func (c *RateController) countBetSum0(bet *roomConf.BetType0Data) int {

	return bet.Banker +
		bet.Player +
		bet.Tie +
		bet.BankerPair +
		bet.PlayerPair +
		bet.Big +
		bet.Small +
		bet.AnyPair +
		bet.PerfectPair +
		bet.SuperSix
}

func (c *RateController) countActiveBet0(bet *roomConf.BetType0Data) int {

	return bet.Banker +
		bet.Player +
		bet.Tie +
		bet.BankerPair +
		bet.PlayerPair +
		bet.Big +
		bet.Small +
		bet.AnyPair +
		bet.PerfectPair +
		bet.SuperSix
}

//type1
func (c *RateController) countType1(credit float32, betData *roomConf.BetType1Data, resultData *pb.RoundResultType1Data) (betCredit, activeCredit int, prizeCredit, balanceCredit float32) {

	if resultData.Result == 0 && betData.Dragon > 0 {
		prizeCredit += float32(betData.Dragon)*pRate1Dragon + float32(betData.Dragon)

	} else if resultData.Result == 1 && betData.Tiger > 0 {
		prizeCredit += float32(betData.Tiger)*pRate1Tiger + float32(betData.Tiger)

	} else if resultData.Result == 2 {
		//開和 壓龍 虎 返回一半
		if betData.Tie > 0 {
			prizeCredit += float32(betData.Tie)*pRate1Tie + float32(betData.Tie)

		}
		if betData.Dragon > 0 {
			prizeCredit += float32(betData.Dragon / 2)
		}
		if betData.Tiger > 0 {
			prizeCredit += float32(betData.Tiger / 2)
		}
	}

	//單雙
	if resultData.DragonOddEven == 1 {
		if betData.DragonOdd > 0 {
			prizeCredit += float32(betData.DragonOdd)*pRate1DragonOdd + float32(betData.DragonOdd)
		}
	} else if resultData.DragonOddEven == 0 {
		if betData.DragonEven > 0 {
			prizeCredit += float32(betData.DragonEven)*pRate1DragonEven + float32(betData.DragonEven)
		}
	}
	//紅黑
	if resultData.DragonRedBlack == 0 {
		if betData.DragonRed > 0 {
			prizeCredit += float32(betData.DragonRed)*pRate1DragonRed + float32(betData.DragonRed)
		}
	} else if resultData.DragonRedBlack == 1 {
		if betData.DragonBlack > 0 {
			prizeCredit += float32(betData.DragonBlack)*pRate1DragonBlack + float32(betData.DragonBlack)
		}
	}

	if resultData.TigerOddEven == 1 {
		if betData.TigerOdd > 0 {
			prizeCredit += float32(betData.TigerOdd)*pRate1TigerOdd + float32(betData.TigerOdd)
		}
	} else if resultData.TigerOddEven == 0 && betData.TigerEven > 0 {
		prizeCredit += float32(betData.TigerEven)*pRate1TigerEven + float32(betData.TigerEven)
	}

	if resultData.TigerRedBlack == 0 {
		if betData.TigerRed > 0 {
			prizeCredit += float32(betData.TigerRed)*pRate1TigerRed + float32(betData.TigerRed)
		} else if resultData.TigerRedBlack == 1 && betData.TigerBlack > 0 {
			prizeCredit += float32(betData.TigerBlack)*pRate1TigerBlack + float32(betData.TigerBlack)
		}
	}
	betCredit = c.countBetSum1(betData)
	activeCredit = c.countActiveBet1(betData)
	balanceCredit = credit - float32(betCredit) + prizeCredit

	return

}
func (c *RateController) countBetSum1(bet *roomConf.BetType1Data) int {

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
		bet.TigerEven
}
func (c *RateController) countActiveBet1(bet *roomConf.BetType1Data) int {

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
		bet.TigerEven
}

//type2
func (c *RateController) countType2(credit float32, betData *roomConf.BetType2Data, resultData *pb.RoundResultType2Data) (betCredit, activeCredit int, prizeCredit, balanceCredit float32) {

	payout0, loseout0 := c.countType2Real(betData.Owner1, resultData.Owner0, resultData.Owner1)
	payout1, loseout1 := c.countType2Real(betData.Owner2, resultData.Owner0, resultData.Owner2)
	payout2, loseout2 := c.countType2Real(betData.Owner3, resultData.Owner0, resultData.Owner3)

	prizeCredit = payout0 + payout1 + payout2
	balanceCredit = credit - (loseout0 + loseout1 + loseout2) + (prizeCredit)

	betCredit = c.countBetSum2(betData)
	activeCredit = c.countActiveBet2(betData)

	return

}
func (c *RateController) countBetSum2(bet *roomConf.BetType2Data) int {

	return bet.Owner1.LoseEqual +
		bet.Owner1.LoseTimes +
		bet.Owner1.WinEqual +
		bet.Owner1.WinTimes +

		bet.Owner2.LoseEqual +
		bet.Owner2.LoseTimes +
		bet.Owner2.WinEqual +
		bet.Owner2.WinTimes +

		bet.Owner3.LoseEqual +
		bet.Owner3.LoseTimes +
		bet.Owner3.WinEqual +
		bet.Owner3.WinTimes
}

func (c *RateController) countActiveBet2(bet *roomConf.BetType2Data) int {

	return bet.Owner1.LoseEqual +
		bet.Owner1.LoseTimes +
		bet.Owner1.WinEqual +
		bet.Owner1.WinTimes +

		bet.Owner2.LoseEqual +
		bet.Owner2.LoseTimes +
		bet.Owner2.WinEqual +
		bet.Owner2.WinTimes +

		bet.Owner3.LoseEqual +
		bet.Owner3.LoseTimes +
		bet.Owner3.WinEqual +
		bet.Owner3.WinTimes
}

//betOwner 指定index投注者投注資料
//resultBanker 莊家輸贏資料
//resultPlayer 指定index閒家輸贏資料

func (c *RateController) countType2Real(betOwner roomConf.BetType2DataOwner, resultBanker *pb.RoundResultType2Owner, resultPlayer *pb.RoundResultType2Owner) (payout, loseout float32) {

	//算沒下中的
	//player贏
	if resultPlayer.Result == 1 {

		//下player輸
		if betOwner.LoseEqual > 0 {

			//下 輸 平倍
			loseout += float32(betOwner.LoseEqual)
		}
		//下輸 翻倍
		if betOwner.LoseTimes > 0 {

			//贏家牌型其他
			if resultBanker.Pattern < int32(c.conf.PatternType2Niu7()) {

				loseout += float32(betOwner.LoseTimes) * pRate2NoHitOther

				//贏家牌型牛牛
			} else if resultBanker.Pattern == int32(c.conf.PatternType2NiuNiu()) {
				loseout += float32(betOwner.LoseTimes) * pRate2NoHitNiuNiu

				//贏家5公
			} else if resultBanker.Pattern == int32(c.conf.PatternType25Face()) {
				loseout += float32(betOwner.LoseTimes) * pRate2NoHit5Face

				//贏家牛789
			} else {
				loseout += float32(betOwner.LoseTimes) * pRate2NoHitNiu789
			}
		}
	} else {
		//player輸

		//下贏平倍
		if betOwner.WinEqual > 0 {

			//下 閒贏 平倍
			loseout += float32(betOwner.WinEqual)
		}

		//下贏 翻倍
		if betOwner.WinTimes > 0 {

			//贏家牌型其他
			if resultBanker.Pattern < int32(c.conf.PatternType2Niu7()) {

				loseout += float32(betOwner.WinTimes) * pRate2NoHitOther

				//贏家牌型牛牛
			} else if resultBanker.Pattern == int32(c.conf.PatternType2NiuNiu()) {
				loseout += float32(betOwner.WinTimes) * pRate2NoHitNiuNiu

				//贏家5公
			} else if resultBanker.Pattern == int32(c.conf.PatternType25Face()) {
				loseout += float32(betOwner.WinTimes) * pRate2NoHit5Face

				//贏家牛789
			} else {
				loseout += float32(betOwner.WinTimes) * pRate2NoHitNiu789
			}
		}

	}

	//算下中的

	//target player  win
	if resultPlayer.Result == 1 {

		//有下 贏 平倍
		if betOwner.WinEqual > 0 {
			payout += float32(betOwner.WinEqual) * pRate2HitWinEqual
		}

		//有下 贏 翻倍
		if betOwner.WinTimes > 0 {

			//5公
			if resultPlayer.Pattern == int32(c.conf.PatternType25Face()) {
				payout += float32(betOwner.WinTimes) * pRate2Hit5Face

				//牛牛
			} else if resultPlayer.Pattern == int32(c.conf.PatternType2NiuNiu()) {
				payout += float32(betOwner.WinTimes) * pRate2HitNiuNiu

				//牛789
			} else if resultPlayer.Pattern >= int32(c.conf.PatternType2Niu7()) {
				payout += float32(betOwner.WinTimes) * pRate2HitNiu789

				//其他
			} else {
				payout += float32(betOwner.WinTimes) * pRate2HitOther

			}
		}
		//owner1 lose
	} else {

		//有下 輸 平倍
		if betOwner.LoseEqual > 0 {
			//平倍
			payout += float32(betOwner.LoseEqual) * pRate2HitWinEqual
		}

		//有下 輸 翻倍
		if betOwner.LoseTimes > 0 {

			//莊5公
			if resultBanker.Pattern == int32(c.conf.PatternType25Face()) {
				payout += float32(betOwner.LoseTimes) * pRate2Hit5Face

				//莊牛牛
			} else if resultBanker.Pattern == int32(c.conf.PatternType2NiuNiu()) {
				payout += float32(betOwner.LoseTimes) * pRate2HitNiuNiu

				//莊牛789
			} else if resultBanker.Pattern >= int32(c.conf.PatternType2Niu7()) {
				payout += float32(betOwner.LoseTimes) * pRate2HitNiu789

				//莊其他
			} else {
				payout += float32(betOwner.LoseTimes) * pRate2HitOther

			}
		}

	}

	return
}

//type6
func (c *RateController) countType6(credit float32, betData *roomConf.BetType6Data, resultData *pb.RoundResultType6Data) (betCredit, activeCredit int, prizeCredit, balanceCredit float32) {

	//var payout float32 = 0

	//大小  1=big, 0=small, -1=圍骰通殺
	if resultData.BigSmall == 1 && betData.Big > 0 {
		prizeCredit += float32(betData.Big)*pRate6BigSmall + float32(betData.Big)
	} else if resultData.BigSmall == 0 && betData.Small > 0 {
		prizeCredit += float32(betData.Small)*pRate6BigSmall + float32(betData.Small)
	}

	//單雙 1=單 0=雙 -1=圍骰通殺
	if resultData.OddEven == 1 && betData.Odd > 0 {
		prizeCredit += float32(betData.Odd)*pRate6OddEven + float32(betData.Odd)
	} else if resultData.OddEven == 0 && betData.Even > 0 {
		prizeCredit += float32(betData.Even)*pRate6OddEven + float32(betData.Even)
	}

	//總和 3 / 18 沒有總和點數, 用index 算
	if resultData.Sum != 3 && resultData.Sum != 18 {
		//sum 的位置有投注，s=投多少
		if s := betData.Sum[resultData.Sum-4]; s > 0 {
			pr := c.getType6SumPayRate(int(resultData.Sum))

			prizeCredit += float32(s)*pr + float32(s)
		}
	}
	//單骰
	for _, v := range resultData.Dice {
		//這個 index 有投注
		if s := betData.Dice[v-1]; s > 0 {
			prizeCredit += float32(s)*pRate6OneDice + float32(s)
		}
	}

	//對子 有對子，且該位置index 有下注
	if resultData.Pair > 0 {
		if s := betData.Pair[resultData.Pair-1]; s > 0 {
			prizeCredit += float32(s)*pRate6Pair + float32(s)
		}
	}

	//牌九 有牌九
	if len(resultData.Paigow) > 0 {
		for _, v := range resultData.Paigow {

			i := c.getType6PaigowToBetDataIndex(v.Result)

			//這個index 有投注
			if s := betData.Paigow[i]; s > 0 {
				prizeCredit += float32(s)*pRate6Paigow + float32(s)
			}
		}
	}

	//圍骰
	if resultData.Triple > 0 {
		//有下全圍
		if alls := betData.Triple[6]; alls > 0 {
			prizeCredit += float32(alls)*pRate6TripleAll + float32(alls)
		}
		//有中單個圍骰

		if s := betData.Triple[resultData.Triple-1]; s > 0 {
			prizeCredit += float32(s)*pRate6Triple + float32(s)
		}
	}

	betCredit = c.countBetSum6(betData)
	activeCredit = c.countActiveBet6(betData)
	balanceCredit = credit - float32(betCredit) + prizeCredit

	return

}

func (c *RateController) countBetSum6(bet *roomConf.BetType6Data) int {

	ssum := 0
	for _, s := range bet.Sum {
		ssum += s
	}
	for _, s := range bet.Dice {
		ssum += s
	}
	for _, s := range bet.Triple {
		ssum += s
	}
	for _, s := range bet.Pair {
		ssum += s
	}
	for _, s := range bet.Paigow {
		ssum += s
	}

	return bet.Big +
		bet.Small +
		bet.Odd +
		bet.Even +
		ssum
}
func (c *RateController) countActiveBet6(bet *roomConf.BetType6Data) int {

	ssum := 0
	for _, s := range bet.Sum {
		ssum += s
	}
	for _, s := range bet.Dice {
		ssum += s
	}
	for _, s := range bet.Triple {
		ssum += s
	}
	for _, s := range bet.Pair {
		ssum += s
	}
	for _, s := range bet.Paigow {
		ssum += s
	}

	return bet.Big +
		bet.Small +
		bet.Odd +
		bet.Even +
		ssum
}

func (c *RateController) getType6SumPayRate(sum int) float32 {
	if sum == 4 || sum == 17 {
		return pRate6Sum0417
	}
	if sum == 5 || sum == 16 {
		return pRate6Sum0516
	}
	if sum == 6 || sum == 15 {
		return pRate6Sum0615
	}
	if sum == 7 || sum == 14 {
		return pRate6Sum0714
	}
	if sum == 8 || sum == 13 {
		return pRate6Sum0813
	}
	if sum == 9 || sum == 10 || sum == 11 || sum == 12 {
		return pRate6Sum09101112
	}
	return 0
}
func (c *RateController) getType6PaigowToBetDataIndex(paigow []int32) int {

	//[]int32 to []int
	pg := []int{}
	for _, v := range paigow {
		pg = append(pg, int(v))
	}
	sort.Ints(pg)
	f := pg[0]
	s := pg[1]
	if f == 1 {
		switch s {
		case 2:
			return 0
		case 3:
			return 1
		case 4:
			return 2
		case 5:
			return 3
		case 6:
			return 4
		}
	}
	if f == 2 {
		switch s {
		case 3:
			return 5
		case 4:
			return 6
		case 5:
			return 7
		case 6:
			return 8
		}
	}
	if f == 3 {
		switch s {
		case 4:
			return 9
		case 5:
			return 10
		case 6:
			return 11
		}
	}
	if f == 4 {
		switch s {
		case 5:
			return 12
		case 6:
			return 13
		}
	}
	if f == 5 && s == 6 {
		return 14
	}

	return 14

}

//type7
func (c *RateController) countType7(credit float32, betData *roomConf.BetType7Data, resultData *pb.RoundResultType7Data) (betCredit, activeCredit int, prizeCredit, balanceCredit float32) {

	//var payout float32 = 0
	var result = int(resultData.Result)

	//大小  1=big, 0=small  -1=沒有=點數0
	if resultData.BigSmall == 1 && betData.Big > 0 {
		prizeCredit += float32(betData.Big)*pRate7BigSmall + float32(betData.Big)
	} else if resultData.BigSmall == 0 && betData.Small > 0 {
		prizeCredit += float32(betData.Small)*pRate7BigSmall + float32(betData.Small)
	}

	//單雙 ,  -1=沒有=點數0
	if resultData.OddEven == 1 && betData.Odd > 0 {
		prizeCredit += float32(betData.Odd)*pRate7OddEven + float32(betData.Odd)
	} else if resultData.OddEven == 0 && betData.Even > 0 {
		prizeCredit += float32(betData.Even)*pRate7OddEven + float32(betData.Even)
	}

	//紅黑  -1=沒有=點數0
	if resultData.RedBlack == 0 && betData.Red > 0 {
		prizeCredit += float32(betData.Red)*pRate7RedBlack + float32(betData.Red)

	} else if resultData.RedBlack == 1 && betData.Black > 0 {
		prizeCredit += float32(betData.Black)*pRate7RedBlack + float32(betData.Black)
	}

	//打 0=1打, 1=2打 ,2=3打 ,-1=沒有=點數0
	if resultData.Dozen == 0 && betData.Dozen[0] > 0 {
		prizeCredit += float32(betData.Dozen[0])*pRate7Dozen + float32(betData.Dozen[0])

	} else if resultData.Dozen == 1 && betData.Dozen[1] > 0 {
		prizeCredit += float32(betData.Dozen[1])*pRate7Dozen + float32(betData.Dozen[1])

	} else if resultData.Dozen == 2 && betData.Dozen[2] > 0 {
		prizeCredit += float32(betData.Dozen[2])*pRate7Dozen + float32(betData.Dozen[2])

	}

	//列 0=1列, 1=2列 ,2=3列 ,-1=沒有=點數0
	if resultData.Column == 0 && betData.Column[0] > 0 {
		prizeCredit += float32(betData.Column[0])*pRate7Column + float32(betData.Column[0])

	} else if resultData.Column == 1 && betData.Column[1] > 0 {
		prizeCredit += float32(betData.Column[1])*pRate7Column + float32(betData.Column[1])

	} else if resultData.Column == 2 && betData.Column[2] > 0 {
		prizeCredit += float32(betData.Column[2])*pRate7Column + float32(betData.Column[2])
	}

	//點數 1號碼  resultData.Result=號碼
	one := betData.One[resultData.Result]
	if one > 0 {
		prizeCredit += float32(one)*pRate7One + float32(one)

	}

	//2號碼
	prizeCredit += c.countType7TwoPayout(betData, result)
	//3號碼
	prizeCredit += c.countType7ThreePayout(betData, result)
	//4號碼
	prizeCredit += c.countType7FourPayout(betData, result)
	//6號碼
	prizeCredit += c.countType7SixPayout(betData, result)

	betCredit = c.countBetSum7(betData)
	activeCredit = c.countActiveBet7(betData)
	balanceCredit = credit - float32(betCredit) + prizeCredit

	return
}

func (c *RateController) countBetSum7(bet *roomConf.BetType7Data) int {

	ssum := 0
	for _, s := range bet.One {
		ssum += s
	}
	for _, s := range bet.Two {
		ssum += s
	}
	for _, s := range bet.Three {
		ssum += s
	}
	for _, s := range bet.Four {
		ssum += s
	}
	for _, s := range bet.Six {
		ssum += s
	}
	for _, s := range bet.Column {
		ssum += s
	}
	for _, s := range bet.Dozen {
		ssum += s
	}

	return bet.Big +
		bet.Small +
		bet.Odd +
		bet.Even +
		bet.Red +
		bet.Black +
		ssum
}

func (c *RateController) countActiveBet7(bet *roomConf.BetType7Data) int {

	ssum := 0
	for _, s := range bet.One {
		ssum += s
	}
	for _, s := range bet.Two {
		ssum += s
	}
	for _, s := range bet.Three {
		ssum += s
	}
	for _, s := range bet.Four {
		ssum += s
	}
	for _, s := range bet.Six {
		ssum += s
	}
	for _, s := range bet.Column {
		ssum += s
	}
	for _, s := range bet.Dozen {
		ssum += s
	}

	return bet.Big +
		bet.Small +
		bet.Odd +
		bet.Even +
		bet.Red +
		bet.Black +
		ssum
}

//計算type7 two 投注結果 ，two index 要看企劃案
func (c *RateController) countType7TwoPayout(betData *roomConf.BetType7Data, result int) (payout float32) {

	resultIndexes := c.countType7ResultIndexTwo(result)

	//用resultIndex去算，最多只有四次
	for _, v := range resultIndexes {
		//v=result two index
		//如果結果的index有投注
		//計算金額
		if bet := betData.Two[v]; bet > 0 {
			payout += float32(bet)*pRate7Two + float32(bet)
		}
	}

	return

}

//type7 下注法 three 投注結果
func (c *RateController) countType7ThreePayout(betData *roomConf.BetType7Data, result int) (payout float32) {

	resultIndexes := c.countType7ResultIndexThree(result)

	//用resultIndex去算
	for _, v := range resultIndexes {
		if bet := betData.Three[v]; bet > 0 {
			payout += float32(bet)*pRate7Three + float32(bet)
		}
	}
	return
}

//type7 下注法 four 投注結果
func (c *RateController) countType7FourPayout(betData *roomConf.BetType7Data, result int) (payout float32) {

	resultIndexes := c.countType7ResultIndexFour(result)

	//用resultIndex去算
	for _, v := range resultIndexes {
		if bet := betData.Four[v]; bet > 0 {
			payout += float32(bet)*pRate7Four + float32(bet)
		}
	}
	return
}

//type7 下注法 four 投注結果
func (c *RateController) countType7SixPayout(betData *roomConf.BetType7Data, result int) (payout float32) {

	resultIndexes := c.countType7ResultIndexSix(result)

	//用resultIndex去算
	for _, v := range resultIndexes {
		if bet := betData.Six[v]; bet > 0 {
			payout += float32(bet)*pRate7Six + float32(bet)

		}
	}
	return
}

//取得type7 two 的result 的index ，最多4個index
func (c *RateController) countType7ResultIndexTwo(result int) (indexes []int) {
	n := 0
	for i, v := range type7TwoValue {
		if v[0] == result {
			indexes = append(indexes, i)
			n += 1
		} else if v[1] == result {
			indexes = append(indexes, i)
			n += 1
		}
		if n >= 4 {
			break
		}
	}

	return indexes
}

//取得type7 three 的result 的index ，最多3個index , 2號
func (c *RateController) countType7ResultIndexThree(result int) (indexes []int) {

	if result <= 3 {
		if result == 0 {
			return []int{0, 1}
		}
		if result == 1 {
			return []int{0, 2}
		}
		if result == 2 {
			return []int{0, 1, 2}
		} else {
			//result ==3
			return []int{1, 2}
		}
	}

	//result >3
	for i, v := range type7ThreeValue {
		if result >= v[0] && result <= v[2] {
			indexes = append(indexes, i)
			return
		}
	}

	return indexes
}

//取得type7 four 的result 的index ，最多4個index
func (c *RateController) countType7ResultIndexFour(result int) (indexes []int) {

	if result == 0 {
		return []int{0}
	}
	//餘數==2，最多4個index，其餘最多2個index
	max := 2 //最大次數
	remainder := result % 3
	if remainder == 2 {
		max = 4
	}

	count := 0 //次數
	for i, v := range type7FourValue {

		if v[0] == result {
			indexes = append(indexes, i)
			count += 1
		} else if v[1] == result {
			indexes = append(indexes, i)
			count += 1
		} else if v[2] == result {
			indexes = append(indexes, i)
			count += 1
		} else if v[3] == result {
			indexes = append(indexes, i)
			count += 1
		}

		if count >= max {
			break
		}
	}

	return indexes
}

//取得type7 six 的result 的index ，最多2個index
func (c *RateController) countType7ResultIndexSix(result int) (indexes []int) {

	//result 0 沒中
	if result == 0 {
		return
	}
	max := 2   //index最大次數
	count := 0 //index次數
	for i, v := range type7SixValue {

		if result >= v[0] && result <= v[5] {
			indexes = append(indexes, i)
			count += 1
		}

		if count >= max {
			break
		}
	}

	return indexes
}

//取得type7 column 的result 的index ，最多1個index
func (c *RateController) countType7ResultIndexColumn(result int) int {

	//result 0 沒中
	if result == 0 {
		return -1
	}

	//用餘數判斷第幾列
	if remainder := result % 3; remainder == 0 {
		//第三列
		return 2
	} else {
		//其餘兩列
		return remainder - 1
	}

}

//取得type7 dozen 的result 的index ，最多1個index
func (c *RateController) countType7ResultIndexDozen(result int) int {

	if result <= 12 {

		//result 0 沒中
		if result == 0 {
			return -1
		}

		return 0

	} else if result <= 24 {
		return 1
	} else {
		return 2
	}

}

package fake

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/cruisechang/liveServer/config/roomConf"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	nexSpace "github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

//加新game
//main加上需要的
//createRoomField
//createHistoryResult
//createTypeData

//room
type RoomField struct {
	RoomID       int    `json:"roomID"`
	Name         string `json:"name"`
	Type         int    `json:"type"`
	Active       int    `json:"active"`
	HLSURL       string `json:"hlsURL"`
	Dealer       roomConf.Dealer
	BetCoundDown int `json:"betCoundDown"`

	TypeData      interface{} //type data for different type room,  bet limit , boot, round....
	ResultHistory interface{} //  result history
	Boot          int64
	Round         int64
}

func CreateRoomField(id int, name string, typ int) RoomField {
	return RoomField{
		RoomID: id,
		Name:   name,
		Type:   typ,
		//
		Active: 0,
		HLSURL: "google.com",
		Dealer: roomConf.Dealer{
			DealerID:    100,
			Name:        "Alisa",
			PortraitURL: "yahoo.com",
		},
		BetCoundDown:  20,
		TypeData:      createTypeData(typ),
		ResultHistory: createHistoryResult(typ),
	}
}

func createHistoryResult(roomType int) interface{} {
	switch roomType {
	case 0:
		return roomConf.HistoryResultType0{
			{2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1},
			{2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0},
			{2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0},
		}
	case 1:
		return roomConf.HistoryResultType1{
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1,
		}
	case 2:
		return roomConf.HistoryResultType2{
			{1, 1, 0, 1, 0, 1, 1, 1}, {1, 1, 0, 2, 1, 2, 1, 2}, {1, 1, 0, 3, 0, 3, 1, 3}, {1, 1, 1, 4, 0, 4, 1, 4}, {1, 1, 0, 5, 0, 5, 1, 5},
			{1, 1, 0, 6, 1, 6, 1, 6}, {1, 1, 1, 7, 1, 7, 1, 7}, {1, 1, 0, 8, 1, 8, 1, 8}, {1, 1, 1, 9, 0, 9, 1, 9}, {1, 1, 0, 10, 0, 10, 1, 10},
			{1, 1, 1, 11, 0, 11, 1, 11},
		}
	case 6:
		var hrt6 []*roomConf.HistoryResultType6
		hrt6 = append(hrt6,
			&roomConf.HistoryResultType6{
				HallID:   0,
				RoomID:   1,
				Dice:     []int{1, 2, 2},
				Sum:      5,
				BigSmall: 0,
				OddEven:  1,
			})

		hrt6 = append(hrt6,
			&roomConf.HistoryResultType6{
				HallID:   0,
				RoomID:   1,
				Dice:     []int{1, 2, 3},
				Sum:      6,
				BigSmall: 0,
				OddEven:  0,
			})

		hrt6 = append(hrt6,
			&roomConf.HistoryResultType6{
				HallID:   0,
				RoomID:   1,
				Dice:     []int{4, 2, 3},
				Sum:      9,
				BigSmall: 1,
				OddEven:  1,
			})

		return hrt6
	case 7:
		return roomConf.HistoryResultType7{1, 0, 10, 33, 25, 13, 22, 1, 9, 26, 30, 30, 30}
	default:
		return nil
	}
}
func createTypeData(roomType int) interface{} {
	switch roomType {
	case 0:
		return &roomConf.TypeData0{
			Boot:             10,
			Round:            28,
			BetLimit:         []int{1, 10000},
			BankerLimit:      []int{1, 10000},
			PlayerLimit:      []int{1, 10000},
			TieLimit:         []int{1, 1000},
			BankerPairLimit:  []int{1, 5000},
			PlayerPairLimit:  []int{1, 5000},
			AnyPairLimit:     []int{1, 1000},
			PerfectPairLimit: []int{1, 1000},
			SuperSixLimit:    []int{1, 1000},
			BigSmallLimit:    []int{1, 1000},
		}
	case 1:
		return &roomConf.TypeData1{
			Boot:          10,
			Round:         48,
			BetLimit:      []int{1, 50000},
			DragonLimit:   []int{1, 50000},
			TigerLimit:    []int{1, 50000},
			TieLimit:      []int{1, 50000},
			OddEvenLimit:  []int{1, 50000},
			RedBlackLimit: []int{1, 50000},
		}
	case 2:
		return &roomConf.TypeData2{
			Boot:      10,
			Round:     10,
			BetLimit:  []int{1, 50000},
			WinTimes:  []int{1, 50000},
			WinEqual:  []int{1, 50000},
			LoseTimes: []int{1, 50000},
			LoseEqual: []int{1, 50000},
		}
	case 6:
		return &roomConf.TypeData6{
			Boot:     1,
			Round:    1,
			BetLimit: []int{1, 50000},

			BigSmallLimit:    []int{1, 50000},
			OddEvenLimit:     []int{1, 50000},
			Sum0417Limit:     []int{1, 1000},
			Sum0516Limit:     []int{1, 3000},
			Sum0615Limit:     []int{1, 3500},
			Sum0714Limit:     []int{1, 4000},
			Sum0813Limit:     []int{1, 7000},
			Sum09101112Limit: []int{1, 9000},

			DiceLimit:      []int{1, 20000},
			PairLimit:      []int{1, 7000},
			PaigowLimit:    []int{1, 10000},
			TripleLimit:    []int{1, 500},
			TripleAllLimit: []int{1, 20000},
		}
	case 7:
		return &roomConf.TypeData7{
			Boot:     10,
			Round:    10,
			BetLimit: []int{1, 50000},

			OneLimit:      []int{1, 50000},
			TwoLimit:      []int{1, 50000},
			ThreeLimit:    []int{1, 50000},
			FourLimit:     []int{1, 50000},
			SixLimit:      []int{1, 50000},
			ColumnLimit:   []int{1, 50000},
			DozenLimit:    []int{1, 50000},
			BigSmallLimit: []int{1, 50000},

			OddEvenLimit:  []int{1, 50000},
			RedBlackLimit: []int{1, 50000},
		}

	default:
		return nil
	}
}
func CreateRoom(nex nexSpace.Nex, rh *roomCtrl.Controller, data RoomField) entity.Room {

	rm := nex.GetRoomManager()
	r, _ := rm.CreateRoom(data.RoomID, data.Type, data.Name)
	//r.SetType(data.Type)

	rh.SetHLSURL(r, data.HLSURL)
	rh.SetDealer(r, &roomConf.Dealer{data.Dealer.DealerID, data.Dealer.Name, data.Dealer.PortraitURL})
	rh.SetBetCountdown(r, data.BetCoundDown)
	rh.SetTypeData(r, data.TypeData)
	rh.SetHistoryResult(r, data.ResultHistory)

	return r
}

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

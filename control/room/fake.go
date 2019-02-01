package room

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/cruisechang/liveServer/config/roomConf"
	nexSpace "github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

//room
type roomField struct {
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

func createRoomField(id int, name string, typ int) roomField {
	return roomField{
		RoomID: id,
		Name:   name,
		Type:   typ,
		//
		Active: 1,
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
		return roomConf.HistoryResultType0{{0, 1, 1}, {1, 0, 0}, {2, 0, 1}}
	case 1:
		return roomConf.HistoryResultType1{0, 1, 2, 0, 1, 2, 0, 1, 2, 1}
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
	default:
		return nil
	}
	return nil
}
func createTypeData(roomType int) interface{} {
	switch roomType {
	case 0:
		return &roomConf.TypeData0{
			Boot:             1,
			Round:            1,
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
			1,
			1,
			[]int{1, 50000},
			[]int{1, 50000},
			[]int{1, 50000},
			[]int{1, 50000},
			[]int{1, 50000},
			[]int{1, 50000},
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
			Boot:          1,
			Round:         1,
			BetLimit:      []int{1, 50000},
			OneLimit:      []int{1, 10000},
			TwoLimit:      []int{1, 10000},
			ThreeLimit:    []int{1, 10000},
			FourLimit:     []int{1, 10000},
			SixLimit:      []int{1, 10000},
			ColumnLimit:   []int{1, 10000},
			DozenLimit:    []int{1, 10000},
			BigSmallLimit: []int{1, 10000},
			OddEvenLimit:  []int{1, 10000},
			RedBlackLimit: []int{1, 10000},
		}

	default:
		return nil
	}
}
func createRoom(nex nexSpace.Nex, rh *Controller, data roomField) entity.Room {

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

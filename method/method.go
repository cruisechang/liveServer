package method

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/dbConf"
	"github.com/cruisechang/liveServer/config/roomConf"
	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	nexSpace "github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func QueryHalls(bp *control.DBController, method, path string) ([]*dbConf.HallData, error) {

	logPrefix := "QueryHalls"

	resp, err := bp.Do(method, path, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, errors.New(fmt.Sprintf("%s http status code wrong %d", logPrefix, resp.StatusCode))
	}
	if resp.Body == nil {
		return nil, errors.New(fmt.Sprintf("%s response body==nil", logPrefix))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	resData := &dbConf.ResponseHallData{}
	err = json.Unmarshal(body, resData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s unmarshal responseData error=%s", logPrefix, err.Error()))
	}
	return resData.Data, nil
}

func QueryRooms(bp *control.DBController, method, path string) ([]*dbConf.RoomData, error) {
	logPrefix := "QueryRooms "

	resp, err := bp.Do(method, path, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, errors.New(fmt.Sprintf("%s http status code wrong %d", logPrefix, resp.StatusCode))
	}
	if resp.Body == nil {
		return nil, errors.New(fmt.Sprintf("%s response body==nil", logPrefix))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	resData := &dbConf.ResponseRoomData{}
	err = json.Unmarshal(body, resData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s unmarshal responseData error=%s", logPrefix, err.Error()))
	}
	return resData.Data, nil
}

func QueryLimitations(bp *control.DBController, method, path string) ([]*dbConf.Limitation, error) {

	logPrefix := "QueryLimitations "

	resp, err := bp.Do(method, path, nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, errors.New(fmt.Sprintf("%s http status code wrong %d", logPrefix, resp.StatusCode))
	}
	if resp.Body == nil {
		return nil, errors.New(fmt.Sprintf("%s response body==nil", logPrefix))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	resData := &dbConf.ResponseLimitation{}
	err = json.Unmarshal(body, resData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s unmarshal responseData error=%s", logPrefix, err.Error()))
	}
	return resData.Data, nil
}

func TransferLimitations(limitations []*dbConf.Limitation) (data []*config.TransferredLimitation, err error) {

	res := []*config.TransferredLimitation{}

	for _, v := range limitations {
		switch v.LimitationID {

		case config.LimitationID0:
			d := &config.Limitation0{}
			err := json.Unmarshal([]byte(v.Limitation), d)
			if err != nil {
				return nil, err
			}

			s := &config.TransferredLimitation{
				LimitationID: v.LimitationID,
				Limitation:   d,
			}
			res = append(res, s)

		case config.LimitationID100:
			d := &config.Limitation100{}
			err := json.Unmarshal([]byte(v.Limitation), d)
			if err != nil {
				return nil, err
			}
			s := &config.TransferredLimitation{
				LimitationID: v.LimitationID,
				Limitation:   d,
			}
			res = append(res, s)
		case config.LimitationID200:
			d := &config.Limitation200{}
			err := json.Unmarshal([]byte(v.Limitation), d)
			if err != nil {
				return nil, err
			}
			s := &config.TransferredLimitation{
				LimitationID: v.LimitationID,
				Limitation:   d,
			}
			res = append(res, s)
		case config.LimitationID600:
			d := &config.Limitation600{}
			err := json.Unmarshal([]byte(v.Limitation), d)
			if err != nil {
				return nil, err
			}
			s := &config.TransferredLimitation{
				LimitationID: v.LimitationID,
				Limitation:   d,
			}
			res = append(res, s)
		case config.LimitationID700:
			d := &config.Limitation700{}
			err := json.Unmarshal([]byte(v.Limitation), d)
			if err != nil {
				return nil, err
			}
			s := &config.TransferredLimitation{
				LimitationID: v.LimitationID,
				Limitation:   d,
			}
			res = append(res, s)
		default:
			return nil, errors.New("TransferLimitations switch not found")
		}
	}
	return res, nil
}

func QueryDealers(bp *control.DBController, method, path string) ([]*dbConf.Dealer, error) {

	logPrefix := "QueryDealers "

	resp, err := bp.Do(method, path, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, errors.New(fmt.Sprintf("%s http status code wrong %d", logPrefix, resp.StatusCode))
	}
	if resp.Body == nil {
		return nil, errors.New(fmt.Sprintf("%s response body==nil", logPrefix))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	resData := &dbConf.DealerData{}
	err = json.Unmarshal(body, resData)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s unmarshal responseData error=%s", logPrefix, err.Error()))
	}
	return resData.Data, nil
}

//加新game
//main加上需要的
//createRoomField
//createHistoryResult
//createTypeData

//room
type RoomField struct {
	HallID       uint   `json:"hallID"`
	RoomID       uint   `json:"roomID"`
	Name         string `json:"name"`
	Type         uint   `json:"type"`
	Active       uint   `json:"active"`
	HLSURL       string `json:"hlsURL"`
	Dealer       roomConf.Dealer
	BetCoundDown uint `json:"betCoundDown"`

	TypeData      interface{} //type data for different type room,  bet limit , boot, round....
	ResultHistory interface{} //  result history
	Boot          int64
	Round         int64
}

func GetLimitation(limitationID int, limitations []*config.TransferredLimitation) interface{} {

	for _, v := range limitations {
		if limitationID == v.LimitationID {

			return v.Limitation
		}
	}
	return nil
}
func GetDealer(dealerID int, dealers []*dbConf.Dealer) (roomConf.Dealer, error) {

	for _, v := range dealers {
		if dealerID == v.DealerID {

			return roomConf.Dealer{
				DealerID:    v.DealerID,
				Name:        v.Name,
				PortraitURL: v.PortraitURL,
			}, nil
		}
	}
	return roomConf.Dealer{}, errors.New("dealer not found")
}

func CreateTypeData(boot int, round int64, limitation interface{}) (interface{}, error) {

	d, ok := limitation.(*config.Limitation0)
	if ok {
		return &roomConf.TypeData0{
			Boot:             boot,
			Round:            round,
			BetLimit:         d.BetLimit,
			BankerLimit:      d.BankerLimit,
			PlayerLimit:      d.PlayerLimit,
			TieLimit:         d.TieLimit,
			BankerPairLimit:  d.BankerPairLimit,
			PlayerPairLimit:  d.PlayerPairLimit,
			AnyPairLimit:     d.AnyPairLimit,
			PerfectPairLimit: d.PerfectPairLimit,
			SuperSixLimit:    d.SuperSixLimit,
			BigSmallLimit:    d.BigSmallLimit,
		}, nil
	}
	d1, ok := limitation.(*config.Limitation100)
	if ok {
		return &roomConf.TypeData1{
			Boot:          boot,
			Round:         round,
			BetLimit:      d1.BetLimit,
			DragonLimit:   d1.DragonLimit,
			TigerLimit:    d1.TigerLimit,
			TieLimit:      d1.TieLimit,
			OddEvenLimit:  d1.OddEvenLimit,
			RedBlackLimit: d1.RedBlackLimit,
		}, nil
	}

	d2, ok := limitation.(*config.Limitation200)
	if ok {

		return &roomConf.TypeData2{
			Boot:           boot,
			Round:          round,
			BetLimit:       d2.BetLimit,
			WinTimesLimit:  d2.WinTimesLimit,
			WinEqualLimit:  d2.WinEqualLimit,
			LoseTimesLimit: d2.LoseTimesLimit,
			LoseEqualLimit: d2.LoseEqualLimit,
		}, nil
	}

	d6, ok := limitation.(*config.Limitation600)
	if ok {
		return &roomConf.TypeData6{
			Boot:     boot,
			Round:    round,
			BetLimit: d6.BetLimit,

			BigSmallLimit:    d6.BigSmallLimit,
			OddEvenLimit:     d6.OddEvenLimit,
			Sum0417Limit:     d6.Sum0417Limit,
			Sum0516Limit:     d6.Sum0516Limit,
			Sum0615Limit:     d6.Sum0615Limit,
			Sum0714Limit:     d6.Sum0813Limit,
			Sum0813Limit:     d6.Sum0813Limit,
			Sum09101112Limit: d6.Sum09101112Limit,

			DiceLimit:      d6.DiceLimit,
			PairLimit:      d6.PairLimit,
			PaigowLimit:    d6.PaigowLimit,
			TripleLimit:    d6.TripleLimit,
			TripleAllLimit: d6.TripleAllLimit,
		}, nil
	}

	d7, ok := limitation.(*config.Limitation700)
	if ok {
		return &roomConf.TypeData7{
			Boot:     boot,
			Round:    round,
			BetLimit: d7.BetLimit,

			OneLimit:      d7.OneLimit,
			TwoLimit:      d7.TwoLimit,
			ThreeLimit:    d7.ThreeLimit,
			FourLimit:     d7.FourLimit,
			SixLimit:      d7.SixLimit,
			ColumnLimit:   d7.ColumnLimit,
			DozenLimit:    d7.DozenLimit,
			BigSmallLimit: d7.BigSmallLimit,

			OddEvenLimit:  d7.OddEvenLimit,
			RedBlackLimit: d7.RedBlackLimit,
		}, nil
	}

	return nil, errors.New("CreateTypeData error")
}
func CreateHistoryResult(roomType int) interface{} {
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

func CreateRoomParameters(data *dbConf.RoomData, typeData interface{}, dealer roomConf.Dealer, historyResult interface{}) *RoomField {
	return &RoomField{
		HallID: data.HallID,
		RoomID: data.RoomID,
		Name:   data.Name,
		Type:   data.RoomType,
		//
		Active:        data.Active,
		HLSURL:        data.HLSURL,
		Dealer:        dealer,
		BetCoundDown:  data.BetCountdown,
		TypeData:      typeData,
		ResultHistory: historyResult,
	}
}
func CreateRoom(rm nexSpace.RoomManager, rh *roomCtrl.Controller, data *RoomField) (entity.Room, error) {

	r, err := rm.CreateRoom(int(data.RoomID), int(data.Type), data.Name)
	if err != nil {
		return nil, err
	}

	r.SetActive(int(data.Active))
	r.SetHallID(int(data.HallID))
	rh.SetHLSURL(r, data.HLSURL)
	rh.SetDealer(r, &roomConf.Dealer{data.Dealer.DealerID, data.Dealer.Name, data.Dealer.PortraitURL})
	rh.SetBetCountdown(r, int(data.BetCoundDown))
	rh.SetTypeData(r, data.TypeData)
	rh.SetHistoryResult(r, data.ResultHistory)

	return r, nil
}

//createRooms
func CreateRooms(rm nexSpace.RoomManager, rCtrl *roomCtrl.Controller, roomData []*dbConf.RoomData, transferredLims []*config.TransferredLimitation, dealers []*dbConf.Dealer) ([]entity.Room, error) {
	var rs []entity.Room

	for _, r := range roomData {

		li := GetLimitation(int(r.LimitationID), transferredLims)
		if li == nil {
			return nil, errors.New(fmt.Sprintf("GetLimitation ==nil,limitationID=%d", r.LimitationID))
		}
		de, err := GetDealer(int(r.DealerID), dealers)

		if err != nil {
			//exit(202, errors.New(fmt.Sprintf("GetDealer ==nil, deaerID=%d", r.DealerID)))

		}

		typeData, err := CreateTypeData(int(r.Boot), int64(r.RoundID), li)
		if err != nil {
			return nil, fmt.Errorf(fmt.Sprintf("CreateTypeData error, roundID=%d", r.RoomID))

		}

		historyResult := CreateHistoryResult(int(r.RoomType))

		param := CreateRoomParameters(r, typeData, de, historyResult)

		rr, err := CreateRoom(rm, rCtrl, param)
		if err != nil {
			return nil, fmt.Errorf("CreateRoom error %s, roomID=%d", err.Error(), r.RoomID)
		}

		//innit boot round
		if rCtrl.NeedInitBootRound(rr) {
			rCtrl.SetRound(rr, 181211100100028)
		}
		rs = append(rs, rr)
	}
	return rs, nil
}

func CreateHalls(hm nexSpace.HallManager, hallData []*dbConf.HallData) {
	for _, h := range hallData {

		ha, _ := hm.CreateHall(int(h.HallID), h.Name)
		ha.SetActive(int(h.Active))
	}
}

func SetHallRoom(halls []entity.Hall, rooms []entity.Room) {
	for _, h := range halls {
		for _, r := range rooms {
			if r.HallID() == h.ID() {
				h.AddRoom(r)
			}
		}
	}
}




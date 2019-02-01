package room

import (
	"errors"
	"reflect"
	"time"

	"github.com/cruisechang/liveServer/config/roomConf"
	"github.com/cruisechang/nex/entity"
)

//type data
func (r *Controller) GetTypeData(room entity.Room) interface{} {
	v, _ := room.GetInterfaceVariable(roomVarTypeData)

	return v
}

//SetTypeData pass a pointer as parameter
func (r *Controller) SetTypeData(room entity.Room, v interface{}) error {

	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return errors.New("parameter is not a pointer")
	}

	room.SetInterfaceVariable(roomVarTypeData, v)
	return nil
}

/*
 boot and round
*/
type BootRoundData struct {
	Boot   int
	Round  int64
	BetMin int
	BetMax int
	Err    error
}

func (r *Controller) GetBootRoundBetMinBetMax(room entity.Room) *BootRoundData {

	data := r.GetTypeData(room)
	res := &BootRoundData{
		Err: errors.New("parameter data error"),
	}

	if typeData, ok := data.(*roomConf.TypeData0); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil

	} else if typeData, ok := data.(*roomConf.TypeData1); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil

	} else if typeData, ok := data.(*roomConf.TypeData2); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil

	} else if typeData, ok := data.(*roomConf.TypeData3); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil

	} else if typeData, ok := data.(*roomConf.TypeData4); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil

	} else if typeData, ok := data.(*roomConf.TypeData5); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil

	} else if typeData, ok := data.(*roomConf.TypeData6); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil

	} else if typeData, ok := data.(*roomConf.TypeData7); ok {
		res.Boot = typeData.Boot
		res.Round = int64(typeData.Round)
		res.BetMin = typeData.BetLimit[0]
		res.BetMax = typeData.BetLimit[1]
		res.Err = nil
	}

	return res

}

func (r *Controller) NeedInitBootRound(room entity.Room) bool {

	data := r.GetBootRoundBetMinBetMax(room)
	if data.Round < 10000 {
		return true
	}
	return false

}

func (r *Controller) InitBootRound(room entity.Room) {

	round := r.initRound(room.ID(), 1)
	r.setBoot(room, 1)
	r.SetRound(room, round)
}

func (r *Controller) setBoot(room entity.Room, boot int) {

	typeData := r.GetTypeData(room)
	r.modifyTypeDataBoot(typeData, boot)
}

func (r *Controller) SetRound(room entity.Room, round int64) {

	typeData := r.GetTypeData(room)
	r.modifyTypeDataRound(typeData, round)
}

//HandleNewBoot count and set new boot to room,
//if room type is type0 or type1 , when got another day, round of the room will be reset too.
func (r *Controller) HandleNewBoot(room entity.Room) (boot int, round int64) {

	brData := r.GetBootRoundBetMinBetMax(room)
	round = brData.Round

	//本日初始round
	ir := r.initRound(room.ID(), 1)
	//本日初始round > old round 表示換日期
	boot = r.countNewBoot(room, brData.Boot, brData.Round, ir)
	round = r.initRound(room.ID(), boot)

	r.setBoot(room, boot)
	r.SetRound(room, round)

	return
}

func (r *Controller) countNewBoot(room entity.Room, oldBoot int, oldRound int64, initRound int64) (newBoot int) {

	//百家 龍虎在換靴時，如果日期沒換，boot直接+1
	//百家 龍虎在換靴時，如果有換日期，reset靴號，局號
	if room.Type() == r.conf.RoomType0() ||
		room.Type() == r.conf.RoomType1() {

		//有換日期，boot init
		if initRound > oldRound {

			newBoot = 1
		} else {
			newBoot = oldBoot + 1
		}
	} else {
		newBoot = oldBoot + 1
	}
	return
}

//HandleNewRound counts and sets new round to the room,
func (r *Controller) HandleNewRound(room entity.Room) (boot int, round int64) {

	brData := r.GetBootRoundBetMinBetMax(room)
	ir := r.initRound(room.ID(), brData.Boot)
	round = r.countNewRound(room, brData.Round, ir)
	r.SetRound(room, round)
	return brData.Boot, round
}

func (r *Controller) countNewRound(room entity.Room, oldRound int64, initRound int64) int64 {

	var newRound int64
	//百家 龍虎 round直接+1
	if room.Type() == r.conf.RoomType0() ||
		room.Type() == r.conf.RoomType1() {
		newRound = oldRound + 1
	} else {
		//第二種，日期更新就reset
		//日期是否換新

		if initRound > oldRound {
			newRound = initRound + 1
		} else {
			newRound = oldRound + 1
		}
	}
	return newRound
}
func (r *Controller) initRound(roomID int, boot int) int64 {
	t := time.Now()

	y := t.Year() % 2000
	dateint := y*10000 + int(t.Month())*100 + t.Day()
	return int64(dateint*1000000000) + int64(roomID*1000000) + int64(boot*10000)
}

func (r *Controller) modifyTypeDataBoot(data interface{}, boot int) (interface{}, error) {

	if typeData, ok := data.(*roomConf.TypeData0); ok {
		typeData.Boot = boot
		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData1); ok {
		typeData.Boot = boot

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData2); ok {
		typeData.Boot = boot

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData3); ok {
		typeData.Boot = boot
		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData4); ok {
		typeData.Boot = boot
		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData5); ok {
		typeData.Boot = boot
		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData6); ok {
		typeData.Boot = boot
		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData7); ok {
		typeData.Boot = boot

		return typeData, nil

	}

	return nil, errors.New("modifyTypeDataBoot error")
}

func (r *Controller) modifyTypeDataRound(data interface{}, round int64) (interface{}, error) {

	if typeData, ok := data.(*roomConf.TypeData0); ok {
		typeData.Round = round

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData1); ok {
		typeData.Round = round

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData2); ok {
		typeData.Round = round

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData3); ok {
		typeData.Round = round

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData4); ok {
		typeData.Round = round

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData5); ok {
		typeData.Round = round

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData6); ok {
		typeData.Round = round

		return typeData, nil

	} else if typeData, ok := data.(*roomConf.TypeData7); ok {
		typeData.Round = round

		return typeData, nil
	}

	return nil, errors.New("modifyTypeDataRound error")
}

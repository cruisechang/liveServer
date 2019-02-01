package room

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/cruisechang/liveServer/config/roomConf"
	"github.com/cruisechang/nex/entity"
)

func (r *Controller) AddHistoryResultType(room entity.Room, data interface{}) error {

	dt, err := room.GetInterfaceVariable(roomHistoryResult)
	if err != nil {
		return err
	}

	sy := sync.RWMutex{}
	sy.Lock()
	defer sy.Unlock()

	switch room.Type() {
	case r.conf.RoomType0():
		m, _ := dt.(roomConf.HistoryResultType0)

		if bd, ok := data.([]int32); ok {

			m = append(m, bd)
			r.SetHistoryResult(room, m)
		} else {
			return errors.New(fmt.Sprintf("passed data error =%s", reflect.TypeOf(data).String()))
		}

	case r.conf.RoomType1():
		m, _ := dt.(roomConf.HistoryResultType1)

		if bd, ok := data.(int32); ok {

			m = append(m, bd)
			r.SetHistoryResult(room, m)
		} else {
			return errors.New(fmt.Sprintf("passed data error =%s", reflect.TypeOf(data).String()))
		}

	case r.conf.RoomType2():
		m, _ := dt.(roomConf.HistoryResultType2)

		if bd, ok := data.([]int32); ok {

			m = append(m, bd)
			r.SetHistoryResult(room, m)
		} else {
			return errors.New(fmt.Sprintf("passed data error =%s", reflect.TypeOf(data).String()))
		}

	case r.conf.RoomType3():

	case r.conf.RoomType4():

	case r.conf.RoomType5():

	case r.conf.RoomType6():
		m, _ := dt.([]*roomConf.HistoryResultType6)

		if bd, ok := data.(*roomConf.HistoryResultType6); ok {

			m = append(m, bd)
			r.SetHistoryResult(room, m)
		} else {
			return errors.New(fmt.Sprintf("passed data error =%s", reflect.TypeOf(data).String()))
		}
	case r.conf.RoomType7():
		m, _ := dt.(roomConf.HistoryResultType7)

		if bd, ok := data.(int32); ok {

			m = append(m, bd)
			r.SetHistoryResult(room, m)
		} else {
			return errors.New(fmt.Sprintf("passed data error =%s", reflect.TypeOf(data).String()))
		}
	default:
		return errors.New(fmt.Sprintf("room type error, type=%d", room.Type()))
	}
	return nil
}

func (r *Controller) SetHistoryResult(room entity.Room, result interface{}) {
	room.SetInterfaceVariable(roomHistoryResult, result)
}

func (r *Controller) GetHistoryResult(room entity.Room) (interface{}, error) {
	return room.GetInterfaceVariable(roomHistoryResult)

}

//各room 的history 裡面
func (r *Controller) SetHistoryResultEmpty(room entity.Room) {

	switch room.Type() {
	case r.conf.RoomType0():
		r.SetHistoryResult(room, roomConf.HistoryResultType0{})
	case r.conf.RoomType1():
		r.SetHistoryResult(room, roomConf.HistoryResultType1{})

	case r.conf.RoomType2():
		r.SetHistoryResult(room, roomConf.HistoryResultType2{})
	case r.conf.RoomType3():

	case r.conf.RoomType4():

	case r.conf.RoomType5():

	case r.conf.RoomType6():
		r.SetHistoryResult(room, []*roomConf.HistoryResultType6{})
	case r.conf.RoomType7():
		r.SetHistoryResult(room, roomConf.HistoryResultType7{})
	}
}


package room

import (
	"errors"
	"sync"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	"github.com/cruisechang/nex/entity"
)

var (
	roomVarBoot         string = "boot"
	roomVarRound        string = "round"
	roomVarHLSUR        string = "hls"
	roomVarBetCountDown string = "betcountdown"
	roomVarDealer       string = "dealer"
	roomVarTypeData     string = "typeData"

	roomHistoryResult string = "historyResult"
	varBet            string = "bet"
)

type Controller struct {
	mutex *sync.RWMutex
	conf  config.Configurer
}

func NewController(conf config.Configurer) *Controller {
	return &Controller{
		mutex: &sync.RWMutex{},
		conf:  conf,
	}
}

func (r *Controller) GetHLSURL(room entity.Room) string {
	hls, _ := room.GetStringVariable(roomVarHLSUR)
	return hls
}
func (r *Controller) SetHLSURL(room entity.Room, hls string) {
	room.SetStringVariable(roomVarHLSUR, hls)
}

func (r *Controller) GetBetCountdown(room entity.Room) int {
	v, _ := room.GetIntVariable(roomVarBetCountDown)
	return v
}
func (r *Controller) SetBetCountdown(room entity.Room, value int) {
	room.SetIntVariable(roomVarBetCountDown, value)
}

//dealer
func (r *Controller) GetDealer(room entity.Room) (dealer roomConf.Dealer, err error) {
	v, _ := room.GetInterfaceVariable(roomVarDealer)

	if de, ok := v.(roomConf.Dealer); ok {
		return de, nil
	}
	return roomConf.Dealer{}, errors.New("GetRoomDealer error")
}

func (r *Controller) SetDealer(room entity.Room, dealer *roomConf.Dealer) {
	room.SetInterfaceVariable(roomVarDealer, dealer)
}

package command

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
	nxLog "github.com/cruisechang/nex/log"
)

type roomInfoProcessor struct {
	BasicProcessor

	roomCtrl *roomCtrl.Controller
}

func NewRoomInfoProcessor(processor BasicProcessor) (*roomInfoProcessor, error) {
	p := &roomInfoProcessor{
		BasicProcessor: processor,
		roomCtrl:       roomCtrl.NewController(processor.GetConfigurer()),
	}

	return p, nil

}

func (p *roomInfoProcessor) Run(obj *nex.CommandObject) error {
	logPrefix := "roomInfoProcessor"
	logger := p.GetLogger()

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s begin", logPrefix))

	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s user==nil", logPrefix))
		return fmt.Errorf("%s user==nil", logPrefix)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdRoomInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s base64 decode cmd data error,user:%s,error:%s", logPrefix, user.Name(), err.Error()))
		return err
	}

	//取 data
	data := []config.RoomInfoCmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdRoomInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json unmarshal cmd data error,user=%s,error=%s, cata=%s", logPrefix, user.Name(), err.Error(), deStr))
		return err
	}

	hm := p.GetHallManager()
	rm := p.GetRoomManager()

	//get rooms
	hallID := data[0].HallID
	roomID := data[0].RoomID
	var hall entity.Hall
	rooms := []entity.Room{}

	//get hall
	if hallID == conf.UnAssigned() {
		p.SendCommand(config.CodeHallIDError, 0, conf.CmdRoomInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s  hallID=%d", logPrefix,hallID))
		return fmt.Errorf("%s hallID error id=%d", logPrefix, hallID)
	}

	//hall id not in hallmanager
	if !hm.ContainHall(hallID) {
		p.SendCommand(config.CodeHallNotFound, 0, conf.CmdRoomInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s  hallManager.ContainHall()!=true,hallID=%d", logPrefix,hallID))
		return fmt.Errorf("%s hallManager.ContainHall(hallID) !=true", logPrefix)
	}

	hall, _ = hm.GetHall(hallID)

	//get all room in this hall
	if roomID == conf.UnAssigned() {
		rooms = hall.GetRooms()
	} else {
		r, _ := rm.GetRoom(roomID)
		rooms = append(rooms, r)
	}

	if len(rooms) == 0 {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s roomManager.GetRoom() len==0 hallID=%d,roomID=%d", logPrefix, hallID, roomID))
	}

	//res data
	resData := []config.RoomInfoResData{}

	for _, r := range rooms {

		typ := r.Type()
		hls := p.roomCtrl.GetHLSURL(r)
		dealer, _ := p.roomCtrl.GetDealer(r)
		betCountDown := p.roomCtrl.GetBetCountdown(r)

		re := config.RoomInfoResData{
			RoomID:       r.ID(),
			Name:         r.Name(),
			Type:         typ,
			Active:       r.Active(),
			HLSURL:       hls,
			Dealer:       dealer,
			BetCoundDown: betCountDown,
		}

		switch typ {
		case conf.RoomType0():
			tpd := p.roomCtrl.GetTypeData(r)
			td, ok := tpd.(*roomConf.TypeData0)
			if !ok {
				return errors.New(fmt.Sprintf("%s getTypeData assertion error roomID=%d, roomType= %d", logPrefix, r.Type(), r.ID()))
			}
			re.TypeData0 = td
		case conf.RoomType1():
			tpd := p.roomCtrl.GetTypeData(r)
			td, ok := tpd.(*roomConf.TypeData1)
			if !ok {
				return errors.New(fmt.Sprintf("%s getTypeData assertion error roomID=%d, roomType= %d", logPrefix, r.Type(), r.ID()))
			}
			re.TypeData1 = td
		case conf.RoomType2():
			tpd := p.roomCtrl.GetTypeData(r)
			td, ok := tpd.(*roomConf.TypeData2)
			if !ok {
				return errors.New(fmt.Sprintf("%s getTypeData assertion error roomID=%d, roomType= %d", logPrefix, r.Type(), r.ID()))
			}
			re.TypeData2 = td
		case conf.RoomType3():
			//tpd, _ := p.roomCtrl.GetRoomTypeData3(r)
			//re.TypeData3 = *tpd
		case conf.RoomType4():
			//tpd, _ := p.roomCtrl.GetRoomTypeData4(r)
			//re.TypeData4 = *tpd
		case conf.RoomType5():
			//tpd, _ := p.roomCtrl.GetRoomTypeData5(r)
			//re.TypeData5 = *tpd
		case conf.RoomType6():
			tpd := p.roomCtrl.GetTypeData(r)
			td, ok := tpd.(*roomConf.TypeData6)
			if !ok {
				return errors.New(fmt.Sprintf("%s getTypeData assertion error roomID=%d, roomType= %d", logPrefix, r.Type(), r.ID()))
			}
			re.TypeData6 = td
		case conf.RoomType7():
			tpd := p.roomCtrl.GetTypeData(r)
			td, ok := tpd.(*roomConf.TypeData7)
			if !ok {
				return errors.New(fmt.Sprintf("%s getTypeData assertion error roomID=%d, roomType= %d", logPrefix, r.Type(), r.ID()))
			}
			re.TypeData7 = td
		}

		resData = append(resData, re)
	}

	b, err := json.Marshal(resData)
	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json marshal res data error,user:%s,error:%s", logPrefix, user.Name(), err.Error()))
		return err
	}

	//[]byte encode to base64 string
	sendData := base64.StdEncoding.EncodeToString(b)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdRoomInfo(), sendData, user, []string{user.ConnID()})
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete  user id=%d,user=%s, resData=%+v ", logPrefix, user.UserID(), user.Name(), resData))

	return nil
}

package command

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
	nxLog "github.com/cruisechang/nex/log"
	"github.com/juju/errors"
)

//loginProcessor implements command.Processor
type resultType7Processor struct {
	BasicProcessor
	roomHandler *roomCtrl.Controller
}

func NewHistoryResultType7Processor(processor BasicProcessor) (*resultType7Processor, error) {
	c := &resultType7Processor{
		BasicProcessor: processor,
		roomHandler:    roomCtrl.NewController(processor.GetConfigurer()),
	}
	return c, nil
}

func (p *resultType7Processor) Run(obj *nex.CommandObject) error {
	conf := p.GetConfigurer()
	logger := p.GetLogger()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("historyResultType7Processor user==nil"))
		return errors.New("historyResultType7Processor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("historyResultType7Processor panic:%v", r))
		}
	}()

	//parsing cmd data
	de, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("historyResultType7Processor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	data := []config.HistoryResultTypeCmdData{}

	if err := json.Unmarshal(de, &data); err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("historyResultType7Processor json Unmarshal cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	hm := p.GetHallManager()
	rm := p.GetRoomManager()

	hallID := data[0].HallID
	roomID := data[0].RoomID

	rooms := []entity.Room{}

	//get hall
	if hallID == conf.UnAssigned() {
		p.SendCommand(config.CodeHallIDError, 0, conf.CmdHistoryResultType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		return errors.New("historyResultType7Processor hall id error")
	}

	//hall id not in hallmanager
	if !hm.ContainHall(hallID) {
		p.SendCommand(config.CodeHallNotFound, 0, conf.CmdHistoryResultType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		return errors.New("historyResultType7Processor hall not found")
	}

	hall, _ := hm.GetHall(hallID)

	//all room in this hall
	if roomID == conf.UnAssigned() {
		rs := hall.GetRooms()

		//get target type room
		for _, v := range rs {
			if r, ok := rm.GetRoom(v.ID()); ok {
				if r.Type() == conf.RoomType7() {
					rooms = append(rooms, r)
				}
			}
		}
	} else {
		//get room and check type
		if r, ok := rm.GetRoom(data[0].RoomID); ok {
			if r.Type() == conf.RoomType7() {
				rooms = append(rooms, r)
			}
		}
	}
	if len(rooms) == 0 {
		p.SendCommand(config.CodeRoomNotFound, 0, conf.CmdHistoryResultType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		return errors.New("historyResultType7Processor room not found")
	}

	var resData []config.HistoryResultType7ResData
	{
	}

	//all rooms
	for _, v := range rooms {

		hrs, err := p.roomHandler.GetHistoryResult(v)

		if err != nil {
			logger.LogFile(nxLog.LevelError, fmt.Sprintf("historyResultType7Processor get history result error=%s ", err.Error()))
			return err
		}

		//this room's history result
		hrt, ok := hrs.(roomConf.HistoryResultType7)
		if !ok {
			logger.LogFile(nxLog.LevelError, fmt.Sprintf("historyResultType7Processor history result assertion error "))
			return errors.New("historyResultType7Processor assertion error")
		}

		resData = append(resData,
			config.HistoryResultType7ResData{
				HallID: hallID,
				RoomID: v.ID(),
				Result: hrt,
			})
	}

	b, err := json.Marshal(resData)
	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("historyResultType7Processor json marshal res data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	sendData := base64.StdEncoding.EncodeToString(b)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdHistoryResultType7(), sendData, user, []string{user.ConnID()})

	//roadMap
	rmToClient := []*config.RoadMapType7ResData{}
	for _, v := range rooms {
		if got, ok := p.RoadMapDataType7(v.ID()); ok {

			rmToClient = append(rmToClient, got)
		}
	}

	if len(rmToClient) > 0 {
		b, err := json.Marshal(rmToClient)
		if err != nil {
			p.DisconnectUser(user.UserID())
			logger.LogFile(nxLog.LevelError, fmt.Sprintf("roadMapType7 json marshal roadmap error,user:%s,error:%s", user.Name(), err.Error()))
			return err
		}

		sendData := base64.StdEncoding.EncodeToString(b)

		p.SendCommand(config.CodeSuccess, 0, conf.CmdRoadMapType7(), sendData, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("roadMapType7 resData=%s ",  string(b)))
	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("historyResultType7Processor complete user id :%d, user=%s,resData=%+v ", user.UserID(), user.Name(), resData))

	return nil
}

package command

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/control"
	ctrlRoom "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
	"github.com/juju/errors"
)

type betType6Processor struct {
	BasicProcessor
	roomCtrl *ctrlRoom.Controller
}

func NewBetType6Processor(processor BasicProcessor) (*betType6Processor, error) {
	p := &betType6Processor{
		BasicProcessor: processor,
		roomCtrl:       ctrlRoom.NewController(processor.GetConfigurer()),
	}

	return p, nil

}

func (p *betType6Processor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor user==nil "))
		return errors.New("betType6Processor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("betType6Processor panic:%v", r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdBetType1(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	//取 data
	var data []config.BetType6CmdData
	{
	}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdBetType6(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor json unmarshal cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	} else if len(data) < 1 {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType1(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor received data error,user:%s", user.Name()))
		return err
	}

	bet := data[0]
	room, ok := p.GetRoomManager().GetRoom(bet.RoomID)
	if !ok {
		return errors.New("get room error")
	}

	//check room active
	if room.Active() != 1 {
		return errors.New("room inactive")
	}

	//check room type
	rTyp := room.Type()
	if rTyp != conf.RoomType6() {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType6(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor received data error,user:%s", user.Name()))
		return err
	}

	//cmd data to room data
	roomBet := control.CmdDataToRoomDataType6(&bet)

	//本次投注總額
	sum, err := control.CountBetSumType6(roomBet)

	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor CountBetSumType6 user:%s, err:%s", user.Name(), err.Error()))
		return errors.New("count user bet sum error")
	}

	//之前投注總額
	oriBet, err := p.roomCtrl.GetUserBetType6(room, user.UserID())
	oriSum, err := control.CountBetSumType6(oriBet)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor CountBetSumType6 ori user:%s, err:%s", user.Name(), err.Error()))
		return errors.New("count user ori bet sum error")
	}
	//檢查總額
	if int(user.Credit()) < (sum + oriSum) {
		p.SendCommand(config.CodeCreditNotEnough, 0, conf.CmdBetType6(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor user credit not enough error,user:%s", user.Name()))
		return errors.New("user credit not enough")
	}

	//room get  user old bet data and add new
	p.roomCtrl.AddBet(room, user.UserID(), roomBet)

	resData := []config.BetType6ResData{
		{
			UserID: user.UserID(),
			RoomID: bet.RoomID,
			Small:  roomBet.Small,
			Big:    roomBet.Big,
			Odd:    roomBet.Odd,
			Even:   roomBet.Even,
			Sum:    roomBet.Sum,
			Dice:   roomBet.Dice,
			Triple: roomBet.Triple,
			Pair:   roomBet.Pair,
			Paigow: roomBet.Paigow,
		},
	}

	//logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("betType6Processor complete  user id=%d,user=%s,bet=%+v", user.UserID(), user.Name(),bet))

	b, err := json.Marshal(resData)
	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType6Processor json marshal res data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	sendData := base64.StdEncoding.EncodeToString(b)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdBetType6(), sendData, user, p.GetRoomReceivers(room.ID()))

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("betType6Processor complete  user id=%d,user=%s, resData=%+v ", user.UserID(), user.Name(), resData))

	return nil
}

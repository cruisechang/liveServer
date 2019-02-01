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

type betType7Processor struct {
	BasicProcessor
	roomCtrl *ctrlRoom.Controller
}

func NewBetType7Processor(processor BasicProcessor) (*betType7Processor, error) {
	p := &betType7Processor{
		BasicProcessor: processor,
		roomCtrl:       ctrlRoom.NewController(processor.GetConfigurer()),
	}

	return p, nil

}

func (p *betType7Processor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor user==nil "))
		return errors.New("betType7Processor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("betType7Processor panic:%v", r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdBetType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	//取 data
	data := []config.BetType7CmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdBetType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor json unmarshal cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	} else if len(data) < 0 {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor received data error,user:%s", user.Name()))
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
	if rTyp != conf.RoomType7() {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor received data error,user:%s", user.Name()))
		return err
	}

	roomBet := control.CmdDataToRoomDataType7(&bet)

	//本次投注總額
	sum, err := control.CountBetSumType7(roomBet)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor CountBetSumType7 user:%s, err:%s", user.Name(), err.Error()))
		return errors.New("count user bet sum error")
	}

	//之前投注總額
	oriBet, err := p.roomCtrl.GetUserBetType7(room, user.UserID())
	oriSum, err := control.CountBetSumType7(oriBet)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor CountBetSumType7 ori user:%s, err:%s", user.Name(), err.Error()))
		return errors.New("count user ori bet sum error")
	}
	//檢查總額
	if int(user.Credit()) < (sum + oriSum) {
		p.SendCommand(config.CodeCreditNotEnough, 0, conf.CmdBetType7(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor user credit not enough error,user:%s", user.Name()))
		return errors.New("user credit not enough")
	}

	//room get  user old bet data and add new
	p.roomCtrl.AddBet(room, user.UserID(), roomBet)

	resData := []config.BetType7ResData{
		{
			UserID: user.UserID(),
			RoomID: bet.RoomID,
			One:    bet.One,
			Two:    bet.Two,
			Three:  bet.Three,
			Four:   bet.Four,
			Six:    bet.Six,
			Column: bet.Column,
			Dozen:  bet.Dozen,
			Big:    bet.Big,
			Small:  bet.Small,
			Odd:    bet.Odd,
			Even:   bet.Even,
			Red:    bet.Red,
			Black:  bet.Black,
		},
	}

	//logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("betType7Processor complete  user id=%d,user=%s,bet=%+v", user.UserID(), user.Name(),bet))

	b, err := json.Marshal(resData)
	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType7Processor json marshal res data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	sendData := base64.StdEncoding.EncodeToString(b)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdBetType7(), sendData, user, p.GetRoomReceivers(room.ID()))

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("betType7Processor complete  user id=%d,user=%s, resData=%+v ", user.UserID(), user.Name(), resData))

	return nil
}

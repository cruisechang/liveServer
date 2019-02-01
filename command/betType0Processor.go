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

type betType0Processor struct {
	BasicProcessor
	roomCtrl *ctrlRoom.Controller
}

func NewBetType0Processor(processor BasicProcessor) (*betType0Processor, error) {
	p := &betType0Processor{
		BasicProcessor: processor,
		roomCtrl:       ctrlRoom.NewController(processor.GetConfigurer()),
	}

	return p, nil

}

func (p *betType0Processor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterHallProcessor user==nil "))
		return errors.New("betType0Processor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("betType0Processor panic:%v", r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	//取 data
	data := []config.BetType0CmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor json unmarshal cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	} else if len(data) < 0 {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor received data error,user:%s", user.Name()))
		return err
	}

	bet := data[0]
	room, ok := p.GetRoomManager().GetRoom(bet.RoomID)
	if !ok {
		p.SendCommand(config.CodeRoomNotFound, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor room not found, roomID=%d, user:%s", bet.RoomID, user.Name()))
		return errors.New("get room error")
	}

	//check room active
	if room.Active() != 1 {

		p.SendCommand(config.CodeRoomInactive, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor room inactive, roomID=%d, user:%s", room.ID(), user.Name()))
		return errors.New("room inactive")
	}

	//check room type
	rTyp := room.Type()
	if rTyp != conf.RoomType0() {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor received data error,user:%s", user.Name()))
		return err
	}

	roomBet := control.CmdDataToRoomDataType0(&bet)

	//撿查投注是否合法
	if !control.CheckBetType0(roomBet) {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor CheckBetType0 user:%s, err", user.Name()))
		return errors.New("check user bet error")
	}
	//取之前投注
	oriBet, err := p.roomCtrl.GetUserBetType0(room, user.UserID())

	//之前投注與現在投注相加
	//addBet,_:=ctrlBet.AddBetUpType0(roomBet,oriBet)

	//撿查投注是否合限紅

	//本次投注總額
	sum, err := control.CountBetSumType0(roomBet)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor CountBetSumType0 user:%s, err:%s", user.Name(), err.Error()))
		return errors.New("count user bet sum error")
	}
	//之前投注總額

	oriSum, err := control.CountBetSumType0(oriBet)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor CountBetSumType0 ori user:%s, err:%s", user.Name(), err.Error()))
		return errors.New("count user ori bet sum error")
	}
	//檢查總額
	if int(user.Credit()) < (sum + oriSum) {
		p.SendCommand(config.CodeCreditNotEnough, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor user credit not enough error,user:%s", user.Name()))
		return errors.New("user credit not enough")
	}

	//room get  user old bet data and add new
	p.roomCtrl.AddBet(room, user.UserID(), roomBet)

	resData := []config.BetType0ResData{
		{
			UserID:      user.UserID(),
			RoomID:      bet.RoomID,
			Banker:      bet.Banker,
			Player:      bet.Player,
			Tie:         bet.Tie,
			BankerPair:  bet.BankerPair,
			PlayerPair:  bet.PlayerPair,
			Big:         bet.Big,
			Small:       bet.Small,
			AnyPair:     bet.AnyPair,
			PerfectPair: bet.PerfectPair,
			SuperSix:    bet.SuperSix,
		},
	}

	b, err := json.Marshal(resData)
	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("betType0Processor json marshal res data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	sendData := base64.StdEncoding.EncodeToString(b)
	p.SendCommand(config.CodeSuccess, 0, conf.CmdBetType0(), sendData, user, p.GetRoomReceivers(room.ID()))

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("betType0Processor complete  user id=%d,user=%s, resData=%+v ", user.UserID(), user.Name(), resData))

	return nil
}

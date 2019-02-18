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
	logPrefix:="betType0"
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s user==nil ",logPrefix))
		return fmt.Errorf("%s user==nil",logPrefix)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix,r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s base64 decode cmd data error,user:%s,error:%s", logPrefix,user.Name(), err.Error()))
		return err
	}

	//取 data
	data := []config.BetType0CmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json unmarshal cmd data error,user:%s,error:%s",logPrefix, user.Name(), err.Error()))
		return err
	} else if len(data) < 0 {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s received data error,user:%s", logPrefix,user.Name()))
		return err
	}

	bet := data[0]
	room, ok := p.GetRoomManager().GetRoom(bet.RoomID)
	if !ok {
		p.SendCommand(config.CodeRoomNotFound, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s room not found, roomID=%d, user:%s", logPrefix,bet.RoomID, user.Name()))
		return fmt.Errorf("%s get room error",logPrefix)
	}

	//check room active
	if room.Active() != 1 {

		p.SendCommand(config.CodeRoomInactive, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s room inactive, roomID=%d, user:%s", logPrefix,room.ID(), user.Name()))
		return fmt.Errorf("%s room inactive",logPrefix)
	}

	//check room type
	rTyp := room.Type()
	if rTyp != conf.RoomType0() {
		p.SendCommand(config.CodeReceivedDataError, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s received data error,user:%s", logPrefix,user.Name()))
		return err
	}

	roomBet := control.CmdDataToRoomDataType0(&bet)

	//撿查投注是否合法
	if !control.CheckBetType0(roomBet) {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s CheckBetType0 user:%s, err", logPrefix,user.Name()))
		return fmt.Errorf("%s check user bet error",logPrefix)
	}
	//取之前投注
	oriBet, err := p.roomCtrl.GetUserBetType0(room, user.UserID())

	//之前投注與現在投注相加
	//addBet,_:=ctrlBet.AddBetUpType0(roomBet,oriBet)

	//撿查投注是否合限紅

	//本次投注總額
	sum, err := control.CountBetSumType0(roomBet)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s CountBetSumType0 user:%s, err:%s", logPrefix,user.Name(), err.Error()))
		return fmt.Errorf("%s count user bet sum error",logPrefix)
	}
	//之前投注總額

	oriSum, err := control.CountBetSumType0(oriBet)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s CountBetSumType0 ori user:%s, err:%s", logPrefix,user.Name(), err.Error()))
		return fmt.Errorf("%s count user ori bet sum error",logPrefix)
	}
	//檢查總額
	if int(user.Credit()) < (sum + oriSum) {
		p.SendCommand(config.CodeCreditNotEnough, 0, conf.CmdBetType0(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s user credit not enough error,user:%s", logPrefix,user.Name()))
		return fmt.Errorf("%s user credit not enough",logPrefix)
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
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json marshal res data error,user:%s,error:%s", logPrefix,user.Name(), err.Error()))
		return err
	}

	sendData := base64.StdEncoding.EncodeToString(b)
	p.SendCommand(config.CodeSuccess, 0, conf.CmdBetType0(), sendData, user, p.GetRoomReceivers(room.ID()))

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete  user id=%d,user=%s, resData=%+v ", logPrefix,user.UserID(), user.Name(), resData))

	return nil
}

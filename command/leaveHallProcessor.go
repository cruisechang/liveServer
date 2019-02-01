package command

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
)

type leaveHallProcessor struct {
	BasicProcessor
}

func NewLeaveHallProcessor(processor BasicProcessor) (*leaveHallProcessor, error) {
	p := &leaveHallProcessor{
		BasicProcessor: processor,
	}

	return p, nil

}

func (p *leaveHallProcessor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("leaveHallProcessor user==nil "))
		return errors.New("leaveHallProcessor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("leaveHallProcessor panic:%v", r))
		}
	}()

	//parsing command data
	_, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdLeaveHall(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("leaveHallProcessor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	h, ok := p.GetHallManager().GetHall(user.HallID())
	if !ok {
		p.SendCommand(config.CodeHallNotFound, 0, conf.CmdLeaveHall(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("leaveHallProcessor hall not found  user:%s", user.Name()))
		return errors.New("hall not found")
	}

	h.RemoveUser(user)
	//empty resData
	//no need

	p.SendCommand(config.CodeSuccess, 0, conf.CmdLeaveHall(), p.DefaultSendData(), user, []string{user.ConnID()})
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("leaveHallProcessor complete  user id=%d,user=%s 2", user.UserID(), user.Name()))

	return nil
}

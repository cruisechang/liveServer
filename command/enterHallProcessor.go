package command

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
	"github.com/juju/errors"
)

type enterHallProcessor struct {
	BasicProcessor
}

func NewEnterHallProcessor(processor BasicProcessor) (*enterHallProcessor, error) {
	p := &enterHallProcessor{
		BasicProcessor: processor,
	}

	return p, nil

}

func (p *enterHallProcessor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterHallProcessor user==nil "))
		return errors.New("enterHallProcessor user ==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("enterHallProcessor panic:%v", r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdEnterHall(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterHallProcessor base64 decode cmd data error,user:%s,error:%s", obj.User.Name(), err.Error()))
		return err
	}

	data := []config.EnterHallCmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdEnterHall(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterHallProcessor json unmarshal cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	//check if targetHallID valid,
	h, ok := p.GetHallManager().GetHall(data[0].HallID)
	if !ok {
		p.SendCommand(config.CodeHallNotFound, 0, conf.CmdEnterHall(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterHallProcessor get hall  error,user:%s", user.Name()))
		return err
	}

	h.AddUser(user)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdEnterHall(), p.DefaultSendData(), user, []string{user.ConnID()})
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("enterHallProcessor complete  user id=%d,user=%s", user.UserID(), user.Name()))

	return nil
}

package command

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
)

type ClientReadyProcessor interface {
	Run(obj *nex.CommandObject) error
}

//loginProcessor implements command.Processor
type clientReadyProcessor struct {
	BasicProcessor
}

func NewClientReadyProcessor(processor BasicProcessor) (ClientReadyProcessor, error) {
	c := &clientReadyProcessor{
		BasicProcessor: processor,
	}
	return c, nil
}

//return error means some thing strange happend
//回傳error是表示程式發生不預期錯誤，
//如果是login 失敗，不回傳錯誤
//faild must logout
//failed must  disconnect client
func (p *clientReadyProcessor) Run(obj *nex.CommandObject) error {

	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("clientReadyProcessor user==nil "))
		return errors.New("clientReadyProcessor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("clientReadyProcessor panic:%v", r))
		}
	}()

	//parsing cmd data
	de, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("clientReadyProcessor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	data := []config.ClientReadyCmdData{}

	if err := json.Unmarshal(de, &data); err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("clientReadyProcessor json Unmarshal cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	//user
	user.SetIntVariable(conf.UserVarClientReady(), conf.ClientReady())

	p.SendCommand(config.CodeSuccess, 0, conf.CmdClientReady(), p.DefaultSendData(), user, []string{user.ConnID()})

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("clientReadyProcessor complete user=%s", user.Name()))

	return nil
}

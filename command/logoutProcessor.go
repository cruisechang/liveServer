package command

import (
	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
	"fmt"
	"encoding/base64"
	"encoding/json"
	"github.com/cruisechang/liveServer/config"
	"errors"
)

//logoutProcessor implements command.Processor
type logoutProcessor struct {
	BasicProcessor
}

func NewLogoutProcessor(processor BasicProcessor) (*logoutProcessor, error) {
	c := &logoutProcessor{
		BasicProcessor: processor,
	}
	return c, nil
}

func (p *logoutProcessor) Run(obj *nex.CommandObject) error {
	logPrefix := "logoutProcessor"
	logger := p.GetLogger()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s user==nil ", logPrefix))
		return errors.New("logoutProcessor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	//錯誤照樣斷線
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s base64 decode cmd data error,user:%s,error:%s", logPrefix, user.Name(), err.Error()))
	}

	data := []config.LogoutCmdData{}

	//錯誤照樣斷線
	if err := json.Unmarshal(deStr, &data); err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json unmarshal cmd data error,user:%s,error:%s", logPrefix, user.Name(), err.Error()))
	}

	//disconnect
	//不logout,等UserLostEvent再logout
	p.DisconnectUser(user.UserID())

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete,user:%s", logPrefix, user.Name()))
	return nil
}

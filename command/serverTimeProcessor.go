package command

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
)

type serverTimeProcessor struct {
	BasicProcessor
}

func NewServerTimeProcessor(processor BasicProcessor) (*serverTimeProcessor, error) {
	p := &serverTimeProcessor{
		BasicProcessor: processor,
	}

	return p, nil

}

func (p *serverTimeProcessor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()

	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("serverTimeProcessor Code Processor user==nil"))
		return errors.New("serverTimeProcessor Code Processor user==nil ")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("serverTimeProcessor panic:%v", r))
		}
	}()

	//parsing command data
	//deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)
	//
	//if err != nil {
	//	p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdServerTime(), p.DefaultSendData(), user, []string{user.ConnID()})
	//	logger.LogFile(nxLog.LevelError, fmt.Sprintf("serverTimeProcessor base64 decode cmd data error,user:%s,error:%s\n", user.Name(), err.Error()))
	//	return err
	//}
	//
	//data := []config.ServerTimeResData{
	//	config.ServerTimeResData{
	//		ServerTime:int(time.Now().UnixNano()/100000),
	//	},
	//}
	//
	//if err := json.Unmarshal(deStr, &data); err != nil {
	//	p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdServerTime(), p.DefaultSendData(), user, []string{user.ConnID()})
	//	logger.LogFile(nxLog.LevelError, fmt.Sprintf("serverTimeProcessor json unmarshal cmd data error,user:%s,error:%s\n", user.Name(), err.Error()))
	//	return err
	//}

	resData := []config.ServerTimeResData{
		config.ServerTimeResData{
			ServerTime: int(time.Now().UnixNano() / 100000),
		},
	}

	b, err := json.Marshal(resData)
	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("serverTimeProcessor json marshal res data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	sendData := base64.StdEncoding.EncodeToString(b)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdServerTime(), sendData, user, []string{user.ConnID()})

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("serverTimeProcessor complete  user id=%d,user=%s", user.UserID(), user.Name()))
	return nil
}

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

type enterRoomProcessor struct {
	BasicProcessor
}

func NewEnterRoomProcessor(processor BasicProcessor) (*enterRoomProcessor, error) {
	p := &enterRoomProcessor{
		BasicProcessor: processor,
	}

	return p, nil

}

func (p *enterRoomProcessor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()

	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterRoomProcessor user==nil "))
		return errors.New("enterRoomProcessor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("enterRoomProcessor panic:%v", r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdEnterRoom(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterRoomProcessor base64 decode cmd data error,user:%s,error:%s", obj.User.Name(), err.Error()))
		return err
	}

	data := []config.EnterRoomCmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdEnterRoom(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterRoomProcessor json unmarshal cmd data error,user:%s,error:%s", obj.User.Name(), err.Error()))
		return err
	}

	room, ok := p.GetRoomManager().GetRoom(data[0].RoomID)
	if !ok {
		p.SendCommand(config.CodeRoomNotFound, 0, conf.CmdEnterRoom(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("enterRoomProcessor get room not found ,user:%s", user.Name()))
		return errors.New("room not found")
	}
	room.AddUser(user)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdEnterRoom(), p.DefaultSendData(), user, []string{user.ConnID()})

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("enterRoomProcessor complete  user id=%d,user=%s", user.UserID(), user.Name()))

	return nil
}

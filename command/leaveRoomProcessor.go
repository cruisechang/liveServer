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

type leaveRoomProcessor struct {
	BasicProcessor
}

func NewLeaveRoomProcessor(processor BasicProcessor) (*leaveRoomProcessor, error) {
	p := &leaveRoomProcessor{
		BasicProcessor: processor,
	}

	return p, nil

}

func (p *leaveRoomProcessor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("leaveRoomProcessor user==nil "))
		return errors.New("leaveRoomProcessor user==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("leaveRoomProcessor panic:%v", r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdLeaveRoom(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("leaveRoomProcessor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	data := []config.LeaveRoomCmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdLeaveRoom(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("leaveRoomProcessor json unmarshal cmd data error,user:%s,error:%s", obj.User.Name(), err.Error()))
		return err
	}

	room, ok := p.GetRoomManager().GetRoom(user.RoomID())
	if !ok {
		p.SendCommand(config.CodeRoomNotFound, 0, conf.CmdLeaveRoom(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("leaveRoomProcessor room not found  user:%s", user.Name()))
		return errors.New("room not found")
	}

	room.RemoveUser(user)
	p.SendCommand(config.CodeSuccess, 0, conf.CmdLeaveRoom(), p.DefaultSendData(), user, []string{user.ConnID()})
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("leaveRoomProcessor complete  user id=%d,user=%s", user.UserID(), user.Name()))

	return nil
}

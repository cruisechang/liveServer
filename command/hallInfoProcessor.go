package command

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
	nxLog "github.com/cruisechang/nex/log"
)

type hallInfoProcessor struct {
	BasicProcessor
}

func NewHallInfoProcessor(processor BasicProcessor) (*hallInfoProcessor, error) {
	p := &hallInfoProcessor{
		BasicProcessor: processor,
	}

	return p, nil

}

func (p *hallInfoProcessor) Run(obj *nex.CommandObject) error {
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("hallInfoProcessor user==nil "))
		return errors.New("hallInfoProcessor user ==nil")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("hallInfoProcessor panic:%v", r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdHallInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("HallInfoProcessor base64 decode cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	//取hall data
	data := []config.HallInfoCmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdHallInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("HallInfoProcessor json unmarshal cmd data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	hm := p.GetHallManager()

	//get hall id
	hallID := data[0].HallID
	var halls []entity.Hall

	//get hall
	if hallID != conf.UnAssigned() {
		h, _ := hm.GetHall(hallID)
		halls = append(halls, h)
	} else {
		hs := hm.GetHalls()

		for _, h := range hs {
			if h, ok := h.(entity.Hall); ok {
				halls = append(halls, h)
			}
		}
	}

	//res data
	resData := []config.HallInfoResData{}

	for _, h := range halls {
		if h, ok := h.(entity.Hall); ok {

			re := config.HallInfoResData{
				HallID: h.ID(),
				Name:   h.Name(),
				Active: h.Active(),
			}
			resData = append(resData, re)
		}
	}

	b, err := json.Marshal(resData)
	if err != nil {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("hallInfoProcessor json marshal res data error,user:%s,error:%s", user.Name(), err.Error()))
		return err
	}

	sendDataStr := base64.StdEncoding.EncodeToString(b)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdHallInfo(), sendDataStr, user, []string{user.ConnID()})
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("hallInfoProcessor complete  user id=%d,user=%s", user.UserID(), user.Name()))

	return nil
}

package builtinevent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cruisechang/liveServer/command"
	"strconv"

	"github.com/cruisechang/nex/builtinEvent"
	nxLog "github.com/cruisechang/nex/log"
)

//NewUserLostEventProcessor returns user lost event structure
func NewUserLostEventProcessor(processor command.BasicProcessor) (*UserLostEventProcessor, error) {
	p := &UserLostEventProcessor{
		BasicProcessor: processor,
	}

	return p, nil
}

//UserLostEventProcessor removes lost user and patch user login
type UserLostEventProcessor struct {
	command.BasicProcessor
}

//Run executes procedure
func (p *UserLostEventProcessor) Run(obj *builtinEvent.EventObject) error {
	logPrefix := "userLostEventProcessor"
	logger := p.GetLogger()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s user==nil", logPrefix))
		return fmt.Errorf("%s user==nil", logPrefix)
	}

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	DBUserID, _ := user.GetInt64Variable(p.GetConfigurer().UserVarDBUserID())

	//remove user
	p.RemoveUser(user.UserID())

	//patch login
	uid := strconv.FormatInt(DBUserID, 10)
	path := fmt.Sprintf("/users/" + uid + "/login")

	bo, err := json.Marshal(struct{ Login int }{0})
	if err != nil {
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s json marshal body erro=%s", logPrefix, err.Error()))
		return fmt.Errorf("%s json marshal body erro=%s", logPrefix, err.Error())
	}
	_, err = p.DBAPIDo("PATCH", path, bytes.NewBuffer(bo))
	if err != nil {
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s db patch login err=%s", logPrefix, err.Error()))
		return fmt.Errorf("%s db patch login err=%s", logPrefix, err.Error())
	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete  user id=%d,user=%s", logPrefix, user.UserID(), user.Name()))

	return nil
}

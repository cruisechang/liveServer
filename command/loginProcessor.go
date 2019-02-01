package command

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/dbConf"
	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
)

//loginProcessor implements command.Processor
type loginProcessor struct {
	BasicProcessor
}

func NewLoginProcessor(processor BasicProcessor) (*loginProcessor, error) {
	c := &loginProcessor{
		BasicProcessor: processor,
	}
	return c, nil
}

//return error means some thing strange happend
//回傳error是表示程式發生不預期錯誤，
//如果是login 失敗，不回傳錯誤
//failed must  disconnect client, 之後會 logout
func (p *loginProcessor) Run(obj *nex.CommandObject) error {
	logPrefix := "loginProcessor"

	conf := p.GetConfigurer()
	logger := p.GetLogger()
	user := obj.User

	code := config.CodeSuccess

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
			code = config.CodePanic
		}
	}()

	if user == nil {
		code = config.CodeUserNil
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s user==nil", logPrefix))
		return errors.New("loginProcessor user==nil")
	}

	//response data
	errResData := []config.LoginCmdResData{{
		UserID: user.UserID(),
	}}
	r, _ := json.Marshal(errResData)

	sendDataStr := base64.StdEncoding.EncodeToString(r)

	//parsing cmd data
	de, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s base64 decode cmd data error=%s", logPrefix, err.Error()))
		return err
	}

	data := []config.LoginCmdData{}

	if err := json.Unmarshal(de, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json Unmarshal cmd data error=%s", logPrefix, err.Error()))
		return err
	}

	accessToken := data[0].SessionID
	//這邊要直接擋掉不合法的


	//get user data from db
	resp, err := p.DBAPIDo("GET", "/users/"+accessToken+"/tokenData", nil)

	if err != nil {
		p.SendCommand(config.CodeDBError, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db response error:%s, token=%s", logPrefix, err.Error(), accessToken))
		return err
	}

	if resp.StatusCode != http.StatusOK {
		p.SendCommand(config.CodeDBHTTPStatusCodeNotOK, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db http status code wrong=%d ,error:%s,token=%s", logPrefix, resp.StatusCode, err.Error(), accessToken))
		return err
	}

	if resp.Body == nil {
		p.SendCommand(config.CodeDBBodyNil, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db response body==nil,error:%s,token=%s", logPrefix, err.Error(), accessToken))
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	//db got
	got := &dbConf.AccessTokenData{}
	err = json.Unmarshal(body, got)
	if err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s unmarshal db got body error:%s", logPrefix, err.Error()))
		return errors.New(fmt.Sprintf("%s unmarshal db got body error=%s", logPrefix, err.Error()))
	}

	if got.Code != config.CodeSuccess {
		p.SendCommand(config.CodeDBQueryFail, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db got code=%d,token=%s", logPrefix, got.Code, accessToken))
		return errors.New(fmt.Sprintf("%s db got code code=%d ,token=%s", logPrefix, got.Code, accessToken))
	}
	if got.Count != 1 {
		p.SendCommand(config.CodeLoginAccessTokenError, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db got count!=1 count=%d,token=%s", logPrefix, got.Count, accessToken))
		return errors.New(fmt.Sprintf("%s db got count!=1, token=%s", logPrefix, accessToken))
	}

	//got data

	if len(got.Data) != 1 {
		p.SendCommand(config.CodeLoginAccessTokenError, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s db got data len!=1 data=%v,token=%s", logPrefix, got.Data, accessToken))
		return errors.New(fmt.Sprintf("%s db got data len !=1, token=%s", logPrefix, accessToken))
	}

	gotData := got.Data[0]
	if gotData.Active != 1 {
		p.SendCommand(config.CodeLoginUseDeactive, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s gor db data active !=1 ,token=%s", logPrefix, accessToken))
		return errors.New(fmt.Sprintf("%s db user deactive, accessToken=%s", logPrefix, accessToken))
	}

	//DB API patch login
	DBUserID := strconv.FormatInt(gotData.UserID, 10)
	path := fmt.Sprintf("/users/" + DBUserID + "/login")

	bo, err := json.Marshal(struct{ Login int }{1})
	if err != nil {
		p.SendCommand(config.CodeLoginUseDeactive, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s json marshal body erro=%s", logPrefix, err.Error()))
		return errors.New(fmt.Sprintf("%s json marshal body erro=%s", logPrefix, err.Error()))
	}
	_, err = p.DBAPIDo("PATCH", path, bytes.NewBuffer(bo))
	if err != nil {
		p.SendCommand(config.CodeLoginUseDeactive, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s db patch login err=%s", logPrefix, err.Error()))
		return errors.New(fmt.Sprintf("%s db patch login err=%s", logPrefix, err.Error()))
	}

	//response data
	resData := []config.LoginCmdResData{{
		UserID: user.UserID(),
	}}

	b, err := json.Marshal(resData)
	if err != nil {
		p.SendCommand(config.CodeMarshalJsonFailed, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s json marshal res data error=%s,token=%s", logPrefix, err.Error(), accessToken))
		return err
	}

	//user
	//userID

	user.SetSessionID(accessToken)
	user.SetAccount(gotData.Account)
	user.SetName(gotData.Name)
	//TODO: 暫時改成給錢
	if gotData.Credit<10000{
		gotData.Credit=10000
	}
	//TODO: 暫時改成給錢

	user.SetCredit(gotData.Credit)
	user.SetInt64Variable(conf.UserVarPartnerID(), gotData.PartnerID)
	user.SetInt64Variable(conf.UserVarDBUserID(), gotData.UserID) //userID from db, user.UserID()  是 nex的userID，不是db的userID
	sendDataStr = base64.StdEncoding.EncodeToString(b)

	//check user is in game
	//已經登入把他斷線

	p.SendCommand(code, 0, conf.CmdLogin(), sendDataStr, user, []string{user.ConnID()})

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s gotData=%+v ", logPrefix, gotData))
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete user account=%s, name=%s", logPrefix, gotData.Account, gotData.Name))

	return nil
}

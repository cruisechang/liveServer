package command

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/dbConf"

	"github.com/cruisechang/nex"
	nxLog "github.com/cruisechang/nex/log"
)

//NewUserInfoProcessor creates struct of user info
func NewUserInfoProcessor(processor BasicProcessor) (*UserInfoProcessor, error) {
	p := &UserInfoProcessor{
		BasicProcessor: processor,
	}

	return p, nil
}

//UserInfoProcessor presents struct of user info
type UserInfoProcessor struct {
	BasicProcessor
}

//Run handles process after receive "userInfo" command
func (p *UserInfoProcessor) Run(obj *nex.CommandObject) error {
	logPrefix := "userInfoProcessor"
	logger := p.GetLogger()
	conf := p.GetConfigurer()
	user := obj.User

	if user == nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s user==nil", logPrefix))
		return errors.New("%s user==nil")
	}

	userDBID, _ := user.GetInt64Variable(conf.UserVarDBUserID())

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v", logPrefix, r))
		}
	}()

	//parsing command data
	deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)

	if err != nil {
		p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdUserInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s base64 decode cmd data error,userID=%d,error:%s", logPrefix, userDBID, err.Error()))
		return err
	}

	data := []config.UserInfoCmdData{}

	if err := json.Unmarshal(deStr, &data); err != nil {
		p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdUserInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s json umarshal cmd data error,userID=%d, error:%s", logPrefix, userDBID, err.Error()))
		return err
	}

	if data[0].UserID < 0 {
		p.DisconnectUser(user.UserID())
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s userID < 0 ", logPrefix))
		return err
	}

	u, ok := p.GetUser(data[0].UserID)
	if !ok {
		p.SendCommand(config.CodeUserNil, 0, conf.CmdUserInfo(), p.DefaultSendData(), user, []string{user.ConnID()})
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s target user not found ID=%d ", logPrefix, data[0].UserID))
		return fmt.Errorf("%s target user not found, ID=%d ", logPrefix, data[0].UserID)
	}

	//assign data
	resData := []config.UserInfoResData{{}}

	resData[0].UserID = data[0].UserID
	resData[0].Credit = u.Credit()
	resData[0].Account = u.Account()

	resB, _ := json.Marshal(resData)
	sendDataStr := base64.StdEncoding.EncodeToString(resB)

	p.SendCommand(config.CodeSuccess, 0, conf.CmdUserInfo(), sendDataStr, user, []string{user.ConnID()})

	//broadcast
	broadcastDB, err := p.queryBroadcasts("GET", dbConf.DBPathBroadcasts)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s broadcast error=%s, userDBID=%d", logPrefix, err.Error(), userDBID))
	} else {
		broadcastData := p.createBroadcasts(conf, broadcastDB)
		if len(broadcastData) > 1 {

			resBC, _ := json.Marshal(broadcastData)
			sendDataStr := base64.StdEncoding.EncodeToString(resBC)
			p.SendCommand(config.CodeSuccess, 0, conf.CmdBroadcast(), sendDataStr, user, []string{user.ConnID()})
			logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s broadcast userDBID=%d, resData=%s", logPrefix, userDBID, string(resBC)))
		}
	}

	//banner
	bannerDB, err := p.queryBanners("GET", dbConf.DBPathBanners)
	if err != nil {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s banner error=%s, userDBID=%d,user=%s", logPrefix, err.Error(), userDBID, user.Name()))
	} else {
		bannerData := p.createBanners(conf, bannerDB)
		if len(bannerData) > 1 {
			resBD, _ := json.Marshal(bannerData)
			sendDataStr := base64.StdEncoding.EncodeToString(resBD)
			p.SendCommand(config.CodeSuccess, 0, conf.CmdBanner(), sendDataStr, user, []string{user.ConnID()})
			logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s banner userDBID=%d, resData=%s", logPrefix, userDBID, string(resBD)))
		}
	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s complete, userDBID=%d", logPrefix, userDBID))

	return nil
}

//broadcast
func (p *UserInfoProcessor) queryBroadcasts(method, path string) ([]*dbConf.Broadcast, error) {

	resp, err := p.DBAPIDo(method, path, nil)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("http status code wrong %d", resp.StatusCode)
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("response body==nil")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	resData := &dbConf.BroadcastData{}
	err = json.Unmarshal(body, resData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal responseData error=%s", err.Error())
	}
	return resData.Data, nil
}

func (p *UserInfoProcessor) createBroadcasts(conf config.Configurer, broadcasts []*dbConf.Broadcast) []*config.Broadcast {

	var sl []*config.Broadcast

	for _, v := range broadcasts {
		if v.Active == 1 {
			sl = append(sl, &config.Broadcast{
				Content:     v.Content,
				Internal:    v.Internal,
				RepeatTimes: v.RepeatTimes,
			})
		}
	}

	return sl
}

//banner
func (p *UserInfoProcessor) queryBanners(method, path string) ([]*dbConf.Banner, error) {

	resp, err := p.DBAPIDo(method, path, nil)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code wrong %d", resp.StatusCode)
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("response body==nil")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	resData := &dbConf.BannerData{}
	err = json.Unmarshal(body, resData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal responseData error=%s", err.Error())
	}

	return resData.Data, nil
}

func (p *UserInfoProcessor) createBanners(conf config.Configurer, banners []*dbConf.Banner) []*config.Banner {

	var sl []*config.Banner

	for _, v := range banners {
		if v.Active == 1 {
			sl = append(sl, &config.Banner{
				PicURL:   v.PicURL,
				LinkURL:  v.LinkURL,
				Platform: v.Platform,
			})
		}
	}

	return sl
}

package command

import (
	"github.com/cruisechang/nex"

)


type heartbeatProcessor struct {
	BasicProcessor
}

func NewHeartbeatProcessor(processor BasicProcessor) (*heartbeatProcessor, error) {
	p:=&heartbeatProcessor{
		BasicProcessor:processor,
	}



	return p,nil

}

func (p *heartbeatProcessor) Run(obj *nex.CommandObject) error {
	//logger:=p.GetLogger()
	//conf:=p.GetConfigurer()
	//
	//user := obj.User
	//
	//if user==nil{
	//	logger.LogFile(nxLog.LevelError, fmt.Sprintf("hallInfoProcessor CodeProcessorUserNil \n"))
	//	return errors.New("heartbeatProcessor CodeProcessorUserNil")
	//}
	//
	//defer func() {
	//	if r := recover(); r != nil {
	//		logger.LogFile(nxLog.LevelPanic,fmt.Sprintf("HeartbeatProcessor panic:%v\n", r))
	//	}
	//}()

	//parsing command data
	//deStr, err := base64.StdEncoding.DecodeString(obj.Cmd.Data)
	//
	//if err != nil {
	//	p.SendCommand(config.CodeBase64DecodeFailed, 0, conf.CmdHeartbeat(), p.DefaultSendData(), user, []string{user.ConnID()})
	//	logger.LogFile(nxLog.LevelError, fmt.Sprintf("HeartbeatProcessor base64 decode cmd data error,user:%s,error:%s\n", user.Name(), err.Error()))
	//	return err
	//}
	//
	//data := []config.HeartbeatCmdData{}
	//
	//if err := json.Unmarshal(deStr, &data); err != nil {
	//	p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdHeartbeat(), p.DefaultSendData(), user, []string{user.ConnID()})
	//	logger.LogFile(nxLog.LevelError, fmt.Sprintf("HeartbeatProcessor json unmarshal cmd data error,user:%s,error:%s\n",user.Name(), err.Error()))
	//	return err
	//}

	//p.SendCommand(config.CodeSuccess, 0, conf.CmdHeartbeat(), p.DefaultSendData(), user, []string{user.ConnID()})


	return nil
}

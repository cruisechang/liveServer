package command

import (
	"encoding/base64"
	"encoding/json"
	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
	"testing"
)


type BasicProcessorTester interface {
	Run(t *testing.T) error
}

type basicProcessorTester struct {
	BasicProcessor
}

func NewBasicProcessorTester(processor BasicProcessor) (BasicProcessorTester, error) {
	p := &basicProcessorTester{
		BasicProcessor: processor,
	}

	return p, nil

}

func TestInheritBasicProcessor(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")
	dbCtrl := control.NewDBController(conf.DBAPIServer())
	rCtrl := roomCtrl.NewController(conf)
	rmc:=control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl,nx.GetLogger())
	bpt, err := NewBasicProcessorTester(NewBasicProcessor(nx, conf,dbCtrl,rmc))
	if err != nil {

	}

	bpt.Run(t)

}

func (p *basicProcessorTester) Run(t *testing.T) error {

conf:=p.GetConfigurer()

	//sendCommand
	resData := []config.LoginCmdResData{config.LoginCmdResData{}}

	b, _ := json.Marshal(resData)

	//[]byte encode to base64 string
	sendDataStr := base64.StdEncoding.EncodeToString(b)

	t.Logf("SendCommand %s\n", sendDataStr)
	p.SendCommand(config.CodeJsonUnmarshalJsonFailed, 0, conf.CmdLogin(), sendDataStr, entity.NewUser(0,"connid"), []string{"uuid"})

	return nil
}



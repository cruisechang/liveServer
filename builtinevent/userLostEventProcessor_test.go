package builtinevent

import (
	"github.com/cruisechang/liveServer/command"
	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/control"
	nexSpace "github.com/cruisechang/nex"
	nexBuiltEvent "github.com/cruisechang/nex/builtinEvent"
	"github.com/cruisechang/nex/entity"
	"testing"
)


func TestNewUserLostEventProcessor(t *testing.T) {

	nex, err := nexSpace.NewNex("nexConfig.json")
	if err != nil {
		t.Errorf("TestLogin error %s \n", err.Error())
		return
	}
	conf,_:=config.NewConfigurer("config.json")
	dbCtrl := control.NewDBController(conf.DBAPIServer())

	ule,_:=NewUserLostEventProcessor(command.NewBasicProcessor(nex,conf,dbCtrl))



	user:=entity.NewUser(0,"connID")
	user.SetAccessID("accessID")

	obj:=nexBuiltEvent.EventObject{
		Code:1,
		User:user,
	}

	t.Logf("%#v",obj)

	ule.Run(&obj)

}


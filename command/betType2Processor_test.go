package command

import (
	"testing"

	"github.com/cruisechang/liveServer/config"

	"encoding/base64"
	"encoding/json"

	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_betType2Processor_Run(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType2(), "name")

	//dbCtrl := control.NewDBController(conf.DBAPIServer())
	//rCtrl := roomCtrl.NewController(conf)
	//rmc:=control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl,nx.GetLogger())
	p, _ := NewBetType2Processor(getBasicProcessor())

	user := entity.NewUser(0, "conn")
	user.SetCredit(100000)

	obj := &[]config.BetType2CmdData{
		{
			RoomID: 0,
			Owner1: config.BetType2CmdDataOwner{WinEqual: 10, WinTimes: 10, LoseEqual: 10, LoseTimes: 10},
			Owner2: config.BetType2CmdDataOwner{WinEqual: 10, WinTimes: 10, LoseEqual: 10, LoseTimes: 10},
			Owner3: config.BetType2CmdDataOwner{WinEqual: 10, WinTimes: 10, LoseEqual: 10, LoseTimes: 10},
		},
	}
	c, err := json.Marshal(obj)
	if err != nil {
		t.Logf("%s", err.Error())
	}
	cs := base64.StdEncoding.EncodeToString(c)

	cmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType2(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType2(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType2(),
			Step:    0,
			Data:    cs,
		},
		User: nil,
	}

	type args struct {
		obj *nex.CommandObject
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{cmd},
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{errcmd},
			wantErr: true,
		},
		{
			name:    "2",
			args:    args{errcmd1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.Run(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("betType2Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

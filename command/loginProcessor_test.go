package command

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_loginProcessor_Run(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType0(), "name")

	//dbCtrl := control.NewDBController(conf.DBAPIServer())
	p, _ := NewLoginProcessor(getBasicProcessor())

	user := entity.NewUser(0, "conn")

	//
	d0 := &[]config.LoginCmdData{{
		SessionID: "769d3e12-4e78-4419-b9b3-77fb99f25e59",
	},
	}
	md0, _ := json.Marshal(d0)
	bd0 := base64.StdEncoding.EncodeToString(md0)

	//
	d1 := &[]config.LoginCmdData{{
		SessionID: "9999999",
	},
	}
	md1, _ := json.Marshal(d1)
	bd1 := base64.StdEncoding.EncodeToString(md1)

	cmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdLogin(),
			Step:    0,
			Data:    bd0,
		},
		User: user,
	}

	cmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdLogin(),
			Step:    0,
			Data:    bd1,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdLogin(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdLogin(),
			Step:    0,
			Data:    bd1,
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
			args:    args{cmd1},
			wantErr: true,
		},
		{
			name:    "2",
			args:    args{errcmd},
			wantErr: true,
		},
		{
			name:    "3",
			args:    args{errcmd1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.Run(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("loginProcessor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

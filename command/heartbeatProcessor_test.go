package command

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_heartbeatProcessor_Run(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType0(), "name")

	//dbCtrl := control.NewDBController(conf.DBAPIServer())
	p, _ := NewHeartbeatProcessor(getBasicProcessor())

	user := entity.NewUser(0, "conn")

	obj := &[]config.LoginCmdData{
		{
			SessionID: "sid",
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
			Command: conf.CmdHeartbeat(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHeartbeat(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHeartbeat(),
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
				t.Errorf("heartbeatProcessor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

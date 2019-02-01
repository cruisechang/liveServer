package command

import (
	"testing"

	"encoding/base64"
	"encoding/json"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_resultType7Processor_Run(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, err := config.NewConfigurer("config.json")

	if err != nil {
		t.Fatalf("err=%s", err.Error())
	}

	rCtrl := roomCtrl.NewController(conf)

	//create hall
	hall, _ := nx.GetHallManager().CreateHall(0, "hall")
	room, _ := nx.GetRoomManager().CreateRoom(0, conf.RoomType7(), "room")
	hall.AddRoom(room)

	resultHistory := roomConf.HistoryResultType7{1, 1, 0, 1, 0, 1, 1, 1, 5, 6, 7, 22, 32, 0, 15, 9, 8}
	rCtrl.SetHistoryResult(room, resultHistory)

	//dbCtrl := control.NewDBController(conf.DBAPIServer())
	p, _ := NewHistoryResultType7Processor(getBasicProcessorByRCTRL(nx,rCtrl))

	user := entity.NewUser(0, "conn")

	obj := &[]config.HistoryResultTypeCmdData{
		{
			HallID: 0,
			RoomID: -1,
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
			Command: conf.CmdHistoryResultType7(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}
	errCmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHistoryResultType7(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}
	errCmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHistoryResultType7(),
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
			args:    args{errCmd},
			wantErr: true,
		},
		{
			name:    "2",
			args:    args{errCmd1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := p.Run(tt.args.obj); (err != nil) != tt.wantErr {
				t.Errorf("resultType7Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

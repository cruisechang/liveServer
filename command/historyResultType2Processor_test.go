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

func Test_resultType2Processor_Run(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, err := config.NewConfigurer("config.json")

	if err != nil {
		t.Fatalf("err=%s", err.Error())
	}

	rCtrl := roomCtrl.NewController(conf)

	//create hall
	hall, _ := nx.GetHallManager().CreateHall(0, "hall")
	room, _ := nx.GetRoomManager().CreateRoom(100, conf.RoomType2(), "room")
	hall.AddRoom(room)

	resultHistory := roomConf.HistoryResultType2{
		{1, 1, 0, 1, 0, 1, 1, 1}, {1, 1, 0, 2, 1, 2, 1, 2}, {1, 1, 0, 3, 0, 3, 1, 3}, {1, 1, 1, 4, 0, 4, 1, 4}, {1, 1, 0, 5, 0, 5, 1, 5},
		{1, 1, 0, 6, 1, 6, 1, 6}, {1, 1, 1, 7, 1, 7, 1, 7}, {1, 1, 0, 8, 1, 8, 1, 8}, {1, 1, 1, 9, 0, 9, 1, 9}, {1, 1, 0, 10, 0, 10, 1, 10},
		{1, 1, 1, 11, 0, 11, 1, 11},
	}
	rCtrl.SetHistoryResult(room, resultHistory)

	//dbCtrl := control.NewDBController(conf.DBAPIServer())
	p, _ := NewHistoryResultType2Processor(getBasicProcessorByRCTRL(nx,rCtrl))

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
			Command: conf.CmdHistoryResultType2(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}
	errCmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHistoryResultType2(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}
	errCmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHistoryResultType2(),
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
				t.Errorf("resultType2Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

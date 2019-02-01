package command

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_resultType6Processor_Run(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	//create room
	rCtrl := roomCtrl.NewController(conf)

	//create hall
	hall, _ := nx.GetHallManager().CreateHall(0, "hall")
	room, _ := nx.GetRoomManager().CreateRoom(0, conf.RoomType6(), "room")
	hall.AddRoom(room)

	var resultHistory []*roomConf.HistoryResultType6
	resultHistory = append(resultHistory,
		&roomConf.HistoryResultType6{
			HallID:   0,
			RoomID:   1,
			Dice:     []int{1, 2, 2},
			Sum:      5,
			BigSmall: 0,
			OddEven:  1,
		})

	resultHistory = append(resultHistory,
		&roomConf.HistoryResultType6{
			HallID:   0,
			RoomID:   1,
			Dice:     []int{1, 2, 3},
			Sum:      6,
			BigSmall: 0,
			OddEven:  0,
		})

	resultHistory = append(resultHistory,
		&roomConf.HistoryResultType6{
			HallID:   0,
			RoomID:   1,
			Dice:     []int{4, 2, 3},
			Sum:      9,
			BigSmall: 1,
			OddEven:  1,
		})

	rCtrl.SetHistoryResult(room, resultHistory)

	//dbCtrl := control.NewDBController(conf.DBAPIServer())
	p, _ := NewHistoryResultType6Processor(getBasicProcessorByRCTRL(nx,rCtrl))

	user := entity.NewUser(0, "conn")

	obj := &[]config.HistoryResultType6ResData{
		{
			HallID:   0,
			RoomID:   0,
			Sum:      6,
			Dice:     []int{1, 2, 3},
			BigSmall: 0,
			OddEven:  0,
		},
		{
			HallID:   0,
			RoomID:   1,
			Sum:      6,
			Dice:     []int{1, 2, 3},
			BigSmall: 0,
			OddEven:  0,
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
			Command: conf.CmdHistoryResultType6(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHistoryResultType6(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdHistoryResultType6(),
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
				t.Errorf("resultType6Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

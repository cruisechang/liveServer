package command

import (
	"testing"

	"github.com/cruisechang/liveServer/config"

	"encoding/base64"
	"encoding/json"

	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_betType7Processor_Run(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	dbCtrl := control.NewDBController(conf.DBAPIServer())
	rCtrl := roomCtrl.NewController(conf)
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType7(), "name")

	p, _ := NewBetType0Processor(NewBasicProcessor(nx, conf, dbCtrl, rmc))

	user := entity.NewUser(0, "conn")
	user.SetCredit(100000)

	obj := &[]config.BetType7CmdData{
		{
			RoomID: 0,
			One:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Two:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Three:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Four:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Six:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Column: []int{0, 0, 0},
			Dozen:  []int{0, 0, 0},
			Big:    10,
			Small:  10,
			Odd:    10,
			Even:   10,
			Red:    10,
			Black:  10,
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
			Command: conf.CmdBetType7(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType7(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType7(),
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
				t.Errorf("betType7Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

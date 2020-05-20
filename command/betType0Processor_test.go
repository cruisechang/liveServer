package command

import (
	"github.com/cruisechang/liveServer/config"
	"testing"

	"encoding/base64"
	"encoding/json"
	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_betType0Processor_Run(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	dbCtrl := control.NewDBController(conf.DBAPIServer())
	rCtrl := roomCtrl.NewController(conf)
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType0(), "name")

	p, _ := NewBetType0Processor(NewBasicProcessor(nx, conf, dbCtrl, rmc))

	user := entity.NewUser(0, "conn")
	user.SetCredit(100000)

	obj := &[]config.BetType0CmdData{
		{
			RoomID:      0,
			Banker:      10,
			Player:      11,
			Tie:         12,
			BankerPair:  13,
			PlayerPair:  14,
			Big:         15,
			Small:       16,
			AnyPair:     17,
			PerfectPair: 18,
			SuperSix:    19,
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
			Command: conf.CmdBetType0(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType0(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType0(),
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
				t.Errorf("betType0Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_betType0Processor_countBetSum(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType0(), "name")

	dbCtrl := control.NewDBController(conf.DBAPIServer())
	rCtrl := roomCtrl.NewController(conf)
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())
	NewBetType0Processor(NewBasicProcessor(nx, conf, dbCtrl, rmc))

	bb := &config.BetType0CmdData{
		Banker:      100,
		Player:      100,
		Tie:         100,
		BankerPair:  100,
		PlayerPair:  100,
		Big:         100,
		Small:       100,
		AnyPair:     100,
		PerfectPair: 100,
		SuperSix:    100,
		Commission:  1,
	}

	type args struct {
		bet *config.BetType0CmdData
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "0",
			args: args{bb},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			roomBet := control.CmdDataToRoomDataType0(tt.args.bet)
			if got, _ := control.CountBetSumType0(roomBet); got != tt.want {
				t.Errorf("betType0Processor.countBetSum() = %v, want %v", got, tt.want)
			}
		})
	}
}

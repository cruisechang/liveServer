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

func Test_betType6Processor_Run(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	dbCtrl := control.NewDBController(conf.DBAPIServer())
	rCtrl := roomCtrl.NewController(conf)
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType6(), "name")

	p, _ := NewBetType0Processor(NewBasicProcessor(nx, conf, dbCtrl, rmc))

	user := entity.NewUser(0, "conn")
	user.SetCredit(100000)

	obj := &[]config.BetType6CmdData{
		{
			RoomID: 0,
			Small:  0,
			Big:    5,
			Odd:    10,
			Even:   10,
			Sum:    []int{0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Dice:   []int{5, 0, 0, 0, 0, 0},
			Triple: []int{10, 0, 0, 0, 0, 0, 0},
			Pair:   []int{10, 0, 0, 0, 0, 0},
			Paigow: []int{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
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
			Command: conf.CmdBetType6(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType6(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType6(),
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
				t.Errorf("betType6Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func Test_betType6Processor_CountBetSum(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType0(), "name")

	dbCtrl := control.NewDBController(conf.DBAPIServer())
	rCtrl := roomCtrl.NewController(conf)
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())
	NewBetType6Processor(NewBasicProcessor(nx, conf, dbCtrl, rmc))

	obj := &config.BetType6CmdData{
		RoomID: 0,
		Small:  0,
		Big:    5,
		Odd:    10,
		Even:   10,
		Sum:    []int{0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Dice:   []int{5, 0, 0, 0, 0, 0},
		Triple: []int{10, 0, 0, 0, 0, 0, 0},
		Pair:   []int{10, 0, 0, 0, 0, 0},
		Paigow: []int{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	obj1 := &config.BetType6CmdData{
		RoomID: 0,
		Small:  0,
		Big:    5,
		Odd:    10,
		Even:   10,
		Sum:    []int{0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		Dice:   []int{5, 0, 0, 0, 0, 0},
		Triple: []int{10, 0, 0, 0, 0, 0, 0},
		Pair:   []int{10, 0, 0, 0, 0, 0},
		Paigow: []int{10, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	obj2 := &config.BetType6CmdData{
		RoomID: 0,
		Small:  0,
		Big:    5,
		Odd:    10,
		Even:   10,
		Sum:    []int{0, 10, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		Dice:   []int{5, 0, 0, 0, 0, 0},
		Triple: []int{10, 0, 0, 0, 0, 0, 1},
		Pair:   []int{10, 0, 0, 0, 0, 1},
		Paigow: []int{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	type args struct {
		bet *config.BetType6CmdData
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "0",
			args: args{obj},
			want: 70,
		},
		{
			name: "1",
			args: args{obj1},
			want: 75,
		},
		{
			name: "2",
			args: args{obj2},
			want: 73,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roomBet := control.CmdDataToRoomDataType6(tt.args.bet)
			if got, _ := control.CountBetSumType6(roomBet); got != tt.want {
				t.Errorf("betType6Processor.CountBetSum() error, got = %d, want %d, name=%s", got, tt.want, tt.name)
			}
		})
	}
}

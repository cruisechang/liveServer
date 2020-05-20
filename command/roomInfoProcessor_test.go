package command

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_roomInfoProcessor_Run(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	dbCtrl := control.NewDBController(conf.DBAPIServer())
	rCtrl := roomCtrl.NewController(conf)
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())

	hallID := 0
	roomID := 1

	hm := nx.GetHallManager()
	hall, _ := hm.CreateHall(hallID, "hallName")
	rm := nx.GetRoomManager()
	room, _ := rm.CreateRoom(roomID, 0, "roomName")

	//create type data
	td := createTypeData()
	rCtrl.SetTypeData(room, td)

	hall.AddRoom(room)

	p, _ := NewRoomInfoProcessor(NewBasicProcessor(nx, conf, dbCtrl, rmc))

	user := entity.NewUser(0, "conn")

	obj := &[]config.RoomInfoCmdData{
		{
			HallID: hallID,
			RoomID: roomID,
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
			Command: conf.CmdRoomInfo(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdRoomInfo(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdRoomInfo(),
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
				t.Errorf("roomInfoProcessor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createTypeData() *roomConf.TypeData0 {
	return &roomConf.TypeData0{
		Boot:             1,
		Round:            1,
		BetLimit:         []int{10, 100},
		BankerLimit:      []int{10, 100},
		PlayerLimit:      []int{10, 100},
		TieLimit:         []int{10, 100},
		BankerPairLimit:  []int{10, 100},
		PlayerPairLimit:  []int{10, 100},
		AnyPairLimit:     []int{10, 100},
		PerfectPairLimit: []int{10, 100},
		SuperSixLimit:    []int{10, 100},
		BigSmallLimit:    []int{10, 100},
	}
}

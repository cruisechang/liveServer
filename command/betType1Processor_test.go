package command

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func Test_betType1Processor_Run(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")

	//create room
	nx.GetRoomManager().CreateRoom(0, conf.RoomType1(), "name")

	p, _ := NewBetType1Processor(getBasicProcessor())

	user := entity.NewUser(0, "conn")
	user.SetCredit(100000)

	obj := &[]config.BetType1CmdData{
		{
			RoomID:      0,
			Dragon:      10,
			Tiger:       10,
			Tie:         10,
			DragonOdd:   10,
			DragonEven:  10,
			DragonRed:   10,
			DragonBlack: 10,
			TigerOdd:    10,
			TigerEven:   10,
			TigerRed:    10,
			TigerBlack:  10,
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
			Command: conf.CmdBetType1(),
			Step:    0,
			Data:    cs,
		},
		User: user,
	}

	errcmd := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType1(),
			Step:    0,
			Data:    "",
		},
		User: user,
	}

	errcmd1 := &nex.CommandObject{
		Cmd: &nex.Command{
			Code:    0,
			Command: conf.CmdBetType1(),
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
				t.Errorf("betType1Processor.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

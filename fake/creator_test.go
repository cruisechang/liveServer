package fake

import (
	"github.com/cruisechang/liveServer/config"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"testing"

	nexSpace "github.com/cruisechang/nex"
)

var (
	testNex  nexSpace.Nex
	testConf config.Configurer
)

func TestPrepare(t *testing.T) {
	tn, err := nexSpace.NewNex(getConfigFilePosition("nexConfig.json"))
	if err != nil {
		t.Fatalf("TestPrepare error =%s", err.Error())
	}
	testNex = tn

	conf, err := config.NewConfigurer("config.json")

	if err != nil {
		t.Fatalf("TestPrepare error =%s", err.Error())
	}
	testConf = conf
}


func TestCreateRoom(t *testing.T) {

	opt0 := CreateRoomField(2001, "龍虎1", testConf.RoomType1())
	opt1 := CreateRoomField(2002, "龍虎2", testConf.RoomType1())

	type args struct {
		nex   nexSpace.Nex
		conf  *roomCtrl.Controller
		field RoomField
	}
	tests := []struct {
		name string
		args args
		wantErr bool
	}{
		{
			name: "0",
			args: args{
				nex:   testNex,
				conf:  roomCtrl.NewController(testConf),
				field: opt0,
			},
		},
		{
			name: "1",
			args: args{
				nex:   testNex,
				conf:  roomCtrl.NewController(testConf),
				field: opt1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := CreateRoom(tt.args.nex, tt.args.conf, tt.args.field)
			t.Logf("TestCreateRoom id=%d,name=%s\n", r.ID(), r.Name())
		})
	}
}

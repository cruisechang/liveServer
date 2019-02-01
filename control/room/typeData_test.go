package room

import (
	"testing"

	"github.com/cruisechang/liveServer/config"

	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
)

func TestController_GetBootRoundBetMinBetMax(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")
	r := NewController(conf)

	r0 := createRoom(nx, r, createRoomField(100, "百家1", conf.RoomType0()))
	r1 := createRoom(nx, r, createRoomField(200, "龍虎1", conf.RoomType1()))

	type args struct {
		room entity.Room
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   int64
		want2   int
		want3   int
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{r0},
			want:    1,
			want1:   1,
			want2:   1,
			want3:   10000,
			wantErr: false,
		},
		{
			name:    "0",
			args:    args{r1},
			want:    1,
			want1:   1,
			want2:   1,
			want3:   50000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			brData := r.GetBootRoundBetMinBetMax(tt.args.room)
			if (brData.Err != nil) != tt.wantErr {
				t.Errorf("Controller.GetBootRoundBetMinBetMax() error = %v, wantErr %v", brData.Err, tt.wantErr)
				return
			}
			if brData.Boot != tt.want {
				t.Errorf("Controller.GetBootRoundBetMinBetMax() got = %v, want %v", brData.Boot, tt.want)
			}
			if brData.Round != tt.want1 {
				t.Errorf("Controller.GetBootRoundBetMinBetMax() got1 = %v, want %v", brData.Round, tt.want1)
			}
			if brData.BetMin != tt.want2 {
				t.Errorf("Controller.GetBootRoundBetMinBetMax() got2 = %v, want %v", brData.BetMin, tt.want2)
			}
			if brData.BetMax != tt.want3 {
				t.Errorf("Controller.GetBootRoundBetMinBetMax() got3 = %v, want %v", brData.BetMax, tt.want3)
			}
		})
	}
}

func TestController_NeedInitBootRound(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")
	r := NewController(conf)

	r0 := createRoom(nx, r, createRoomField(200, "龍虎1", conf.RoomType1()))
	r1 := createRoom(nx, r, createRoomField(201, "龍虎1", conf.RoomType1()))

	r.SetRound(r1, 100000)

	type args struct {
		room entity.Room
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "0",
			args: args{r0},
			want: true,
		},
		{
			name: "1",
			args: args{r1},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.NeedInitBootRound(tt.args.room); got != tt.want {
				t.Errorf("Controller.NeedInitBootRound() = %v, want %v , name = %s", got, tt.want, tt.name)
				t.Errorf("Controller.NeedInitBootRound() = %v , name = %s", r.GetBootRoundBetMinBetMax(tt.args.room), tt.name)
			}
		})
	}
}

func TestController_InitBootRound(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")
	r := NewController(conf)

	r0 := createRoom(nx, r, createRoomField(200, "龍虎1", conf.RoomType1()))
	r.InitBootRound(r0)

	var initRound int64 = 20180815200000000

	type args struct {
		room entity.Room
	}
	tests := []struct {
		name  string
		args  args
		boot  int
		round int64
	}{
		{
			name:  "0",
			args:  args{r0},
			boot:  0,
			round: initRound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.InitBootRound(tt.args.room)

			brData := r.GetBootRoundBetMinBetMax(r0)
			if brData.Boot != tt.boot {
				t.Fatalf("TestController_InitBootRound error, boot got=%d, want=%d", brData.Boot, tt.boot)
			} else if brData.Round != tt.round {
				t.Fatalf("TestController_InitBootRound error, round got=%d, want=%d", brData.Round, tt.round)

			}
		})
	}
}

func TestController_HandleNewBoot(t *testing.T) {

	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")
	r := NewController(conf)

	r0 := createRoom(nx, r, createRoomField(200, "龍虎1", conf.RoomType1()))
	var resultBoot = 1
	var resultRound int64 = 20180815200000000

	type args struct {
		room entity.Room
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int64
	}{
		{
			name:  "0",
			args:  args{r0},
			want:  resultBoot,
			want1: resultRound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, got1 := r.HandleNewBoot(tt.args.room); got != tt.want {
				t.Errorf("Controller.NewBoot() boot = %v, want %v", got, tt.want)
			} else if got1 != tt.want1 {
				t.Fatalf("Controller.NewBoot() round = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestController_HandleNewRound(t *testing.T) {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")
	c := NewController(conf)

	r0 := createRoom(nx, c, createRoomField(200, "龍虎1", conf.RoomType1()))

	brData := c.GetBootRoundBetMinBetMax(r0)

	type args struct {
		room entity.Room
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int64
	}{
		{
			name:  "0",
			args:  args{r0},
			want:  brData.Boot,
			want1: brData.Round + 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if boot, round := c.HandleNewRound(tt.args.room); boot != tt.want {
				t.Fatalf("Controller.NewRound() boot = %v, want %v", boot, tt.want)
			} else if round != tt.want1 {
				t.Fatalf("Controller.NewRound() round = %v, want %v", boot, tt.want)
			}
		})
	}
}

package control

import (
	"github.com/cruisechang/liveServer/config/roomConf"
	"testing"
)

func TestCountBetSumType0(t *testing.T) {

	b0 := &roomConf.BetType0Data{
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
	b1 := &roomConf.BetType0Data{
		Banker:      100,
		Player:      0,
		Tie:         100,
		BankerPair:  100,
		PlayerPair:  100,
		Big:         100,
		Small:       0,
		AnyPair:     100,
		PerfectPair: 100,
		SuperSix:    100,
		Commission:  1,
	}
	b2 := &roomConf.BetType0Data{
		Banker:      50,
		Player:      0,
		Tie:         100,
		BankerPair:  100,
		PlayerPair:  100,
		Big:         50,
		Small:       0,
		AnyPair:     100,
		PerfectPair: 100,
		SuperSix:    20,
		Commission:  1,
	}

	type args struct {
		bet *roomConf.BetType0Data
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
		wantErr bool
	}{
		{
			name:    "0",
			args:    args{bet: b0},
			wantSum: 1000,
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{bet: b1},
			wantSum: 800,
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{bet: b2},
			wantSum: 620,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, err := CountBetSumType0(tt.args.bet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountBetSumType0() error = %v, wantErr %v, name %s", err, tt.wantErr, tt.name)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("CountBetSumType0() = %v, want %v, name %s", gotSum, tt.wantSum, tt.name)
			}
		})
	}
}

/*
func TestCountBetSumType1(t *testing.T) {
	type args struct {
		bet *roomConf.BetType1Data
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, err := CountBetSumType1(tt.args.bet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountBetSumType1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("CountBetSumType1() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func TestCountBetSumType2(t *testing.T) {
	type args struct {
		bet *roomConf.BetType2Data
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, err := CountBetSumType2(tt.args.bet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountBetSumType2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("CountBetSumType2() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func TestCountBetSumType6(t *testing.T) {
	type args struct {
		bet *roomConf.BetType6Data
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, err := CountBetSumType6(tt.args.bet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountBetSumType6() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("CountBetSumType6() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}

func TestCountBetSumType7(t *testing.T) {
	type args struct {
		bet *roomConf.BetType7Data
	}
	tests := []struct {
		name    string
		args    args
		wantSum int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum, err := CountBetSumType7(tt.args.bet)
			if (err != nil) != tt.wantErr {
				t.Errorf("CountBetSumType7() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSum != tt.wantSum {
				t.Errorf("CountBetSumType7() = %v, want %v", gotSum, tt.wantSum)
			}
		})
	}
}
*/
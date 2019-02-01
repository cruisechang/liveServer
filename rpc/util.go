package rpc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/control/room"
	pb "github.com/cruisechang/liveServer/protobuf"
	"github.com/cruisechang/nex/entity"
)

func (s *rpcServer) sendCommand(errorCode, step int, cmdName, dataStr string, sender entity.User, receiveConnID []string) {
	resCmd, _ := s.nex.CreateCommand(errorCode, step, cmdName, dataStr)
	s.nex.SendCommand(resCmd, sender, receiveConnID, true)
}

func getAllReceivers(users []entity.User) []string {

	receivers := []string{}

	//send to users in hall
	for _, u := range users {
		receivers = append(receivers, u.ConnID())
	}

	return receivers
}

func getHallReceivers(users []entity.User, hallID int) []string {

	receivers := []string{}

	//send to users in hall
	for _, u := range users {
		if u.HallID() == hallID {
			receivers = append(receivers, u.ConnID())
		}
	}

	return receivers
}

func getRoomReceivers(users []entity.User, roomID int) []string {

	receivers := []string{}

	//send to users in hall
	for _, u := range users {
		if u.RoomID() == roomID {
			receivers = append(receivers, u.ConnID())
		}
	}

	return receivers
}
func getStatusStart() int64 {
	return time.Now().UnixNano() / 100000
}

func getRoomStatusResponse(room entity.Room, status int, boot int, round int64, statusStart int64) ([]byte, error) {

	resData := []config.RoomStatusResData{
		{
			RoomID:      room.ID(),
			Status:      status,
			StatusStart: statusStart,
			Boot:        boot,
			Round:       round,
		},
	}

	return json.Marshal(resData)
}

func getRoundResultResponse(in interface{}, brData *room.BootRoundData) ([]byte, error) {

	if d, ok := in.(*pb.RoundResultType0Data); ok {
		r := &[]config.RoundResultType0ResData{
			{
				RoomID:      int(d.RoomID),
				Round:       brData.Round,
				Result:      int(d.Result),
				BankerPair:  int(d.BankerPair),
				PlayerPair:  int(d.PlayerPair),
				BigSmall:    int(d.BigSmall),
				AnyPair:     int(d.AnyPair),
				PerfectPair: int(d.PerfectPair),
				SuperSix:    int(d.SuperSix),
				BankerPoint: int(d.BankerPoint),
				PlayerPoint: int(d.PlayerPoint),
			},
		}

		return json.Marshal(r)
	}
	if d, ok := in.(*pb.RoundResultType1Data); ok {
		r := &[]config.RoundResultType1ResData{
			{
				RoomID:         int(d.RoomID),
				Round:          brData.Round,
				Result:         int(d.Result),
				DragonOddEven:  int(d.DragonOddEven),
				DragonRedBlack: int(d.DragonRedBlack),
				TigerOddEven:   int(d.TigerOddEven),
				TigerRedBlack:  int(d.TigerRedBlack),
			},
		}

		return json.Marshal(r)
	}

	if d, ok := in.(*pb.RoundResultType2Data); ok {

		poker0 := []int{}
		poker1 := []int{}
		poker2 := []int{}
		poker3 := []int{}
		for i, p := range d.Owner0.Poker {
			poker0 = append(poker0, int(p))
			poker1 = append(poker1, int(d.Owner1.Poker[i]))
			poker2 = append(poker2, int(d.Owner2.Poker[i]))
			poker3 = append(poker3, int(d.Owner3.Poker[i]))

		}

		r := &[]config.RoundResultType2ResData{
			{
				RoomID: int(d.RoomID),
				Round:  brData.Round,
				Head:   int(d.Head),
				Owner0: config.RoundResultType2Owner{Result: int(d.Owner0.Result), Pattern: int(d.Owner0.Pattern), Poker: poker0},
				Owner1: config.RoundResultType2Owner{Result: int(d.Owner1.Result), Pattern: int(d.Owner1.Pattern), Poker: poker1},
				Owner2: config.RoundResultType2Owner{Result: int(d.Owner2.Result), Pattern: int(d.Owner2.Pattern), Poker: poker2},
				Owner3: config.RoundResultType2Owner{Result: int(d.Owner3.Result), Pattern: int(d.Owner3.Pattern), Poker: poker3},
			},
		}

		return json.Marshal(r)
	}
	if d, ok := in.(*pb.RoundResultType6Data); ok {
		convertedDice := []int{}
		for _, dv := range d.Dice {
			convertedDice = append(convertedDice, int(dv))
		}

		//assing in.Paigow to slice
		pgSl := [][]int{}
		for i, pg := range d.Paigow {
			pgSl = append(pgSl, []int{})

			for _, pg2 := range pg.Result {
				pgSl[i] = append(pgSl[i], int(pg2))

			}
		}

		r := &[]config.RoundResultType6ResData{
			{
				RoomID:   int(d.RoomID),
				Dice:     convertedDice,
				Sum:      int(d.Sum),
				BigSmall: int(d.BigSmall),
				OddEven:  int(d.OddEven),
				Triple:   int(d.Triple),
				Paigow:   pgSl,
				Pair:     int(d.Pair),
			},
		}

		return json.Marshal(r)
	}

	if d, ok := in.(*pb.RoundResultType7Data); ok {
		r := &[]config.RoundResultType7ResData{
			{
				RoomID:   int(d.RoomID),
				Result:   int(d.Result),
				BigSmall: int(d.BigSmall),
				OddEven:  int(d.OddEven),
				RedBlack: int(d.RedBlack),
				Dozen:    int(d.Dozen),
				Column:   int(d.Column),
			},
		}

		return json.Marshal(r)
	}

	return nil, fmt.Errorf("data assertion error")
}

func getRoundRecordData(in interface{}, brData *room.BootRoundData) ([]byte, error) {

	if d, ok := in.(*pb.RoundResultType0Data); ok {
		r := &config.RoundRecordType0{
			brData.Boot,
			brData.Round,
			int(d.Result),
			int(d.BankerPair),
			int(d.PlayerPair),
			int(d.BigSmall),
			int(d.AnyPair),
			int(d.PerfectPair),
			int(d.SuperSix),
			int(d.BankerPoint),
			int(d.PlayerPoint),
		}
		return json.Marshal(r)
	}

	if d, ok := in.(*pb.RoundResultType1Data); ok {
		r := &config.RoundRecordType1{
			brData.Boot,
			brData.Round,
			int(d.Result),
			int(d.DragonOddEven),
			int(d.DragonRedBlack),
			int(d.TigerOddEven),
			int(d.TigerRedBlack),
		}
		return json.Marshal(r)
	}

	if d, ok := in.(*pb.RoundResultType2Data); ok {
		poker0 := []int{}
		poker1 := []int{}
		poker2 := []int{}
		poker3 := []int{}
		for i, p := range d.Owner0.Poker {
			poker0 = append(poker0, int(p))
			poker1 = append(poker1, int(d.Owner1.Poker[i]))
			poker2 = append(poker2, int(d.Owner2.Poker[i]))
			poker3 = append(poker3, int(d.Owner3.Poker[i]))

		}

		r := &config.RoundRecordType2{
			Boot:   brData.Boot,
			Round:  brData.Round,
			Head:   int(d.Head),
			Owner0: config.RoundResultType2Owner{Result: int(d.Owner0.Result), Pattern: int(d.Owner0.Pattern), Poker: poker0},
			Owner1: config.RoundResultType2Owner{Result: int(d.Owner1.Result), Pattern: int(d.Owner1.Pattern), Poker: poker1},
			Owner2: config.RoundResultType2Owner{Result: int(d.Owner2.Result), Pattern: int(d.Owner2.Pattern), Poker: poker2},
			Owner3: config.RoundResultType2Owner{Result: int(d.Owner3.Result), Pattern: int(d.Owner3.Pattern), Poker: poker3},
		}
		return json.Marshal(r)
	}

	if d, ok := in.(*pb.RoundResultType6Data); ok {

		convertedDice := []int{}
		for _, dv := range d.Dice {
			convertedDice = append(convertedDice, int(dv))
		}

		//assing in.Paigow to slice
		pgSl := [][]int{}
		for i, pg := range d.Paigow {
			pgSl = append(pgSl, []int{})

			for _, pg2 := range pg.Result {
				pgSl[i] = append(pgSl[i], int(pg2))

			}
		}

		r := &config.RoundRecordType6{
			Round:    brData.Round,
			Dice:     convertedDice,
			Sum:      int(d.Sum),
			BigSmall: int(d.BigSmall),
			OddEven:  int(d.OddEven),
			Triple:   int(d.Triple),
			Paigow:   pgSl,
			Pair:     int(d.Pair),
		}
		return json.Marshal(r)
	}

	if d, ok := in.(*pb.RoundResultType7Data); ok {
		r := &[]config.RoundRecordType7{
			{
				Round:    brData.Round,
				Result:   int(d.Result),
				BigSmall: int(d.BigSmall),
				OddEven:  int(d.OddEven),
				RedBlack: int(d.RedBlack),
				Dozen:    int(d.Dozen),
				Column:   int(d.Column),
			},
		}

		return json.Marshal(r)
	}

	return nil, fmt.Errorf("data assertion error")
}

func getUserResResult(userID int64, winLose int, resultCredit, balanceCredit float32) config.UserResultResDataResult {
	return config.UserResultResDataResult{
		UserID:        userID,
		Result:        winLose,
		ResultCredit:  resultCredit,
		BalanceCredit: balanceCredit,
	}
}

package rpc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	pb "github.com/cruisechang/liveServer/protobuf"

	nxLog "github.com/cruisechang/nex/log"
	"golang.org/x/net/context"
)

func (s *rpcServer) RoundProcess1(ctx context.Context, in *pb.RoundProcess1Data) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer RoundProcess1 panic:%v\n", r))
		}
	}()

	room, ok := s.nex.GetRoomManager().GetRoom(int(in.RoomID))
	if !ok {
		s.nex.GetLogger().LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer RoundProcess1 roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	receivers := getHallReceivers(s.nex.GetUsers(), room.HallID())

	//send to client
	if len(receivers) > 0 {
		toClient := []config.RoundProcess1ResData{
			{
				RoomID: int(in.RoomID),
				Poker:  int(in.Poker),
				Index:  int(in.Index),
				Owner:  int(in.Owner),
			},
		}

		toClientJ, _ := json.Marshal(toClient)
		sendDataStr := base64.StdEncoding.EncodeToString(toClientJ)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoundProcess1(), sendDataStr, nil, receivers)

		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer RoundProcess1 toClient=%s", string(toClientJ)))
	}
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer RoundProcess1 complete"))

	return &pb.Empty{}, nil
}

func (s *rpcServer) RoundResultType1(ctx context.Context, in *pb.RoundResultType1Data) (*pb.Empty, error) {
	logPrefix := "rpcServer RoundResultType1"
	rr := newRoundResultHandler(s, logPrefix, in.RoomID, in.Round, in)
	return rr.handle()
}

func (s *rpcServer) UpdateResultType1(ctx context.Context, in *pb.UpdateResultType1Data) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer UpdateResultType1 panic:%v", r))
		}
	}()

	//targetHallID := int(in.Round)
	//receivers := []string{}
	//
	//users := s.nex.GetUsers()
	//
	////send to users in hall
	//for _, u := range users {
	//
	//	if u.HallID() == targetHallID {
	//		receivers = append(receivers, u.ConnID())
	//	}
	//}
	//
	//if len(receivers) > 0 {
	//	resData := []config.RoundResultType1ResData{
	//		config.RoundResultType1ResData{
	//			RoomID:         int(in.RoomID),
	//			Round:          in.Round,
	//			Result:         int(in.Result),
	//			DragonPoker:    in.DragonPoker,
	//			DragonRedBlack: int(in.DragonRedBlack),
	//			TigerPoker:     in.TigerPoker,
	//			DragonOddEven:  int(in.DragonOddEven),
	//			TigerRedBlack:  int(in.TigerRedBlack),
	//		},
	//	}
	//
	//	b, _ := json.Marshal(resData)
	//	sendDataStr := base64.StdEncoding.EncodeToString(b)
	//	s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoundResult1(), sendDataStr, nil, receivers)
	//}

	return &pb.Empty{}, nil

}

func (s *rpcServer) HistoryResultType1(ctx context.Context, in *pb.HistoryResultType1Data) (*pb.HistoryResultType1Res, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer HistoryResultType1 panic:%v\n", r))
		}
	}()

	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer HistoryResultType1 roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	result, err := s.roomCtrl.GetHistoryResult(r)
	if err != nil {
		return nil, err
	}
	re, ok := result.(roomConf.HistoryResultType1)
	if !ok {
		return nil, errors.New("GetHistoryResultType1 assertion error")
	}

	return &pb.HistoryResultType1Res{
		Result: re,
	}, nil
}

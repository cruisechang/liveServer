package rpc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	pb "github.com/cruisechang/liveServer/protobuf"
	"github.com/cruisechang/nex/entity"
	nxLog "github.com/cruisechang/nex/log"
	"golang.org/x/net/context"
)

func (s *rpcServer) RoundProcess2(ctx context.Context, in *pb.RoundProcess2Data) (*pb.Empty, error) {
	logPrefix := "rpcServer RoundProcess2"
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("%s panic:%v\n", logPrefix, r))
		}
	}()

	room, ok := s.nex.GetRoomManager().GetRoom(int(in.RoomID))
	if !ok {
		s.nex.GetLogger().LogFile(nxLog.LevelError, fmt.Sprintf("%s roomID not found :%d\n", logPrefix, in.RoomID))
		return nil, errors.New("roomID not found")
	}

	receivers := getHallReceivers(s.nex.GetUsers(), room.HallID())

	//send to client
	if len(receivers) > 0 {
		toClient := []config.RoundProcess2ResData{
			{
				RoomID: int(in.RoomID),
				Poker:  int(in.Poker),
				Index:  int(in.Index),
				Owner:  int(in.Owner),
			},
		}

		toClientJ, _ := json.Marshal(toClient)
		sendDataStr := base64.StdEncoding.EncodeToString(toClientJ)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoundProcess2(), sendDataStr, nil, receivers)

		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s toClient=%s", logPrefix, string(toClientJ)))
	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("%s compplete", logPrefix))

	return &pb.Empty{}, nil
}

func (s *rpcServer) RoundResultType2(ctx context.Context, in *pb.RoundResultType2Data) (*pb.Empty, error) {
	logPrefix := "rpcServer RoundResultType2"

	rr := newRoundResultHandler(s, logPrefix, in.RoomID, in.Round, in)
	return rr.handle()
}

func (s *rpcServer) addHistoryResult2(logPrefix string, room entity.Room, data interface{}) error {
	//room result 增加
	if in, ok := data.(*pb.RoundResultType2Data); ok {
		return s.roomCtrl.AddHistoryResultType(room, []int32{int32(in.Owner0.Result), int32(in.Owner0.Pattern), int32(in.Owner1.Result), int32(in.Owner1.Pattern), int32(in.Owner2.Result), int32(in.Owner2.Pattern), int32(in.Owner3.Result), int32(in.Owner3.Pattern)})

	}
	return fmt.Errorf("%s addHistoryFunc data assertion error=", logPrefix)
}

func (s *rpcServer) UpdateResultType2(ctx context.Context, in *pb.UpdateResultType2Data) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer UpdateResultType2 panic:%v\n", r))
		}
	}()

	return &pb.Empty{}, nil
}
func (s *rpcServer) HistoryResultType2(ctx context.Context, in *pb.HistoryResultType2Data) (*pb.HistoryResultType2Res, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer HistoryResultType2 panic:%v\n", r))
		}
	}()

	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer HistoryResultType2 roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	result, err := s.roomCtrl.GetHistoryResult(r)
	if err != nil {
		return nil, err
	}
	re, ok := result.(roomConf.HistoryResultType2)
	if !ok {
		return nil, errors.New("HistoryResultType2 assertion error")
	}

	var pbAry []*pb.HistoryResultType2ResInnerType
	for _, v := range re {
		pbAry = append(pbAry, &pb.HistoryResultType2ResInnerType{Result: v})

	}

	return &pb.HistoryResultType2Res{
		Result: pbAry,
	}, nil

}

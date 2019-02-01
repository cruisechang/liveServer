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

func (s *rpcServer) RoundResultType7(ctx context.Context, in *pb.RoundResultType7Data) (*pb.Empty, error) {
	logPrefix := "rpcServer RoundResultType6"

	rr := newRoundResultHandler(s, logPrefix, in.RoomID, in.Round, in)
	return rr.handle()
}

func (s *rpcServer) addHistoryResult7(logPrefix string, room entity.Room, data interface{}) error {
	//room result 增加
	if in, ok := data.(*pb.RoundResultType7Data); ok {
		return s.roomCtrl.AddHistoryResultType(room, in.Result)
	}
	return fmt.Errorf("%s addHistoryFunc data assertion error=", logPrefix)
}

func (s *rpcServer) UpdateResultType7(ctx context.Context, in *pb.UpdateResultType7Data) (*pb.Empty, error) {

	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer UpdateResultType7 panic:%v\n", r))
		}
	}()

	return &pb.Empty{}, nil
}

func (s *rpcServer) HistoryResultType7(ctx context.Context, in *pb.HistoryResultType7Data) (*pb.HistoryResultType7Res, error) {

	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer HistoryResultType7 panic:%v\n", r))
		}
	}()

	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer HistoryResultType7 roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	result, err := s.roomCtrl.GetHistoryResult(r)
	if err != nil {
		return nil, err
	}
	re, ok := result.(roomConf.HistoryResultType7)
	if !ok {
		return nil, errors.New("HistoryResultType7 assertion error")
	}

	return &pb.HistoryResultType7Res{
		Result: re,
	}, nil

}

func (s *rpcServer) Rethrow(ctx context.Context, in *pb.RethrowData) (*pb.Empty, error) {

	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer Rethrow panic:%v\n", r))
		}
	}()

	rm := s.nex.GetRoomManager()
	room, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer Rethrow roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	hallReceivers := getHallReceivers(s.nex.GetUsers(), room.HallID())

	//send to client
	if len(hallReceivers) > 0 {
		toClient := []config.RethrowResData{
			{
				HallID: room.HallID(),
				RoomID: room.ID(),
				Round:  in.Round,
			},
		}

		toClientJ, _ := json.Marshal(toClient)
		sendDataStr := base64.StdEncoding.EncodeToString(toClientJ)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRethrow(), sendDataStr, nil, hallReceivers)

		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer Rethrow toClient=%s", string(toClientJ)))
	}
	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer Rethrow complete"))

	return &pb.Empty{}, nil

}

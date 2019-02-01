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

func (s *rpcServer) RoundProcess0(ctx context.Context, in *pb.RoundProcess0Data) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer RoundProcess0 panic:%v\n", r))
		}
	}()

	room, ok := s.nex.GetRoomManager().GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer RoundProcess0 roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	receivers := getHallReceivers(s.nex.GetUsers(), room.HallID())

	//send to client
	if len(receivers) > 0 {
		toClient := []config.RoundProcess0ResData{
			{
				RoomID: int(in.RoomID),
				Poker:  int(in.Poker),
				Index:  int(in.Index),
				Owner:  int(in.Owner),
			},
		}

		toClientJ, _ := json.Marshal(toClient)
		sendDataStr := base64.StdEncoding.EncodeToString(toClientJ)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRoundProcess0(), sendDataStr, nil, receivers)

		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer RoundProcess0  =%s", string(toClientJ)))
	}

	return &pb.Empty{}, nil
}

func (s *rpcServer) RoundResultType0(ctx context.Context, in *pb.RoundResultType0Data) (*pb.Empty, error) {

	logPrefix := "rpcServer RoundResultType0"

	rr := newRoundResultHandler(s, logPrefix, in.RoomID, in.Round, in)
	return rr.handle()
}

func (s *rpcServer) UpdateResultType0(ctx context.Context, in *pb.UpdateResultType0Data) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer UpdateResultType0 panic:%v\n", r))
		}
	}()

	return &pb.Empty{}, nil
}

func (s *rpcServer) HistoryResultType0(ctx context.Context, in *pb.HistoryResultType0Data) (*pb.HistoryResultType0Res, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer HistoryResultType0 panic:%v\n", r))
		}
	}()

	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer HistoryResultType0 roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	result, err := s.roomCtrl.GetHistoryResult(r)
	if err != nil {
		return nil, err
	}
	re, ok := result.(roomConf.HistoryResultType0)
	if !ok {
		return nil, errors.New("HistoryResultType0 assertion error")
	}

	var pbAry []*pb.HistoryResultType0ResInnerType
	for _, v := range re {
		pbAry = append(pbAry, &pb.HistoryResultType0ResInnerType{Result: v})

	}

	return &pb.HistoryResultType0Res{
		Result: pbAry,
	}, nil

}

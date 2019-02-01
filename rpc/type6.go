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

func (s *rpcServer) RoundResultType6(ctx context.Context, in *pb.RoundResultType6Data) (*pb.Empty, error) {
	logPrefix := "rpcServer RoundResultType6"

	rr := newRoundResultHandler(s, logPrefix, in.RoomID, in.Round, in)
	return rr.handle()
}

func (s *rpcServer) addHistoryResult6(logPrefix string, room entity.Room, data interface{}) error {
	//room result 增加
	if in, ok := data.(*pb.RoundResultType6Data); ok {

		convertedDice := []int{}
		for _, dv := range in.Dice {
			convertedDice = append(convertedDice, int(dv))
		}

		//room result 增加
		hrt6 := &roomConf.HistoryResultType6{
			HallID:   room.HallID(),
			RoomID:   room.ID(),
			Dice:     convertedDice,
			Sum:      int(in.Sum),
			BigSmall: int(in.BigSmall),
			OddEven:  int(in.OddEven),
		}
		return s.roomCtrl.AddHistoryResultType(room, hrt6)

	}
	return fmt.Errorf("%s addHistoryFunc data assertion error=", logPrefix)
}

func (s *rpcServer) UpdateResultType6(ctx context.Context, in *pb.UpdateResultType6Data) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer UpdateResultType6 panic:%v\n", r))
		}
	}()

	return &pb.Empty{}, nil
}

func (s *rpcServer) HistoryResultType6(ctx context.Context, in *pb.HistoryResultType6Data) (*pb.HistoryResultType6Res, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer HistoryResultType6 panic:%v\n", r))
		}
	}()

	rm := s.nex.GetRoomManager()
	r, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer HistoryResultType6 roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	result, err := s.roomCtrl.GetHistoryResult(r)
	if err != nil {
		return nil, err
	}
	re, ok := result.([]*roomConf.HistoryResultType6)
	if !ok {
		return nil, errors.New("HistoryResultType6 assertion error")
	}

	var pbAry []int32
	for _, v := range re {
		pbAry = append(pbAry, int32(v.OddEven))

	}

	return &pb.HistoryResultType6Res{
		Result: pbAry,
	}, nil

}

func (s *rpcServer) RerollDice(ctx context.Context, in *pb.RerollDiceData) (*pb.Empty, error) {
	logger := s.nex.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.LogFile(nxLog.LevelPanic, fmt.Sprintf("rpcServer RerollDice panic:%v\n", r))
		}
	}()

	rm := s.nex.GetRoomManager()
	room, ok := rm.GetRoom(int(in.RoomID))
	if !ok {
		logger.LogFile(nxLog.LevelError, fmt.Sprintf("rpcServer RerollDice roomID not found :%d\n", in.RoomID))
		return nil, errors.New("roomID not found")
	}

	hallReceivers := getHallReceivers(s.nex.GetUsers(), room.HallID())

	//send to client
	if len(hallReceivers) > 0 {
		toClient := []config.RerollDiceResData{
			{
				HallID: room.HallID(),
				RoomID: room.ID(),
				Round:  in.Round,
			},
		}

		toClientJ, _ := json.Marshal(toClient)
		sendDataStr := base64.StdEncoding.EncodeToString(toClientJ)
		s.sendCommand(config.CodeSuccess, 0, s.configure.CmdRerollDice(), sendDataStr, nil, hallReceivers)
		logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer RerollDice toClient=%s", string(toClientJ)))
	}

	logger.LogFile(nxLog.LevelInfo, fmt.Sprintf("rpcServer RerollDice complete"))

	return &pb.Empty{}, nil

}

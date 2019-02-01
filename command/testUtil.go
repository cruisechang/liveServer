package command

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/control"
	roomCtrl "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex"
)

func getBasicProcessor() BasicProcessor {
	nx, _ := nex.NewNex(getConfigFilePosition("nexConfig.json"))
	conf, _ := config.NewConfigurer("config.json")
	dbCtrl := control.NewDBController(conf.DBAPIServer())
	//room,_:=nx.GetRoomManager().CreateRoom(0,conf.RoomType7(),"room")
	rCtrl := roomCtrl.NewController(conf)
	//rCtrl.InitRoomBet(room )
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())
	return NewBasicProcessor(nx, conf, dbCtrl, rmc)
}

func getBasicProcessorByRCTRL(nx nex.Nex,rCtrl *roomCtrl.Controller) BasicProcessor {
	conf, _ := config.NewConfigurer("config.json")
	dbCtrl := control.NewDBController(conf.DBAPIServer())
	//room,_:=nx.GetRoomManager().CreateRoom(0,conf.RoomType7(),"room")
	//rCtrl := roomCtrl.NewController(conf)
	//rCtrl.InitRoomBet(room )
	rmc := control.NewRoadMapController(conf.RoadMapAPIHost(), rCtrl, nx.GetLogger())
	return NewBasicProcessor(nx, conf, dbCtrl, rmc)
}

func getConfigFilePosition(fileName string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}

	var buf bytes.Buffer
	buf.WriteString(dir)
	buf.WriteString("/")
	buf.WriteString(fileName)

	return buf.String()
}

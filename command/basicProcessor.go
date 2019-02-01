package command

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/control"
	"github.com/cruisechang/nex"
	"github.com/cruisechang/nex/entity"
	nxLog "github.com/cruisechang/nex/log"
)

//interface for parameters
//type Configurer interface {
//	GetVersion() string
//}

type BasicProcessor interface {
	SendCommand(errorCode, step int, cmdName, dataStr string, sender entity.User, receiveConnID []string)
	DisconnectUser(userID int)
	RemoveUser(userID int)
	GetUsers() []entity.User
	GetUser(userID int) (entity.User, bool)
	GetAllReceivers() []string
	GetHallReceivers(hallID int) []string
	GetRoomReceivers(roomID int) []string

	GetLogger() nxLog.Logger
	GetConfigurer() config.Configurer
	GetHallManager() nex.HallManager
	GetRoomManager() nex.RoomManager

	DefaultSendData() string

	SetDBAPIHost(host string)
	DBAPIDo(method, path string, body io.Reader) (*http.Response, error)

	RequestRoadMap(postData []byte) ([]byte, error)
	RoadMapDataType0(roomID int) (*config.RoadMapType0ResData, bool)
	RoadMapDataType1(roomID int) (*config.RoadMapType1ResData, bool)
	RoadMapDataType2(roomID int) (*config.RoadMapType2ResData, bool)
	RoadMapDataType6(roomID int) (*config.RoadMapType6ResData, bool)
	RoadMapDataType7(roomID int) (*config.RoadMapType7ResData, bool)
	SetRoadMapDataFromResult(roomID, roomType int, resultStr string)
	GetRoadMapRequestData(hallID, roomID, roomType int, historyResult interface{}) ([]byte, error)
}

//BasicProcessor is parent struct for process.
type basicProcessor struct {
	nex             nex.Nex
	configurer      config.Configurer
	dbCtrl          *control.DBController
	roadMapCtrl     *control.RoadMapController
	defaultSendData string
}

func NewBasicProcessor(nex nex.Nex, conf config.Configurer, db *control.DBController, rmc *control.RoadMapController) BasicProcessor {
	//response data
	resData := []config.EmptyResData{{}}
	b, _ := json.Marshal(resData)

	basic := &basicProcessor{
		nex:             nex,
		configurer:      conf,
		defaultSendData: base64.StdEncoding.EncodeToString(b),
		dbCtrl:          db,
		roadMapCtrl:     rmc,
	}

	return basic
}

func (b *basicProcessor) DefaultSendData() string {
	return b.defaultSendData
}
func (b *basicProcessor) GetLogger() nxLog.Logger {
	return b.nex.GetLogger()
}

func (b *basicProcessor) GetConfigurer() config.Configurer {
	return b.configurer
}

func (b *basicProcessor) GetHallManager() nex.HallManager {
	return b.nex.GetHallManager()
}

func (b *basicProcessor) GetRoomManager() nex.RoomManager {
	return b.nex.GetRoomManager()
}

func (b *basicProcessor) Print(msg string) {
	fmt.Printf(msg)
}

func (b *basicProcessor) SendCommand(errorCode, step int, cmdName, dataStr string, sender entity.User, receiveConnID []string) {

	resCmd, _ := b.nex.CreateCommand(errorCode, step, cmdName, dataStr)

	b.nex.SendCommand(resCmd, sender, receiveConnID, true)
}

func (b *basicProcessor) DisconnectUser(userID int) {
	b.nex.DisconnectUser(userID)
}

func (b *basicProcessor) RemoveUser(userID int) {
	b.nex.RemoveUser(userID)
}
func (b *basicProcessor) GetUsers() []entity.User {
	return b.nex.GetUsers()
}
func (b *basicProcessor) GetUser(userID int) (entity.User, bool) {
	return b.nex.GetUser(userID)
}

//func (b *basicProcessor) GetUserConnIDs()[]string {
//	return b.nex.GetUserConnIDs()
//}

func (b *basicProcessor) GetAllReceivers() []string {
	users := b.nex.GetUsers()
	receivers := []string{}

	//send to users in hall
	for _, u := range users {
		receivers = append(receivers, u.ConnID())
	}

	return receivers
}

func (b *basicProcessor) GetHallReceivers(hallID int) []string {

	receivers := []string{}

	users := b.nex.GetUsers()
	//send to users in hall
	for _, u := range users {
		if u.HallID() == hallID {
			receivers = append(receivers, u.ConnID())
		}
	}

	return receivers
}

func (b *basicProcessor) GetRoomReceivers(roomID int) []string {

	receivers := []string{}

	users := b.nex.GetUsers()
	//send to users in hall
	for _, u := range users {
		if u.RoomID() == roomID {
			receivers = append(receivers, u.ConnID())
		}
	}

	return receivers
}

//db get/patch hall, get/patch room  , get/patch user
func (b *basicProcessor) SetDBAPIHost(host string) {
	b.dbCtrl.SetDBAPIHost(host)
}

func (b *basicProcessor) DBAPIDo(method, path string, body io.Reader) (*http.Response, error) {
	return b.dbCtrl.Do(method, path, body)
}

//road map
func (b *basicProcessor) RequestRoadMap(postData []byte) ([]byte, error) {
	return b.roadMapCtrl.RequestRoadMap(postData)
}
func (b *basicProcessor) RoadMapDataType0(roomID int) (*config.RoadMapType0ResData, bool) {
	return b.roadMapCtrl.RoadMapDataType0(roomID)
}
func (b *basicProcessor) RoadMapDataType1(roomID int) (*config.RoadMapType1ResData, bool) {
	return b.roadMapCtrl.RoadMapDataType1(roomID)
}
func (b *basicProcessor) RoadMapDataType2(roomID int) (*config.RoadMapType2ResData, bool) {
	return b.roadMapCtrl.RoadMapDataType2(roomID)
}
func (b *basicProcessor) RoadMapDataType6(roomID int) (*config.RoadMapType6ResData, bool) {
	return b.roadMapCtrl.RoadMapDataType6(roomID)
}
func (b *basicProcessor) RoadMapDataType7(roomID int) (*config.RoadMapType7ResData, bool) {
	return b.roadMapCtrl.RoadMapDataType7(roomID)
}

func (b *basicProcessor) SetRoadMapDataFromResult(roomID, roomType int, resultStr string) {
	b.roadMapCtrl.SetRoadMapDataFromResult(roomID, roomType, resultStr)
}

func (b *basicProcessor) GetRoadMapRequestData(hallID, roomID, roomType int, historyResult interface{}) ([]byte, error) {
	return b.roadMapCtrl.GetRoadMapRequestData(hallID, roomID, roomType, historyResult)
}

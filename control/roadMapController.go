package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	rc "github.com/cruisechang/liveServer/control/room"
	"github.com/cruisechang/nex/entity"
	nxLog "github.com/cruisechang/nex/log"
)

//NewRoadMapController returns RoadMapController structure
func NewRoadMapController(targetHost string, roomCtrl *rc.Controller, logger nxLog.Logger) *RoadMapController {
	return &RoadMapController{
		targetHost:   targetHost,
		roomCtrl:     roomCtrl,
		logger:       logger,
		roadMapData0: make(map[int]*config.RoadMapType0ResData),
		roadMapData1: make(map[int]*config.RoadMapType1ResData),
		roadMapData2: make(map[int]*config.RoadMapType2ResData),
		roadMapData6: make(map[int]*config.RoadMapType6ResData),
		roadMapData7: make(map[int]*config.RoadMapType7ResData),
	}
}

//RoadMapController 負責處理路單
//取room history資料
//post給路單server
//存起來
type RoadMapController struct {
	targetHost   string
	roomCtrl     *rc.Controller
	logger       nxLog.Logger
	roadMapData0 map[int]*config.RoadMapType0ResData //從roadMap server取回的資料，要回傳給client,用roomID當key
	roadMapData1 map[int]*config.RoadMapType1ResData
	roadMapData2 map[int]*config.RoadMapType2ResData
	roadMapData6 map[int]*config.RoadMapType6ResData
	roadMapData7 map[int]*config.RoadMapType7ResData
}

//InitRoadMapData gets road map data from road map server when server start.
func (c *RoadMapController) InitRoadMapData(rooms []entity.Room) error {

	logPrefix := "RoadMapController InitRoadMapData "

	for _, v := range rooms {

		hrs, err := c.roomCtrl.GetHistoryResult(v)
		if err != nil {
			//c.logger.LogFile(nxLog.LevelError, fmt.Sprintf("%s GetHistoryResult() error=%s", logPrefix, err.Error()))
			return fmt.Errorf("%s roomCtroller.GetHistoryResult() error=%s", logPrefix, err.Error())
		}

		postData, err := c.GetRoadMapRequestData(v.HallID(), v.ID(), v.Type(), hrs)
		if err != nil {
			return fmt.Errorf("%s getRoadMapRequestData() error=%s", logPrefix, err.Error())
		}

		body, err := c.RequestRoadMap(postData)
		if err != nil {
			return fmt.Errorf("%s requestRoadMap() error=%s", logPrefix, err.Error())
		}

		resData := &config.RoadMapResponse{}

		err = json.Unmarshal(body, resData)
		if err != nil {
			return fmt.Errorf("%s json.Unmarshal() response body error=%s", logPrefix, err.Error())
		}

		//resData.RoomType
		if resData.Code != config.CodeSuccess {
			return fmt.Errorf("%s roadMapResponse code error got=%d, got body=%s post data=%s", logPrefix, resData.Code, string(body), string(postData))
		}

		 err = c.SetRoadMapDataFromResult(v.ID(), v.Type(), resData.Result)
		if err != nil {
			return fmt.Errorf("%s result2Data() error=%s", logPrefix, err.Error())
		}

		log.Printf("roadMap resData result=%s", resData.Result)

	}

	return nil
}

//取畫路單的資料
func (c *RoadMapController) GetRoadMapRequestData(hallID, roomID, roomType int, historyResult interface{}) ([]byte, error) {

	postData := &config.RoadMapRequest{
		HallID:   hallID,
		RoomID:   roomID,
		RoomType: roomType,
		Result:   "",
	}

	//type6 要改，其他不用改
	if hrt, ok := historyResult.([]*roomConf.HistoryResultType6); ok {

		trans := []*config.HistoryResultType6{}

		for _, v := range hrt {

			trans = append(trans, &config.HistoryResultType6{
				Dice:     v.Dice,
				Sum:      v.Sum,
				BigSmall: v.BigSmall,
				OddEven:  v.OddEven,
			})
		}

		mb, err := json.Marshal(trans)
		if err != nil {
			return nil, fmt.Errorf("")
		}
		postData.Result = string(mb)
	} else {
		mb, _ := json.Marshal(historyResult)
		postData.Result = string(mb)
	}

	return json.Marshal(postData)
}
func (c *RoadMapController) request(body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.targetHost, body)
	if err != nil {

	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Do(req)
}
func (c *RoadMapController) RequestRoadMap(postData []byte) ([]byte, error) {
	res, err := c.request(bytes.NewBuffer(postData))
	if err != nil {
		return nil, err

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}



func (c *RoadMapController) RoadMapDataType0(roomID int) (*config.RoadMapType0ResData, bool) {
	if v, ok := c.roadMapData0[roomID]; ok {
		return v, true
	}
	return nil, false
}
func (c *RoadMapController) RoadMapDataType1(roomID int) (*config.RoadMapType1ResData, bool) {
	if v, ok := c.roadMapData1[roomID]; ok {
		return v, true
	}
	return nil, false
}
func (c *RoadMapController) RoadMapDataType2(roomID int) (*config.RoadMapType2ResData, bool) {
	if v, ok := c.roadMapData2[roomID]; ok {
		return v, true
	}
	return nil, false
}
func (c *RoadMapController) RoadMapDataType6(roomID int) (*config.RoadMapType6ResData, bool) {
	if v, ok := c.roadMapData6[roomID]; ok {
		return v, true
	}
	return nil, false
}
func (c *RoadMapController) RoadMapDataType7(roomID int) (*config.RoadMapType7ResData, bool) {
	if v, ok := c.roadMapData7[roomID]; ok {
		return v, true
	}
	return nil, false
}

/*
func (c *RoadMapController) SetRoadMapData(roomID int, roomType int, data interface{}) error {

	switch roomType {
	case 0:
		if d, ok := data.(*config.RoadMapType0ResData); ok {
			c.roadMapData0[roomID] = d
		}
		return nil
	case 1:
		if d, ok := data.(*config.RoadMapType1ResData); ok {
			c.roadMapData1[roomID] = d
		}
		return nil
	case 2:
		if d, ok := data.(*config.RoadMapType2ResData); ok {
			c.roadMapData2[roomID] = d
		}
		return nil
	case 6:
		if d, ok := data.(*config.RoadMapType6ResData); ok {
			c.roadMapData6[roomID] = d
		}
		return nil
	case 7:
		if d, ok := data.(*config.RoadMapType7ResData); ok {
			c.roadMapData7[roomID] = d
		}
		return nil
	default:
		return fmt.Errorf("room type error got=%d", roomType)
	}
}
*/

func (c *RoadMapController) SetRoadMapDataFromResult(roomID, roomType int, jsonStr string) error {

	switch roomType {
	case 0:
		rm := &config.RoadMapType0ResData{}
		err := json.Unmarshal([]byte(jsonStr), rm)
		if err != nil {
			return err
		}
		c.roadMapData0[roomID] = rm
		return nil
	case 1:
		rm := &config.RoadMapType1ResData{}
		err := json.Unmarshal([]byte(jsonStr), rm)
		if err != nil {
			return err
		}
		c.roadMapData1[roomID] = rm
		return nil
	case 2:
		rm := &config.RoadMapType2ResData{}
		err := json.Unmarshal([]byte(jsonStr), rm)
		if err != nil {
			return err
		}
		c.roadMapData2[roomID] = rm
		return nil
	case 6:
		rm := &config.RoadMapType6ResData{}
		err := json.Unmarshal([]byte(jsonStr), rm)
		if err != nil {
			return err
		}
		c.roadMapData6[roomID] = rm
		return nil
	case 7:
		rm := &config.RoadMapType7ResData{}
		err := json.Unmarshal([]byte(jsonStr), rm)
		if err != nil {
			return err
		}
		c.roadMapData7[roomID] = rm
		return nil
	default:
		return fmt.Errorf("room type not found=%d", roomType)
	}
}

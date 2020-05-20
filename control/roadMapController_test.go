package control

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cruisechang/liveServer/config"
	"github.com/cruisechang/liveServer/config/roomConf"
	rc "github.com/cruisechang/liveServer/control/room"
)

func createRoadMapController() *RoadMapController {

	conf, _ := config.NewConfigurer("config.json")
	rCtrl := rc.NewController(conf)
	return NewRoadMapController(conf.RoadMapAPIHost(), rCtrl)

}

/*
func TestRoadMapController_InitRoadMapData(t *testing.T) {
	type fields struct {
		targetHost  string
		roomCtrl    *rc.Controller
		roadMapData map[int][]byte
	}
	type args struct {
		rooms []entity.Room
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RoadMapController{
				targetHost:  tt.fields.targetHost,
				roomCtrl:    tt.fields.roomCtrl,
				roadMapData: tt.fields.roadMapData,
			}
			if err := c.InitRoadMapData(tt.args.rooms); (err != nil) != tt.wantErr {
				t.Errorf("RoadMapController.InitRoadMapData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
*/
//func TestRoadMapController_getHistoryResultPostData(t *testing.T) {
//	type fields struct {
//		targetHost  string
//		roomCtrl    *rc.Controller
//		roadMapData map[int][]byte
//	}
//	type args struct {
//		hallID        int
//		roomID        int
//		roomType      int
//		historyResult interface{}
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    []byte
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &RoadMapController{
//				targetHost:  tt.fields.targetHost,
//				roomCtrl:    tt.fields.roomCtrl,
//				roadMapData: tt.fields.roadMapData,
//			}
//			got, err := c.getRoadMapRequestData(tt.args.hallID, tt.args.roomID, tt.args.roomType, tt.args.historyResult)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("RoadMapController.getHistoryResultPostData() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("RoadMapController.getHistoryResultPostData() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestRoadMapController_RoadMapData(t *testing.T) {
//	type fields struct {
//		targetHost  string
//		roomCtrl    *rc.Controller
//		roadMapData map[int][]byte
//	}
//	type args struct {
//		roomID int
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   []byte
//		want1  bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &RoadMapController{
//				targetHost:  tt.fields.targetHost,
//				roomCtrl:    tt.fields.roomCtrl,
//				roadMapData: tt.fields.roadMapData,
//			}
//			got, got1 := c.RoadMapData(tt.args.roomID)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("RoadMapController.RoadMapData() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("RoadMapController.RoadMapData() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func TestRoadMapController_SetRoadMapData(t *testing.T) {
//	type fields struct {
//		targetHost  string
//		roomCtrl    *rc.Controller
//		roadMapData map[int][]byte
//	}
//	type args struct {
//		roomID int
//		data   []byte
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := &RoadMapController{
//				targetHost:  tt.fields.targetHost,
//				roomCtrl:    tt.fields.roomCtrl,
//				roadMapData: tt.fields.roadMapData,
//			}
//			c.SetRoadMapData(tt.args.roomID, tt.args.data)
//		})
//	}
//}

func TestRoadMapController_Request(t *testing.T) {

	rmc := createRoadMapController()

	data0 := createHistoryResult(0)
	data1 := createHistoryResult(1)
	data2 := createHistoryResult(2)
	data6 := createHistoryResult(6)
	data7 := createHistoryResult(7)

	bd0, _ := rmc.GetRoadMapRequestData(100, 1, 0, data0)
	bd1, _ := rmc.GetRoadMapRequestData(100, 1, 1, data1)
	bd2, _ := rmc.GetRoadMapRequestData(100, 1, 2, data2)
	bd6, _ := rmc.GetRoadMapRequestData(100, 1, 6, data6)
	bd7, _ := rmc.GetRoadMapRequestData(100, 1, 7, data7)

	tests := []struct {
		name    string
		code int
		reqData *bytes.Buffer
		wantErr bool
	}{

		{
			name:    "0",
			code:config.CodeSuccess,
			reqData: bytes.NewBuffer(bd0),
			wantErr: false,
		},
		{
			name:    "1",
			code:config.CodeSuccess,
			reqData: bytes.NewBuffer(bd1),
			wantErr: false,
		},
		{
			name:    "2",
			code:config.CodeSuccess,
			reqData: bytes.NewBuffer(bd2),
			wantErr: false,
		},
		{
			name:    "3",
			code:config.CodeSuccess,
			reqData: bytes.NewBuffer(bd6),
			wantErr: false,
		},
		{
			name:    "4",
			code:config.CodeSuccess,
			reqData: bytes.NewBuffer(bd7),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := rmc.Request(tt.reqData)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RoadMapController.Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if res.StatusCode != http.StatusOK {

				t.Fatalf("http status error, got %d want %d  name=%s", res.StatusCode, http.StatusOK, tt.name)
			}

			body, _ := ioutil.ReadAll(res.Body)

			resData := &config.RoadMapResponse{}
			err = json.Unmarshal(body, resData)
			if err != nil {
				t.Fatalf("handler unmarshal responseData error=%s name=%s", err.Error(), tt.name)
			}

			if resData.Code!=tt.code{

				t.Fatalf("handle responseData code error got %d want %d name=%s", resData.Code,tt.code, tt.name)
			}

			t.Logf("resData result=%s",resData.Result)
		})
	}
}

func createHistoryResult(roomType int) interface{} {
	switch roomType {
	case 0:
		return roomConf.HistoryResultType0{
			{2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1}, {2, 0, 1},
			{2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0},
			{2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0}, {2, 1, 0},
		}
	case 1:
		return roomConf.HistoryResultType1{
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1, 2, 1,
			0, 1, 2, 0, 1, 2, 0, 1,
		}
	case 2:
		return roomConf.HistoryResultType2{
			{1, 1, 0, 1, 0, 1, 1, 1}, {1, 1, 0, 2, 1, 2, 1, 2}, {1, 1, 0, 3, 0, 3, 1, 3}, {1, 1, 1, 4, 0, 4, 1, 4}, {1, 1, 0, 5, 0, 5, 1, 5},
			{1, 1, 0, 6, 1, 6, 1, 6}, {1, 1, 1, 7, 1, 7, 1, 7}, {1, 1, 0, 8, 1, 8, 1, 8}, {1, 1, 1, 9, 0, 9, 1, 9}, {1, 1, 0, 10, 0, 10, 1, 10},
			{1, 1, 1, 11, 0, 11, 1, 11},
		}
	case 6:
		var hrt6 []*roomConf.HistoryResultType6
		hrt6 = append(hrt6,
			&roomConf.HistoryResultType6{
				HallID:   0,
				RoomID:   1,
				Dice:     []int{1, 2, 2},
				Sum:      5,
				BigSmall: 0,
				OddEven:  1,
			})

		hrt6 = append(hrt6,
			&roomConf.HistoryResultType6{
				HallID:   0,
				RoomID:   1,
				Dice:     []int{1, 2, 3},
				Sum:      6,
				BigSmall: 0,
				OddEven:  0,
			})

		hrt6 = append(hrt6,
			&roomConf.HistoryResultType6{
				HallID:   0,
				RoomID:   1,
				Dice:     []int{4, 2, 3},
				Sum:      9,
				BigSmall: 1,
				OddEven:  1,
			})

		return hrt6
	case 7:
		return roomConf.HistoryResultType7{1, 0, 10, 33, 25, 13, 22, 1, 9, 26, 30, 30, 30}
	default:
		return nil
	}
}

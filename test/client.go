package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/cruisechang/nex"

	"github.com/cruisechang/goutil/net"
	gorillaWebsocket "github.com/gorilla/websocket"
)

var addr = flag.String("addr", "211.75.180.91:17003", "http service address")

//var addr = flag.String("addr", "127.0.0.1:17003", "http service address")

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("client panic:%v", r)
		}
	}()

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
	log.Printf("connecting to %s", u.String())

	//connect
	c, _, err := gorillaWebsocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	//wait for 3 sec
	time.Sleep(time.Duration(time.Second * 1))

	log.Printf("GetSessionID\n")

	GetSessionID(c)

	go func() {

		//recover panic
		defer func() {
			if r := recover(); r != nil {
				c.Close()
				c = nil
				log.Printf("go panic:%v", r)
			}
		}()

		for {

			_, resBuff, _ := c.ReadMessage()
			b, _ := base64.StdEncoding.DecodeString(string(resBuff))

			cmd := &nex.Command{}

			unMarshalErr := json.Unmarshal(b, cmd)

			log.Printf("ReadMessage\n")

			if unMarshalErr != nil {
				log.Printf("unmarshal error:%s\n", unMarshalErr.Error())
			} else {

				log.Printf("res cmd:%#v\n", cmd)

				deData, _ := base64.StdEncoding.DecodeString(cmd.Data)
				log.Printf("cmd data:%s\n", string(deData))

				switch cmd.Command {
				case "Login":

					log.Printf("Send command :Heartbeat\n")
					c.WriteMessage(gorillaWebsocket.TextMessage, createTextMessage("Heartbeat", createHeartbeatData()))
				case "Heartbeat":
					log.Printf("Send command :GameInfo\n")
					c.WriteMessage(gorillaWebsocket.TextMessage, createTextMessage("GameInfo", createGameInfoData()))

				case "GameInfo":
					log.Printf("Send command :UserInfo\n")
					c.WriteMessage(gorillaWebsocket.TextMessage, createTextMessage("UserInfo", createUserInfoData()))

				case "UserInfo":
					log.Printf("Send command :Bet\n")
					c.WriteMessage(gorillaWebsocket.TextMessage, createTextMessage("Bet", createBetData()))

				case "Bet":
					log.Printf("Send command : Logout\n")
					c.WriteMessage(gorillaWebsocket.TextMessage, createTextMessage("Logout", createLogoutData()))

				case "Logout":

					done <- struct{}{}

				}
			}

		}
	}()

	//ticker
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
		case <-interrupt:

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(gorillaWebsocket.CloseMessage, gorillaWebsocket.FormatCloseMessage(gorillaWebsocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}

func GetSessionID(c *gorillaWebsocket.Conn) {

	requestURL := "http://api.3ag.co/checkOrCreateGameAccount"

	pa := &struct {
		UserName string `json:"userName"`
	}{

		UserName: "testUser0",
	}
	js, _ := json.Marshal(pa)

	jsStr := string(js)
	res, err := net.HTTPPost(requestURL, "accessToken=NDkwODVhMjItZGY4NC00ZTBhLWI0MTktYjFjY2YyNWIyZjY2&params="+jsStr)

	if err != nil {
		return
	}

	resData := &struct {
		SessionID   string `json:"sessionId"`
		GamePageUrl string `json:"gamePageUrl"`
	}{}

	json.Unmarshal([]byte(res), resData)

	//userKill(resData.SessionID)
	//return resData.SessionID

	c.WriteMessage(gorillaWebsocket.TextMessage, createTextMessage("Login", createLoginData(resData.SessionID)))

}

func createLoginData(sessionID string) string {

	l := LoginCmdData{
		GameId:    "slot0",
		SessionId: sessionID,
		Lang:      "zh-hans",
	}
	d := []LoginCmdData{l}

	b, err := json.Marshal(&d)
	if err != nil {
		log.Printf("marshal LoginCmdData error:%s\n", err.Error())
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b)
}

func createHeartbeatData() string {

	l := HeartbeatCmdData{
		SessionId: "sessinID",
	}
	d := []HeartbeatCmdData{l}

	b, err := json.Marshal(&d)
	if err != nil {
		log.Printf("marshal HeartbeatCmdData error:%s\n", err.Error())
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b)
}

func createUserInfoData() string {

	l := UserInfoCmdData{
		SessionId: "sessinID",
	}
	d := []UserInfoCmdData{l}

	b, err := json.Marshal(&d)
	if err != nil {
		log.Printf("marshal UserInfoCmdData error:%s\n", err.Error())
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b)
}

func createGameInfoData() string {

	l := GameInfoCmdData{
		GameId: "slot0",
	}
	d := []GameInfoCmdData{l}

	b, err := json.Marshal(&d)
	if err != nil {
		log.Printf("marshal GameInfoCmdData error:%s\n", err.Error())
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b)
}

func createBetData() string {

	l := BetCmdData{
		Mode:      0,
		BetCredit: 10.13,
	}
	d := []BetCmdData{l}

	b, err := json.Marshal(&d)
	if err != nil {
		log.Printf("marshal createBetData error:%s\n", err.Error())
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b)
}

func createLogoutData() string {

	l := LogoutCmdData{
		Type: 0,
	}
	d := []LogoutCmdData{l}

	b, err := json.Marshal(&d)
	if err != nil {
		log.Printf("marshal LogoutCmdData error:%s\n", err.Error())
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b)
}

func createTextMessage(cmd string, dataStr string) []byte {
	c, _ := createCommand(0, 1, cmd, dataStr)

	ce, _ := json.Marshal(c)

	str := base64.StdEncoding.EncodeToString(ce)

	return []byte(str)

}

func createCompletePacket() []byte {
	d, _ := createLoginCommandData("gameId", "sessionID", "lang")
	c, _ := createCommand(0, 1, "Login", d)
	ce, _ := json.Marshal(c)
	return composePacket(0, ce)

}
func composePacket(packetType int, packetBody []byte) []byte {
	pl := uint32(len(packetBody))

	headBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(headBytes, pl)

	typeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, uint16(packetType))

	res := make([]byte, pl+4+1)
	copy(res[0:4], headBytes[:])
	copy(res[4:5], typeBytes[1:])
	copy(res[5:], packetBody[:])

	return res

}

func createCommand(code int, step int, cmdName string, data string) (*nex.Command, error) {

	cmd := &nex.Command{
		Code:    code,
		Step:    step,
		Command: cmdName,
		Data:    data,
	}
	return cmd, nil
}

func createLoginCommandData(gameID, sessinID, lang string) (string, error) {

	l := LoginCmdData{
		GameId:    gameID,
		SessionId: sessinID,
		Lang:      lang,
	}
	d := []LoginCmdData{l}

	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b), nil
}

type LoginCmdData struct {
	GameId    string
	SessionId string
	Lang      string
}
type LoginCmdResData struct {
}

type LogoutCmdData struct {
	Type int
}

type GameInfoCmdData struct {
	GameId string
}

type GameInfoResData struct {
	Name       string
	CategoryId int
	Wheel      int
	Row        int
	LineMax    int
	BetChip    []float32
	PayLineNum int
}

type UserInfoCmdData struct {
	SessionId string
}

type UserInfoResData struct {
	Name          string
	Credit        float32
	PresentCredit float32
	IsTest        int
	LastLogin     string
}

type BetCmdData struct {
	Mode      int
	BetCredit float32
}

type BetResData struct {
	ShowCredit    int
	ResultCredit  float32
	BalanceCredit float32
	Symbol        [][]int
	LinkLine      [][]int
	FreeSpinTimes int
}

type FreeSpinResultResData struct {
	ShowCredit int
}

type HeartbeatCmdData struct {
	SessionId string
}
type HeartbeatResData struct {
}

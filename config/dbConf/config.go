package dbConf

const (
	DBPathLimitations = "/limitations"
	DBPathDealers     = "/dealers"
	DBPathRooms       = "/rooms"
	DBPathHalls       = "/halls"
	DBPathBanners     = "/banners"
	DBPathBroadcasts  = "/broadcasts"

	CodeSuccess = 0

	DefaultJSON="[]"
)

type GotData struct {
	Code    int         `json:"code"`
	Count   int         `json:"count"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//hall
type ResponseHallData struct {
	Code    int         `json:"code"`
	Count   int         `json:"count"`
	Message string      `json:"message"`
	Data    []*HallData `json:"data"`
}

type HallData struct {
	HallID     uint   `json:"hallID"`
	Name       string `json:"name"`
	Active     uint   `json:"active"`
	CreateDate string `json:"createDate"`
}

//room
type ResponseRoomData struct {
	Code    int         `json:"code"`
	Count   int         `json:"count"`
	Message string      `json:"message"`
	Data    []*RoomData `json:"data"`
}

type RoomData struct {
	RoomID        uint   `json:"roomID"`
	HallID        uint   `json:"hallID"`
	Name          string `json:"name"`
	RoomType      uint   `json:"roomType"`
	Active        uint   `json:"active"`
	HLSURL        string `json:"hlsURL"`
	Boot          uint   `json:"boot"`
	RoundID       uint64 `json:"round"`
	Status        int    `json:"status"`
	BetCountdown  uint   `json:"betCountdown"`
	DealerID      uint   `json:"dealerID"`
	LimitationID  uint   `json:"limitationID"`
	HistoryResult string `json:"historyResult"`
	CreateDate    string `json:"createDate"`
}
type Status struct {
	Status int `json:"status"`
}
type RoomNewRoundData struct {
	Boot    int   `json:"boot"`
	RoundID int64 `json:"round"`
	Status  int   `json:"status"`
}

//limitation
type ResponseLimitation struct {
	Code    int           `json:"code"`
	Count   int           `json:"count"`
	Message string        `json:"message"`
	Data    []*Limitation `json:"data"`
}

type Limitation struct {
	LimitationID int    `json:"limitationID"`
	Limitation   string `json:"limitation"`
}

//dealer
type DealerData struct {
	Code    int       `json:"code"`
	Count   int       `json:"count"`
	Message string    `json:"message"`
	Data    []*Dealer `json:"data"`
}
type Dealer struct {
	DealerID    int    `json:"dealerID"`
	Name        string `json:"name"`
	Account     string `json:"account"`
	Active      uint   `json:"active"`
	PortraitURL string `json:"portraitURL"`
	CreateDate  string `json:"createDate"`
}

//accessToken
type AccessTokenData struct {
	Code    int            `json:"code"`
	Count   int            `json:"count"`
	Message string         `json:"message"`
	Data    []*AccessToken `json:"data"`
}

type AccessToken struct {
	UserID            int64   `json:"userID"`
	Account           string  `json:"account"`
	Credit            float32 `json:"credit"`
	Name              string  `json:"name"`
	PartnerID         int64   `json:"partnerID"`
	Active            int     `json:"active"`
	AccessTokenExpire string  `json:"accessTokenExpire"`
}

//banner
type BannerData struct {
	Code    int       `json:"code"`
	Count   int       `json:"count"`
	Message string    `json:"message"`
	Data    []*Banner `json:"data"`
}
type Banner struct {
	BannerID    uint   `json:"bannerID"`
	PicURL      string `json:"picURL"`
	LinkURL     string `json:"linkURL"`
	Description string `json:"description"`
	Platform    uint   `json:"platform"`
	Active      uint   `json:"active"`
	CreateDate  string `json:"createDate"`
}

//banner
type BroadcastData struct {
	Code    int          `json:"code"`
	Count   int          `json:"count"`
	Message string       `json:"message"`
	Data    []*Broadcast `json:"data"`
}
type Broadcast struct {
	BroadcastID uint   `json:"broadcastID"`
	Content     string `json:"content"`
	Internal    int    `json:"internal"`
	RepeatTimes int    `json:"repeatTimes"`
	Active      uint   `json:"active"`
	CreateDate  string `json:"createDate"`
}

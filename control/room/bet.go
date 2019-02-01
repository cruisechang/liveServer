package room

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/cruisechang/liveServer/config/roomConf"
	"github.com/cruisechang/nex/entity"
)

const (
	type6DiceLen   = 6
	type6PairLen   = 6
	type6TripleLen = 7
	type6SumLen    = 14
	type6PaigowLen = 15

	type7OneLen    = 37
	type7TwoLen    = 60
	type7ThreeLen  = 14
	type7FourLen   = 23
	type7SixLen    = 11
	type7ColumnLen = 3
	type7DozenLen  = 3
)

// bet
func (r *Controller) InitRoomBet(room entity.Room) {
	//map[userID]*roomConf.BetType1Data
	switch room.Type() {
	case r.conf.RoomType0():
		room.SetInterfaceVariable(varBet, make(map[int]*roomConf.BetType0Data))

	case r.conf.RoomType1():
		room.SetInterfaceVariable(varBet, make(map[int]*roomConf.BetType1Data))

	case r.conf.RoomType2():
		room.SetInterfaceVariable(varBet, make(map[int]*roomConf.BetType2Data))

	case r.conf.RoomType6():
		room.SetInterfaceVariable(varBet, make(map[int]*roomConf.BetType6Data))

	case r.conf.RoomType7():
		room.SetInterfaceVariable(varBet, make(map[int]*roomConf.BetType7Data))
	}
}

func (r *Controller) AddBet(room entity.Room, userID int, cmdBetData interface{}) error {
	switch room.Type() {
	case r.conf.RoomType0():
		return r.addBetType0(room, userID, cmdBetData)

	case r.conf.RoomType1():
		return r.addBetType1(room, userID, cmdBetData)

	case r.conf.RoomType2():
		return r.addBetType2(room, userID, cmdBetData)

	case r.conf.RoomType6():
		return r.addBetType6(room, userID, cmdBetData)

	case r.conf.RoomType7():
		return r.addBetType7(room, userID, cmdBetData)

	default:
		return errors.New("AddBet room type not found")
	}
}

//type0
func (r *Controller) addBetType0(room entity.Room, userID int, betData interface{}) error {

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	m, _ := dt.(map[int]*roomConf.BetType0Data)

	if bd, ok := betData.(*roomConf.BetType0Data); ok {

		if d, ok := m[userID]; ok {
			d = r.addBetType0Real(d, bd)
			m[userID] = d
		} else {
			m[userID] = bd
		}
	} else {
		return errors.New(fmt.Sprintf("addBetType0 passed cmdBetData error =%s", reflect.TypeOf(betData).String()))
	}

	return nil
}
func (r *Controller) addBetType0Real(data *roomConf.BetType0Data, betData *roomConf.BetType0Data) *roomConf.BetType0Data {

	return &roomConf.BetType0Data{
		Banker:      data.Banker + betData.Banker,
		Player:      data.Player + betData.Player,
		Tie:         data.Tie + betData.Tie,
		BankerPair:  data.BankerPair + betData.BankerPair,
		PlayerPair:  data.PlayerPair + betData.PlayerPair,
		Big:         data.Big + betData.Big,
		Small:       data.Small + betData.Small,
		AnyPair:     data.AnyPair + betData.AnyPair,
		PerfectPair: data.PerfectPair + betData.PerfectPair,
		SuperSix:    data.SuperSix + betData.SuperSix,
	}
}

//type1
func (r *Controller) addBetType1(room entity.Room, userID int, cmdBetData interface{}) error {

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	m, _ := dt.(map[int]*roomConf.BetType1Data)

	if bd, ok := cmdBetData.(*roomConf.BetType1Data); ok {

		if d, ok := m[userID]; ok {
			d = r.addBetType1Real(d, bd)
			m[userID] = d
		} else {
			m[userID] = bd
		}
	} else {
		return errors.New(fmt.Sprintf("addBetType1 passed cmdBetData error =%s", reflect.TypeOf(cmdBetData).String()))
	}

	return nil
}

func (r *Controller) addBetType1Real(data *roomConf.BetType1Data, cmdData *roomConf.BetType1Data) *roomConf.BetType1Data {

	return &roomConf.BetType1Data{
		Dragon:      data.Dragon + cmdData.Dragon,
		Tiger:       data.Tiger + cmdData.Tiger,
		Tie:         data.Tie + cmdData.Tie,
		DragonOdd:   data.DragonOdd + cmdData.DragonOdd,
		DragonEven:  data.DragonEven + cmdData.DragonEven,
		DragonRed:   data.DragonRed + cmdData.DragonRed,
		DragonBlack: data.DragonBlack + cmdData.DragonBlack,
		TigerOdd:    data.TigerOdd + cmdData.Dragon,
		TigerEven:   data.TigerEven + cmdData.TigerEven,
		TigerRed:    data.TigerRed + cmdData.TigerRed,
		TigerBlack:  data.TigerBlack + cmdData.TigerBlack,
	}
}

//type2
func (r *Controller) addBetType2(room entity.Room, userID int, cmdBetData interface{}) error {

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	m, _ := dt.(map[int]*roomConf.BetType2Data)

	if bd, ok := cmdBetData.(*roomConf.BetType2Data); ok {

		if d, ok := m[userID]; ok {
			d = r.addBetType2Real(d, bd)
			m[userID] = d
		} else {
			m[userID] = bd
		}
	} else {
		return errors.New(fmt.Sprintf("addBetType2 passed cmdBetData error =%s", reflect.TypeOf(cmdBetData).String()))
	}

	return nil
}
func (r *Controller) addBetType2Real(data *roomConf.BetType2Data, cmdData *roomConf.BetType2Data) *roomConf.BetType2Data {

	return &roomConf.BetType2Data{
		Owner1: roomConf.BetType2DataOwner{WinEqual: data.Owner1.WinEqual + cmdData.Owner1.WinEqual, WinTimes: data.Owner1.WinTimes + cmdData.Owner1.WinTimes, LoseEqual: data.Owner1.LoseEqual + cmdData.Owner1.LoseEqual, LoseTimes: data.Owner1.LoseTimes + cmdData.Owner1.LoseTimes},
		Owner2: roomConf.BetType2DataOwner{WinEqual: data.Owner2.WinEqual + cmdData.Owner2.WinEqual, WinTimes: data.Owner2.WinTimes + cmdData.Owner2.WinTimes, LoseEqual: data.Owner2.LoseEqual + cmdData.Owner2.LoseEqual, LoseTimes: data.Owner2.LoseTimes + cmdData.Owner2.LoseTimes},
		Owner3: roomConf.BetType2DataOwner{WinEqual: data.Owner3.WinEqual + cmdData.Owner3.WinEqual, WinTimes: data.Owner3.WinTimes + cmdData.Owner3.WinTimes, LoseEqual: data.Owner3.LoseEqual + cmdData.Owner3.LoseEqual, LoseTimes: data.Owner3.LoseTimes + cmdData.Owner3.LoseTimes},
	}
}

//CountUserBetSumType2 return bet sum
//func (r *Controller) CountUserBetSumType2(userBet *roomConf.BetType2Data) (sum int, err error) {
//	owner1 := userBet.Owner1
//	owner2 := userBet.Owner2
//	owner3 := userBet.Owner3
//	sum += owner1.LoseEqual + owner1.WinEqual
//	sum += owner2.LoseEqual + owner2.WinEqual
//	sum += owner3.LoseEqual + owner3.WinEqual
//
//	//翻倍要夠五倍才行
//	sum += (owner1.LoseTimes + owner1.WinTimes +
//		owner2.LoseTimes + owner2.WinTimes +
//		owner3.LoseTimes + owner3.WinTimes) * 5
//	return sum, nil
//}

//type6
func (r *Controller) addBetType6(room entity.Room, userID int, betData interface{}) error {

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	m, _ := dt.(map[int]*roomConf.BetType6Data)

	if bd, ok := betData.(*roomConf.BetType6Data); ok {

		if d, ok := m[userID]; ok {
			d = r.addBetType6Real(d, bd)
			m[userID] = d
		} else {
			m[userID] = bd
		}
	} else {
		return errors.New(fmt.Sprintf("addBetType6 passed cmdBetData error =%s", reflect.TypeOf(betData).String()))
	}

	return nil
}

func (r *Controller) addBetType6Real(data *roomConf.BetType6Data, cmdData *roomConf.BetType6Data) *roomConf.BetType6Data {

	res := &roomConf.BetType6Data{
		Small: data.Small + cmdData.Small,
		Big:   data.Big + cmdData.Big,
		Odd:   data.Odd + cmdData.Odd,
		Even:  data.Even + cmdData.Even,
	}

	//dice and pair
	for i, v := range data.Dice {
		res.Dice[i] = v + cmdData.Dice[i]
		res.Pair[i] = v + cmdData.Pair[i]
	}
	//triple
	for i, v := range data.Triple {
		res.Triple[i] = v + cmdData.Triple[i]
	}
	//sum
	for i, v := range data.Sum {
		res.Sum[i] = v + cmdData.Sum[i]
	}

	//paigow
	for i, v := range data.Paigow {
		res.Paigow[i] = v + cmdData.Paigow[i]
	}

	return res
}

//type7
func (r *Controller) addBetType7(room entity.Room, userID int, betData interface{}) error {

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	m, _ := dt.(map[int]*roomConf.BetType7Data)

	if bd, ok := betData.(*roomConf.BetType7Data); ok {

		if d, ok := m[userID]; ok {
			d = r.addBetType7Real(d, bd)
			m[userID] = d
		} else {
			m[userID] = bd
		}
	} else {
		return errors.New(fmt.Sprintf("addBetType7 passed cmdBetData error =%s", reflect.TypeOf(betData).String()))
	}

	return nil
}
func (r *Controller) addBetType7Real(data *roomConf.BetType7Data, cmdData *roomConf.BetType7Data) *roomConf.BetType7Data {

	res := &roomConf.BetType7Data{
		Small: data.Small + cmdData.Small,
		Big:   data.Big + cmdData.Big,
		Odd:   data.Odd + cmdData.Odd,
		Even:  data.Even + cmdData.Even,
		Red:   data.Red + cmdData.Red,
		Black: data.Black + cmdData.Black,
	}

	for i, v := range data.One {
		res.One[i] = v + cmdData.One[i]
	}

	for i, v := range data.Two {
		res.Two[i] = v + cmdData.Two[i]
	}

	for i, v := range data.Three {
		res.Three[i] = v + cmdData.Three[i]
	}

	for i, v := range data.Four {
		res.Four[i] = v + cmdData.Four[i]
	}

	for i, v := range data.Six {
		res.Six[i] = v + cmdData.Six[i]
	}

	for i, v := range data.Column {
		res.Column[i] = v + cmdData.Column[i]
	}

	for i, v := range data.Dozen {
		res.Dozen[i] = v + cmdData.Dozen[i]
	}

	return res
}

//CountUserBetSumType7 return bet sum
//func (r *Controller) CountUserBetSumType7(userBet *roomConf.BetType7Data) (sum int, err error) {
//
//	for _, v := range userBet.One {
//		sum += v
//	}
//	for _, v := range userBet.Two {
//		sum += v
//	}
//	for _, v := range userBet.Three {
//		sum += v
//	}
//	for _, v := range userBet.Four {
//		sum += v
//	}
//	for _, v := range userBet.Six {
//		sum += v
//	}
//	for _, v := range userBet.Column {
//		sum += v
//	}
//	for _, v := range userBet.Dozen {
//		sum += v
//	}
//
//	sum += userBet.Big
//	sum += userBet.Small
//	sum += userBet.Odd
//	sum += userBet.Even
//	sum += userBet.Red
//	sum += userBet.Black
//
//	return sum, nil
//}

//GetBet returns interface{} of room bet data
func (r *Controller) GetBet(room entity.Room) (interface{}, error) {
	return room.GetInterfaceVariable(varBet)
}

//GetUserBetType0 returns room bet typ (*roomConf.BetType0Data) of target room
func (r *Controller) GetUserBetType0(room entity.Room, userID int) (*roomConf.BetType0Data, error) {

	if room.Type() != r.conf.RoomType0() {
		return nil, errors.New(fmt.Sprintf("Room type error want %d, got %d", r.conf.RoomType0(), room.Type()))
	}

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if m, ok := dt.(map[int]*roomConf.BetType0Data); ok {

		if d, ok := m[userID]; ok {
			return d, nil
		}
		//沒資料
		return r.GetEmptypBetType0(), nil

	}
	return nil, errors.New("get user bet error")
}

//GetUserBetType1 returns room bet typ (*roomConf.BetType1Data) of target room
func (r *Controller) GetUserBetType1(room entity.Room, userID int) (*roomConf.BetType1Data, error) {

	if room.Type() != r.conf.RoomType1() {
		return nil, errors.New(fmt.Sprintf("Room type error want %d, got %d", r.conf.RoomType1(), room.Type()))
	}

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if m, ok := dt.(map[int]*roomConf.BetType1Data); ok {

		if d, ok := m[userID]; ok {
			return d, nil
		}
		//沒資料
		return r.GetEmptypBetType1(), nil
	}
	return nil, errors.New("get user bet error")
}

//GetUserBetType2 returns room bet typ (*roomConf.BetType2Data) of target room
func (r *Controller) GetUserBetType2(room entity.Room, userID int) (*roomConf.BetType2Data, error) {

	if room.Type() != r.conf.RoomType2() {
		return nil, errors.New(fmt.Sprintf("Room type error want %d, got %d", r.conf.RoomType2(), room.Type()))
	}

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if m, ok := dt.(map[int]*roomConf.BetType2Data); ok {

		if d, ok := m[userID]; ok {
			return d, nil
		}
		//沒資料
		return r.GetEmptypBetType2(), nil
	}
	return nil, errors.New("get user bet error")
}

//GetUserBetType6 returns room bet typ (*roomConf.BetType6Data) of target room
func (r *Controller) GetUserBetType6(room entity.Room, userID int) (*roomConf.BetType6Data, error) {

	if room.Type() != r.conf.RoomType6() {
		return nil, errors.New(fmt.Sprintf("Room type error want %d, got %d", r.conf.RoomType6(), room.Type()))
	}

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if m, ok := dt.(map[int]*roomConf.BetType6Data); ok {

		if d, ok := m[userID]; ok {
			return d, nil
		}
		//沒資料
		return r.GetEmptypBetType6(), nil
	}
	return nil, errors.New("get user bet error")
}

//GetUserBetSum returns target user's bet sum of the target room
func (r *Controller) GetUserBetType7(room entity.Room, userID int) (*roomConf.BetType7Data, error) {

	if room.Type() != r.conf.RoomType7() {
		return nil, errors.New(fmt.Sprintf("Room type error want %d, got %d", r.conf.RoomType7(), room.Type()))
	}

	dt, err := room.GetInterfaceVariable(varBet)
	if err != nil {
		return nil, err
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if m, ok := dt.(map[int]*roomConf.BetType7Data); ok {
		if d, ok := m[userID]; ok {
			return d, nil
		}
		//沒資料
		return r.GetEmptypBetType7(), nil
	}

	return nil, errors.New("get user bet error")
}

//GetUserBetSum returns target user's bet sum of the target room
//func (r *Controller) CountUserBetSum(room entity.Room, userID int) (int, error) {
//
//	dt, err := room.GetInterfaceVariable(varBet)
//	if err != nil {
//		return -1, err
//	}
//
//	switch room.Type() {
//	case r.conf.RoomType0():
//	case r.conf.RoomType1():
//	case r.conf.RoomType2():
//		if m, ok := dt.(map[int]*roomConf.BetType2Data); ok {
//			if d, ok := m[userID]; ok {
//				return r.CountUserBetSumType2(d)
//			}
//			return 0, nil //之前沒有
//		}
//	case r.conf.RoomType6():
//	case r.conf.RoomType7():
//		if m, ok := dt.(map[int]*roomConf.BetType7Data); ok {
//			if d, ok := m[userID]; ok {
//				return r.CountUserBetSumType7(d)
//			}
//			return 0, nil //之前沒有
//		}
//
//	}
//	return -1, errors.New("GetUserBetSum error")
//}

//emptyBet
func (r *Controller) GetEmptypBetType0() *roomConf.BetType0Data {
	return &roomConf.BetType0Data{}
}
func (r *Controller) GetEmptypBetType1() *roomConf.BetType1Data {
	return &roomConf.BetType1Data{}
}
func (r *Controller) GetEmptypBetType2() *roomConf.BetType2Data {
	return &roomConf.BetType2Data{
		Owner1: roomConf.BetType2DataOwner{},
		Owner2: roomConf.BetType2DataOwner{},
		Owner3: roomConf.BetType2DataOwner{},
	}
}
func (r *Controller) GetEmptypBetType6() *roomConf.BetType6Data {

	return &roomConf.BetType6Data{
		Small:  0,
		Big:    0,
		Odd:    0,
		Even:   0,
		Dice:   make([]int, type6DiceLen),
		Pair:   make([]int, type6PairLen),
		Triple: make([]int, type6TripleLen),
		Sum:    make([]int, type6SumLen),
		Paigow: make([]int, type6PaigowLen),
	}
}
func (r *Controller) GetEmptypBetType7() *roomConf.BetType7Data {
	return &roomConf.BetType7Data{
		Big:    0,
		Small:  0,
		Odd:    0,
		Even:   0,
		Red:    0,
		Black:  0,
		One:    make([]int, type7OneLen),
		Two:    make([]int, type7TwoLen),
		Three:  make([]int, type7ThreeLen),
		Four:   make([]int, type7FourLen),
		Six:    make([]int, type7SixLen),
		Column: make([]int, type7ColumnLen),
		Dozen:  make([]int, type7DozenLen),
	}
}

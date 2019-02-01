package main

import (
	"reflect"
	"testing"
	"time"
)

func TestUnixEpoch(t *testing.T) {
	nano := time.Now().UnixNano()
	milli := nano / 1000000
	milli2 := int(nano / 1000000)
	macro := nano / 100000

	t.Logf("unix ecpoch=%v", time.Now().Unix())
	t.Logf("unix milli, ecpoch=%v", milli)
	t.Logf("unix milli, ecpoch=%v", milli2)
	t.Logf("unix macro, ecpoch=%v", macro)
	t.Logf("unix nano, ecpoch=%v", nano)

	t.Logf("type %s ", reflect.TypeOf(milli).String())

}

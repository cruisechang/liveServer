package config

import (
	"math"
	"testing"

	"github.com/shopspring/decimal"
)

var (
	testConf Configurer
)

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

func TestLoadConfigData(t *testing.T) {

	c, err := NewConfigurer("config.json")
	if err != nil {
		t.Fatalf("TestLoadConfigData error:%v", err)
	}

	testConf = c

	var ch float64 = 0.029
	ch0 := decimal.NewFromFloat(ch)
	t.Logf("decimal test decimal.NewFromFloat:%s", ch0.RoundBank(2).String())

}

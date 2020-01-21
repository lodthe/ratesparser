package exchangers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/lodthe/ratesparser/models"
)

const (
	BINANCE_API_URL = "https://api.binance.com/api/v3/"
)

type Binance struct{Counter int}

type BinanceTickerPriceResponse struct {
	Symbol string
	Price  string
}

func (exchanger *Binance) Name() string {
	return "Binance"
}

func (exchanger *Binance) GetRatePrice(pair models.Pair) (float64, error) {
	queryString := BINANCE_API_URL + "ticker/price?symbol=" + pair.SpendCurrency + pair.BuyCurrency
	response, err := http.Get(queryString)
	if err != nil {
		return 0, fmt.Errorf("cannot get %s rate for %v", exchanger.Name(), pair)
	}
	defer response.Body.Close()

	var binanceResponse BinanceTickerPriceResponse
	responseBody, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(responseBody, &binanceResponse)
	result, err := strconv.ParseFloat(binanceResponse.Price, 64)
	if err != nil   {
		return 0, fmt.Errorf("cannot parse price for %v from %s response", pair, exchanger.Name())
	}
	return result, nil
}

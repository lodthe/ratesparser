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
	EXMO_API_URL = "https://api.exmo.com/v1/ticker/"
)

type Exmo struct{}

type ExmoTickerPriceResponse struct {
	Buy_price  string
}

func (exchanger *Exmo) Name() string {
	return "Exmo"
}

func (exchanger *Exmo) GetRatePrice(pair models.Pair) (float64, error) {
	response, err := http.Get(EXMO_API_URL)
	if err != nil {
		return 0, fmt.Errorf("cannot get %s rate for %v", exchanger.Name(), pair)
	}
	defer response.Body.Close()

	var binanceResponse map[string]ExmoTickerPriceResponse
	responseBody, _ := ioutil.ReadAll(response.Body)
	_ = json.Unmarshal(responseBody, &binanceResponse)

	symbol := pair.SpendCurrency + "_" + pair.BuyCurrency
	ticker, found := binanceResponse[symbol]
	if found == false {
		return 0, fmt.Errorf("there is no ticker for %s on %s", pair, exchanger.Name())
	}
	result, err := strconv.ParseFloat(ticker.Buy_price, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse price for %v from %s response", pair, exchanger.Name())
	}
	return result, nil
}

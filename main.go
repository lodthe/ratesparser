package main

import (
	"encoding/json"
	"fmt"
	"github.com/lodthe/ratesparser/controllers"
	"github.com/lodthe/ratesparser/exchangers"
	"github.com/lodthe/ratesparser/models"
	"io"
	"log"
	"net/http"
)

var (
	EXCHANGERS = []exchangers.Exchanger{
		&exchangers.Binance{},
		&exchangers.Exmo{},
	}
	PAIRS = []models.Pair{
		{SpendCurrency:"BTC", BuyCurrency:"USDT"},
		{SpendCurrency:"ETH", BuyCurrency:"BTC"},
	}
)

var controller = &controllers.RatesController{Platforms: EXCHANGERS, Pairs:PAIRS}

func GetRates(w http.ResponseWriter, r *http.Request) {
	rates := make(chan *models.Rate)
	response := make([]*models.Rate, 0)
	go controller.GetAllRates(rates)

	for rate := range rates {
		response = append(response, rate)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	if responseJson, err := json.Marshal(response); err == nil {
		io.WriteString(w, string(responseJson))
	} else {
		fmt.Println(err)
		io.WriteString(w, "Cannot parse response")
	}
}

func main() {
	controller.Init()
	go controller.Run(300)
	http.HandleFunc("/get_rates", GetRates)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

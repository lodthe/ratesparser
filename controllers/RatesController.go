package controllers

import (
	"fmt"
	"sync"
	"time"

	"github.com/lodthe/ratesparser/exchangers"
	"github.com/lodthe/ratesparser/models"
)

type RatesController struct {
	Platforms []exchangers.Exchanger
	Pairs     []models.Pair
	sync.Mutex
	ratesCache map[string]map[models.Pair]*models.Rate // [platform][pair] keeps Rate
}

func (controller *RatesController) Init() {
	controller.ratesCache = make(map[string]map[models.Pair]*models.Rate)
}

//Run starts getting relevant rates from exchangers a fixed rate
func (controller *RatesController) Run(delay float64) {
	for _, exchanger := range controller.Platforms {
		for _, pair := range controller.Pairs {
			go func(exchanger exchangers.Exchanger, pair models.Pair) {
				for {
					go func() {
						if price, err := exchanger.GetRatePrice(pair); err == nil {
							fmt.Printf("successfully got rate for %s on %s\n", pair, exchanger.Name())
							controller.UpdateRate(models.Rate{
								Pair:          pair,
								Time:          time.Now(),
								ExchangerName: exchanger.Name(),
								Price:         price,
							})
						} else {
							fmt.Println(err)
						}
					}()

					time.Sleep(time.Second * time.Duration(delay))
				}
			}(exchanger, pair)
		}
	}
}

//getOrCreateExchangerMap creates map with rates by exchanger name
//if it doesn't exist and returns it
func (controller *RatesController) getOrCreateExchangerMap(exchangerName string) map[models.Pair]*models.Rate {
	if exchangerMap, found := controller.ratesCache[exchangerName]; found == false {
		exchangerMap = make(map[models.Pair]*models.Rate)
		controller.ratesCache[exchangerName] = exchangerMap
		return exchangerMap
	} else {
		return exchangerMap
	}
}

//UpdateRate updates information about rate
func (controller *RatesController) UpdateRate(rate models.Rate) {
	controller.Lock()
	defer controller.Unlock()

	controller.getOrCreateExchangerMap(rate.ExchangerName)[rate.Pair] = &rate
}

//GetRate returns pair rate by it's exchanger name
func (controller *RatesController) GetRate(exchangerName string, pair models.Pair) (*models.Rate, error) {
	controller.Lock()
	defer controller.Unlock()

	result, ok := controller.getOrCreateExchangerMap(exchangerName)[pair]
	if ok == false {
		return result, fmt.Errorf("there is no rate for %s on %s", pair, exchangerName)
	}
	return result, nil
}

//GetAllRates collects cached rates and returns them
func (controller *RatesController) GetAllRates(rates chan *models.Rate) {
	wg := sync.WaitGroup{}

	for _, exchanger := range controller.Platforms {
		for _, pair := range controller.Pairs {
			wg.Add(1)

			go func(exchangerName string, pair models.Pair) {
				defer wg.Done()

				if rate, err := controller.GetRate(exchangerName, pair); err == nil {
					rates <- rate
				} else {
					fmt.Println(err)
				}
			}(exchanger.Name(), pair)
		}
	}

	wg.Wait()
	close(rates)
}

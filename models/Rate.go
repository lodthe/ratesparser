package models

import (
	"time"
)

type Rate struct {
	Pair          Pair
	Time          time.Time
	ExchangerName string
	Price         float64
}

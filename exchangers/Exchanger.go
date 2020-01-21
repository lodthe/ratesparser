package exchangers

import "github.com/lodthe/ratesparser/models"

type Exchanger interface {
	Name() string
	GetRatePrice(pair models.Pair) (float64, error)
}

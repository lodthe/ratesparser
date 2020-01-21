package models

import "fmt"

type Pair struct {
	SpendCurrency string
	BuyCurrency   string
}

func (pair Pair) String() string {
	return fmt.Sprintf("Pair(%s/%s)", pair.SpendCurrency, pair.BuyCurrency)
}
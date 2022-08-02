package priceprinter

import (
	"fmt"

	"tickerprice/cmd/fairprice/internal/types"
)

// PricePrinter is a price printer to the standard output device.
type PricePrinter struct{}

// New creates a new initialized instance of PricePrinter.
func New() *PricePrinter {
	return &PricePrinter{}
}

// Print prints the timestamp and price to the standard output device.
func (p *PricePrinter) Print(tickers <-chan types.TickerPrice) {
	for ticker := range tickers {
		fmt.Printf("%d, %v\n", ticker.Time.Unix(), ticker.Price)
	}
}

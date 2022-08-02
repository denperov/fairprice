package mockpricesource

import (
	"context"
	"fmt"
	"time"

	"tickerprice/cmd/fairprice/internal/types"
)

// MockPriceSource is mock price source used for debugging.
type MockPriceSource struct {
	price    float64
	interval time.Duration
}

// New creates a new initialized instance of MockPriceSource.
func New(price float64, interval time.Duration) *MockPriceSource {
	return &MockPriceSource{
		price:    price,
		interval: interval,
	}
}

// SubscribePriceStream subscribes to price updates from the source.
func (d *MockPriceSource) SubscribePriceStream(
	ctx context.Context,
	ticker types.Ticker,
) (<-chan types.TickerPrice, <-chan error) {
	tickerPrices := make(chan types.TickerPrice)
	tickerErrors := make(chan error)

	go func() {
		defer func() {
			close(tickerPrices)
			close(tickerErrors)
		}()

		for {
			tickerPrices <- types.TickerPrice{
				Ticker: ticker,
				Time:   time.Now(),
				Price:  fmt.Sprintf("%6f", d.price),
			}

			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				return
			}
		}
	}()

	return tickerPrices, tickerErrors
}

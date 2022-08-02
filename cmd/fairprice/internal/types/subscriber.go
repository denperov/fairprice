package types

import "context"

//go:generate moq -pkg types_test -out mocks_test.go . PriceStreamSubscriber

type PriceStreamSubscriber interface {
	SubscribePriceStream(context.Context, Ticker) (<-chan TickerPrice, <-chan error)
}

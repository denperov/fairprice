package fairpricesource

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"tickerprice/cmd/fairprice/internal/types"
	"tickerprice/internal/log"
)

//go:generate moq -pkg fairpricesource_test -out mocks_test.go . PriceAlgorithm PriceStorage
//go:generate moq -pkg fairpricesource_test -out mocks_types_test.go ../types PriceStreamSubscriber

// PriceAlgorithm is an algorithm for calculating a fair price based on an array of prices.
type PriceAlgorithm interface {
	// CalculatePrice calculates a fair price based on prices from different sources.
	CalculatePrice(prices map[types.SourceID]float64) (float64, error)
}

// PriceStorage is a storage for prices.
type PriceStorage interface {
	AddPrice(ticker types.Ticker, timeslot types.Timeslot, sourceID types.SourceID, price string)
	GetPrices(ticker types.Ticker, timeslot types.Timeslot) map[types.SourceID]string
	RemovePrices(ticker types.Ticker, timeslot types.Timeslot)
}

// FairPriceSource is the source of the aggregated price from other sources.
type FairPriceSource struct {
	algorithm   PriceAlgorithm
	storage     PriceStorage
	subscribers map[types.SourceID]types.PriceStreamSubscriber
	timeNowFunc func() time.Time
}

// New creates a new initialized instance of FairPriceSource.
func New(
	algorithm PriceAlgorithm,
	storage PriceStorage,
	subscribers map[types.SourceID]types.PriceStreamSubscriber,
	timeNowFunc func() time.Time,
) *FairPriceSource {
	return &FairPriceSource{
		algorithm:   algorithm,
		storage:     storage,
		subscribers: subscribers,
		timeNowFunc: timeNowFunc,
	}
}

// SubscribePriceStream subscribes to price updates from the source.
func (p *FairPriceSource) SubscribePriceStream(
	ctx context.Context,
	ticker types.Ticker,
) (<-chan types.TickerPrice, <-chan error) {
	subscribersWaitGroup := sync.WaitGroup{}

	for sourceID, subscriber := range p.subscribers {
		subscribersWaitGroup.Add(1)

		go func(sourceID types.SourceID, subscriber types.PriceStreamSubscriber) {
			defer subscribersWaitGroup.Done()

			p.runSubscriber(ctx, ticker, sourceID, subscriber)
		}(sourceID, subscriber)
	}

	outTickerPrices := make(chan types.TickerPrice)
	outTickerErrors := make(chan error)

	go func() {
		defer func() {
			subscribersWaitGroup.Wait()

			close(outTickerPrices)
			close(outTickerErrors)
		}()

		p.runPublisher(ctx, ticker, outTickerPrices)
	}()

	return outTickerPrices, outTickerErrors
}

func (p *FairPriceSource) runSubscriber(
	ctx context.Context,
	ticker types.Ticker,
	sourceID types.SourceID,
	subscriber types.PriceStreamSubscriber,
) {
	reconnectWithDelay(ctx, func() {
		tickerPrices, tickerErrors := subscriber.SubscribePriceStream(ctx, ticker)

		for tickerPrice := range tickerPrices {
			timeslot := calculateTimeslot(tickerPrice.Time)

			p.storage.AddPrice(ticker, timeslot, sourceID, tickerPrice.Price)
		}

		// stream can return an error, in that case the channel is closed
		for tickerError := range tickerErrors {
			log.Errorf(ctx, "subscription: %v", tickerError)
		}
	})
}

func (p *FairPriceSource) runPublisher(
	ctx context.Context,
	ticker types.Ticker,
	outTickerPrices chan<- types.TickerPrice,
) {
	p.executeAtTimeslotEnd(ctx, func(timeslot types.Timeslot) {
		stringPrices := p.storage.GetPrices(ticker, timeslot)

		prices := parsePrices(ctx, stringPrices)

		fairPrice, err := p.algorithm.CalculatePrice(prices)
		if err != nil {
			log.Errorf(ctx, "calculate fair price: %v", err)
			return
		}

		fairTickerPrice := types.TickerPrice{
			Ticker: ticker,
			Time:   timeslot.ToTime(),
			Price:  formatPrice(fairPrice),
		}

		select {
		case <-ctx.Done():
			return

		case outTickerPrices <- fairTickerPrice:
			p.storage.RemovePrices(ticker, timeslot)
		}
	})
}

func (p *FairPriceSource) executeAtTimeslotEnd(ctx context.Context, fn func(timeslot types.Timeslot)) {
	currentTimeslot := calculateTimeslot(p.timeNowFunc())

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(time.Second):
			timeslot := calculateTimeslot(p.timeNowFunc())

			// wait for the next timeslot
			if timeslot != currentTimeslot {
				fn(currentTimeslot)

				currentTimeslot = timeslot
			}
		}
	}
}

func parsePrices(ctx context.Context, stringPrices map[types.SourceID]string) map[types.SourceID]float64 {
	prices := make(map[types.SourceID]float64, len(stringPrices))

	for sourceID, stringPrice := range stringPrices {
		price, err := parsePrice(stringPrice)
		if err != nil {
			log.Errorf(ctx, "parse price: %v", err)
			continue
		}

		prices[sourceID] = price
	}

	return prices
}

func parsePrice(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}

	return f, nil
}

func formatPrice(f float64) string {
	return strconv.FormatFloat(f, 'f', 10, 64)
}

func calculateTimeslot(t time.Time) types.Timeslot {
	t = t.UTC()

	year, month, day := t.Date()

	// the start of the time slot is the same as the start of the minute in UTC
	return types.Timeslot(time.Date(year, month, day, t.Hour(), t.Minute(), 0, 0, t.Location()).Unix())
}

func reconnectWithDelay(ctx context.Context, connect func()) {
	// check context cancellation before repeat
	for ctx.Err() == nil {
		connect()

		// avoid the retries storm
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}
	}
}

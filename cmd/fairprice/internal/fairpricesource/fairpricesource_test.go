package fairpricesource_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"tickerprice/cmd/fairprice/internal/fairpricesource"
	"tickerprice/cmd/fairprice/internal/types"
)

func TestFairPriceSource_SubscribePriceStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mockTicker = types.Ticker("ticker_1")

		mockPrice1    = "1.0"
		mockPrice2    = "3.0"
		mockFairPrice = "2.0000000000"

		mockPriceFloat1    = 1.0
		mockPriceFloat2    = 3.0
		mockFairPriceFloat = 2.0

		mockTimeslot = types.Timeslot(60)

		mockTickerPrice1 = types.TickerPrice{
			Ticker: mockTicker,
			Time:   time.Unix(62, 0),
			Price:  mockPrice1,
		}

		mockTickerPrice2 = types.TickerPrice{
			Ticker: mockTicker,
			Time:   time.Unix(63, 0),
			Price:  mockPrice2,
		}

		mockSource1 = &PriceStreamSubscriberMock{
			SubscribePriceStreamFunc: func(
				ctx context.Context,
				ticker types.Ticker,
			) (
				<-chan types.TickerPrice,
				<-chan error,
			) {
				tickers := make(chan types.TickerPrice, 1)
				errors := make(chan error, 1)

				go func() {
					<-ctx.Done()
					close(tickers)
					close(errors)
				}()

				tickers <- mockTickerPrice1

				return tickers, errors
			},
		}

		mockSource2 = &PriceStreamSubscriberMock{
			SubscribePriceStreamFunc: func(
				ctx context.Context,
				ticker types.Ticker,
			) (
				<-chan types.TickerPrice,
				<-chan error,
			) {
				tickers := make(chan types.TickerPrice, 1)
				errors := make(chan error, 1)

				go func() {
					<-ctx.Done()
					close(tickers)
					close(errors)
				}()

				tickers <- mockTickerPrice2

				return tickers, errors
			},
		}

		mockSourceID1 = types.SourceID("source_1")
		mockSourceID2 = types.SourceID("source_2")

		mockStorage = &PriceStorageMock{
			AddPriceFunc: func(
				ticker types.Ticker,
				timeslot types.Timeslot,
				sourceID types.SourceID,
				price string,
			) {
				switch sourceID {
				case mockSourceID1:
					assert.Equal(t, mockTicker, ticker)
					assert.Equal(t, mockTimeslot, timeslot)
					assert.Equal(t, mockPrice1, price)
				case mockSourceID2:
					assert.Equal(t, mockTicker, ticker)
					assert.Equal(t, mockTimeslot, timeslot)
					assert.Equal(t, mockPrice2, price)
				default:
					t.Fail()
				}
			},
			GetPricesFunc: func(
				ticker types.Ticker,
				timeslot types.Timeslot,
			) map[types.SourceID]string {
				return map[types.SourceID]string{
					mockSourceID1: mockPrice1,
					mockSourceID2: mockPrice2,
				}
			},
			RemovePricesFunc: func(
				ticker types.Ticker,
				timeslot types.Timeslot,
			) {
				assert.Equal(t, mockTicker, ticker)
				assert.Equal(t, mockTimeslot, timeslot)
			},
		}

		mockAlgorithm = &PriceAlgorithmMock{
			CalculatePriceFunc: func(prices map[types.SourceID]float64) (float64, error) {
				expectedPrices := map[types.SourceID]float64{
					mockSourceID1: mockPriceFloat1,
					mockSourceID2: mockPriceFloat2,
				}

				assert.Equal(t, expectedPrices, prices)

				return mockFairPriceFloat, nil
			},
		}

		mockSubscribers = map[types.SourceID]types.PriceStreamSubscriber{
			mockSourceID1: mockSource1,
			mockSourceID2: mockSource2,
		}
	)

	mockTimeNow := time.Unix(119, 0)
	go func() {
		time.Sleep(2 * time.Second)

		mockTimeNow = time.Unix(121, 0)
	}()

	mockTimeNowFunc := func() time.Time {
		return mockTimeNow
	}

	fairPriceSource := fairpricesource.New(mockAlgorithm, mockStorage, mockSubscribers, mockTimeNowFunc)

	tickerPrices, tickerErrors := fairPriceSource.SubscribePriceStream(ctx, mockTicker)

	var resultTickerPrices []types.TickerPrice

	for tickerPrice := range tickerPrices {
		resultTickerPrices = append(resultTickerPrices, tickerPrice)

		cancel()
	}

	var resultTickerError []error

	for tickerError := range tickerErrors {
		resultTickerError = append(resultTickerError, tickerError)
	}

	assert.Empty(t, tickerErrors)
	if assert.Equal(t, 1, len(resultTickerPrices)) {
		assert.Equal(t, mockFairPrice, resultTickerPrices[0].Price)
		assert.Equal(t, mockTimeslot, types.Timeslot(resultTickerPrices[0].Time.Unix()))
	}
}

package memstorage

import (
	"tickerprice/cmd/fairprice/internal/types"
)

// MemoryStorage is a thread-safe storage of prices grouped by ticker, timeslot and source.
type MemoryStorage struct {
	tickers *collection[types.Ticker, *collection[types.Timeslot, *collection[types.SourceID, string]]]
}

// New creates a new initialized instance of MemoryStorage.
func New() *MemoryStorage {
	return &MemoryStorage{
		tickers: createTickerCollection(),
	}
}

// AddPrice adds a new or updates an existing price in the store.
func (s *MemoryStorage) AddPrice(
	ticker types.Ticker,
	timeslot types.Timeslot,
	sourceID types.SourceID,
	price string,
) {
	timeslots := s.tickers.GetOrCreate(ticker, createTimeslotCollection)

	sources := timeslots.GetOrCreate(timeslot, createSourceCollection)

	sources.Set(sourceID, price)
}

// GetPrices returns all prices related to ticker and timeslot.
func (s *MemoryStorage) GetPrices(ticker types.Ticker, timeslot types.Timeslot) map[types.SourceID]string {
	timeslots, ok := s.tickers.Get(ticker)
	if !ok {
		return nil
	}

	sources, ok := timeslots.Get(timeslot)
	if !ok {
		return nil
	}

	return sources.Map()
}

// RemovePrices removes all prices related to ticker and timeslot.
func (s *MemoryStorage) RemovePrices(ticker types.Ticker, timeslot types.Timeslot) {
	timeslots, ok := s.tickers.Get(ticker)
	if !ok {
		return
	}

	timeslots.Del(timeslot)
}

func createTickerCollection() *collection[types.Ticker, *collection[types.Timeslot, *collection[types.SourceID, string]]] {
	return newCollection[types.Ticker, *collection[types.Timeslot, *collection[types.SourceID, string]]]()
}

func createTimeslotCollection() *collection[types.Timeslot, *collection[types.SourceID, string]] {
	return newCollection[types.Timeslot, *collection[types.SourceID, string]]()
}

func createSourceCollection() *collection[types.SourceID, string] {
	return newCollection[types.SourceID, string]()
}

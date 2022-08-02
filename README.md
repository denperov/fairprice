# Fair Price
An example of an aggregator of exchange rates from various sources.

## Description
The aggregator receives data from various sources and every minute build a "fair" price for a given pair.
Bars timestamps are solid minute and provided in on-line manner.

## Requirements for sources
- Data from the streams can come with delays, but strictly in increasing time order for each stream.
- Stream can return an error, in that case the channel is closed.

## Interface Modifications
- A context has been added to the interface to notify the price source when the subscription has ended and allow it to gracefully close channels.
- The requirements for channels returned upon subscription have been changed to read-only.

```golang
type PriceStreamSubscriber interface {
	SubscribePriceStream(context.Context, Ticker) (<-chan TickerPrice, <-chan error)
}

```

## Requirements
- Golang 1.18 or above.
- [MOQ](https://github.com/matryer/moq) to generate mock for interfaces in unit-tests.

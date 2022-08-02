package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"tickerprice/cmd/fairprice/internal/averagealgorithm"
	"tickerprice/cmd/fairprice/internal/fairpricesource"
	"tickerprice/cmd/fairprice/internal/memstorage"
	"tickerprice/cmd/fairprice/internal/mockpricesource"
	"tickerprice/cmd/fairprice/internal/priceprinter"
	"tickerprice/cmd/fairprice/internal/types"
	"tickerprice/internal/log"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	priceSourceA := mockpricesource.New(1.1, 1*time.Second)

	priceSourceB := mockpricesource.New(1.2, 2*time.Second)

	priceSourceC := mockpricesource.New(1.6, 3*time.Second)

	subscribers := map[types.SourceID]types.PriceStreamSubscriber{
		"source_a": priceSourceA,
		"source_b": priceSourceB,
		"source_c": priceSourceC,
	}

	algorithm := averagealgorithm.New()

	storage := memstorage.New()

	fairPriceSource := fairpricesource.New(algorithm, storage, subscribers, time.Now)

	printer := priceprinter.New()

	tickers, errs := fairPriceSource.SubscribePriceStream(ctx, types.BTCUSDTicker)

	printer.Print(tickers)

	for err := range errs {
		log.Errorf(ctx, "fair price subscription: %v", err)
	}
}

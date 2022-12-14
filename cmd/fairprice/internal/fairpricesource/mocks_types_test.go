// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package fairpricesource_test

import (
	"context"
	"sync"
	"tickerprice/cmd/fairprice/internal/types"
)

// Ensure, that PriceStreamSubscriberMock does implement types.PriceStreamSubscriber.
// If this is not the case, regenerate this file with moq.
var _ types.PriceStreamSubscriber = &PriceStreamSubscriberMock{}

// PriceStreamSubscriberMock is a mock implementation of types.PriceStreamSubscriber.
//
// 	func TestSomethingThatUsesPriceStreamSubscriber(t *testing.T) {
//
// 		// make and configure a mocked types.PriceStreamSubscriber
// 		mockedPriceStreamSubscriber := &PriceStreamSubscriberMock{
// 			SubscribePriceStreamFunc: func(contextMoqParam context.Context, ticker types.Ticker) (<-chan types.TickerPrice, <-chan error) {
// 				panic("mock out the SubscribePriceStream method")
// 			},
// 		}
//
// 		// use mockedPriceStreamSubscriber in code that requires types.PriceStreamSubscriber
// 		// and then make assertions.
//
// 	}
type PriceStreamSubscriberMock struct {
	// SubscribePriceStreamFunc mocks the SubscribePriceStream method.
	SubscribePriceStreamFunc func(contextMoqParam context.Context, ticker types.Ticker) (<-chan types.TickerPrice, <-chan error)

	// calls tracks calls to the methods.
	calls struct {
		// SubscribePriceStream holds details about calls to the SubscribePriceStream method.
		SubscribePriceStream []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Ticker is the ticker argument value.
			Ticker types.Ticker
		}
	}
	lockSubscribePriceStream sync.RWMutex
}

// SubscribePriceStream calls SubscribePriceStreamFunc.
func (mock *PriceStreamSubscriberMock) SubscribePriceStream(contextMoqParam context.Context, ticker types.Ticker) (<-chan types.TickerPrice, <-chan error) {
	if mock.SubscribePriceStreamFunc == nil {
		panic("PriceStreamSubscriberMock.SubscribePriceStreamFunc: method is nil but PriceStreamSubscriber.SubscribePriceStream was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Ticker          types.Ticker
	}{
		ContextMoqParam: contextMoqParam,
		Ticker:          ticker,
	}
	mock.lockSubscribePriceStream.Lock()
	mock.calls.SubscribePriceStream = append(mock.calls.SubscribePriceStream, callInfo)
	mock.lockSubscribePriceStream.Unlock()
	return mock.SubscribePriceStreamFunc(contextMoqParam, ticker)
}

// SubscribePriceStreamCalls gets all the calls that were made to SubscribePriceStream.
// Check the length with:
//     len(mockedPriceStreamSubscriber.SubscribePriceStreamCalls())
func (mock *PriceStreamSubscriberMock) SubscribePriceStreamCalls() []struct {
	ContextMoqParam context.Context
	Ticker          types.Ticker
} {
	var calls []struct {
		ContextMoqParam context.Context
		Ticker          types.Ticker
	}
	mock.lockSubscribePriceStream.RLock()
	calls = mock.calls.SubscribePriceStream
	mock.lockSubscribePriceStream.RUnlock()
	return calls
}

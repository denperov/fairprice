package averagealgorithm

import (
	"fmt"

	"tickerprice/cmd/fairprice/internal/types"
)

type AverageAlgorithm struct{}

// New creates a new initialized instance of AverageAlgorithm.
func New() *AverageAlgorithm {
	return &AverageAlgorithm{}
}

// CalculatePrice calculates an average price based on prices from different sources.
func (c *AverageAlgorithm) CalculatePrice(prices map[types.SourceID]float64) (float64, error) {
	if len(prices) == 0 {
		return 0, fmt.Errorf("not enough data to calculate an average price")
	}

	var avg float64

	for _, price := range prices {
		avg += price
	}

	return avg / float64(len(prices)), nil
}

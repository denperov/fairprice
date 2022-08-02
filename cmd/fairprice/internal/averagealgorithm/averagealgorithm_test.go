package averagealgorithm_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"tickerprice/cmd/fairprice/internal/averagealgorithm"
	"tickerprice/cmd/fairprice/internal/types"
)

func TestAverageAlgorithm_CalculatePrice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockPrices := map[types.SourceID]float64{
			"a": 1.0,
			"b": 2.0,
			"c": 6.0,
		}

		expectedAveragePrice := 3.0

		algorithm := averagealgorithm.New()

		fairPrice, err := algorithm.CalculatePrice(mockPrices)

		if assert.NoError(t, err) {
			assert.Equal(t, expectedAveragePrice, math.Round(fairPrice))
		}
	})

	t.Run("empty prices", func(t *testing.T) {
		mockPrices := map[types.SourceID]float64{}

		algorithm := averagealgorithm.New()

		_, err := algorithm.CalculatePrice(mockPrices)

		assert.Error(t, err)
	})
}

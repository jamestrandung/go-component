package rounding

import (
	"context"
	"math"
)

// Component is in charge of rounding the calculated
// fare to the hard-coded decimal places.
type Component struct {
	input Input
}

const decimalPlaces = 2

func (c Component) ExecuteSync(ctx context.Context) (any, error) {
	f := c.input.GetRunningFare()

	pow := math.Pow10(decimalPlaces)
	f.Amount = math.Round(f.Amount*pow) / pow

	return nil, nil
}

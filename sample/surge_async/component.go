package surge

import (
	"context"
	"fmt"

	"github.com/jamestrandung/go-component/sample/surge_async/dependencies"
)

// Component is in charge of fetching surge for a particular
// vehicle at a specific location and time.
type Component struct {
	surgeEngine dependencies.ISurgeEngine
	input       Input
}

const fallbackSurge float64 = 1.0

func (c Component) Execute(ctx context.Context) (output, error) {
	surge, err := c.surgeEngine.FetchSurge(
		ctx,
		c.input.GetVehicleTypeID(),
		c.input.GetPickUpLocation(),
		c.input.GetPickUpTime(),
	)

	if err != nil {
		fmt.Printf("error fetching surge: %v \n", err.Error())

		surge = fallbackSurge
	}

	return output{
		surge: surge,
	}, nil
}

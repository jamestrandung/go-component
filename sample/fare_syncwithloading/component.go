package fare

import (
	"context"
	"fmt"

	"github.com/jamestrandung/go-component"
	"github.com/jamestrandung/go-component/sample/fare_syncwithloading/dependencies"
)

// Component is in charge of calculating the fare for
// travelling from A to B on a particular vehicle.
type Component struct {
	configStore dependencies.IConfigStore
	input       Input
}

func (c Component) Load(ctx context.Context) (dependencies.Configs, error) {
	return c.configStore.FetchConfigs(ctx, c.input.GetVehicleTypeID())
}

func (c Component) ExecuteSync(ctx context.Context, data component.LoadData[dependencies.Configs]) (output, error) {
	if data.Err != nil {
		fmt.Printf("error fetching configs for vehicle %v: %v \n", c.input.GetVehicleTypeID(), data.Err.Error())

		return output{}, data.Err
	}

	configs := data.Data

	kmFare := configs.PerKMFare * c.input.GetDistanceInKM()
	minuteFare := configs.PerMinuteFare * c.input.GetDurationInSeconds() / 60

	fareBeforeSurge := configs.StartingFare + kmFare + minuteFare

	c.input.GetRunningFare().Amount = fareBeforeSurge * c.input.GetSurge()

	return output{
		metadata: Metadata{
			AppliedStartingFare:  configs.StartingFare,
			AppliedPerKMFare:     configs.PerKMFare,
			AppliedPerMinuteFare: configs.PerMinuteFare,
		},
	}, nil
}

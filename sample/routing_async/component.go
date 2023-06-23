package routing

import (
	"context"
	"fmt"

	"github.com/jamestrandung/go-component/sample/routing_async/dependencies"
)

// Component is in charge of fetching the distance and duration to
// travel from A to B at a particular time using a specific vehicle.
type Component struct {
	mapService dependencies.IMapService
	input      Input
}

func (c Component) Execute(ctx context.Context) (output, error) {
	travelPlan, err := c.mapService.FetchTravelPlan(
		ctx,
		c.input.GetVehicleTypeID(),
		c.input.GetPickUpLocation(),
		c.input.GetDropOffLocation(),
		c.input.GetPickUpTime(),
	)

	if err != nil {
		fmt.Printf("error fetching travel plan: %v \n", err.Error())

		travelPlan = c.calculateFallbackTravelPlan()
	}

	return output{
		distanceInKM:      travelPlan.DistanceInKM,
		durationInSeconds: travelPlan.DurationInSeconds,
	}, nil
}

const staticSpeedInKMPerHour float64 = 20

func (c Component) calculateFallbackTravelPlan() dependencies.TravelPlan {
	distanceInKM := calculateGEODistance(
		c.input.GetPickUpLocation().Lat,
		c.input.GetPickUpLocation().Lng,
		c.input.GetDropOffLocation().Lat,
		c.input.GetDropOffLocation().Lng,
	)

	durationInSeconds := distanceInKM / staticSpeedInKMPerHour * 60

	return dependencies.TravelPlan{
		DistanceInKM:      distanceInKM,
		DurationInSeconds: durationInSeconds,
	}
}

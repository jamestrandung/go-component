package dependencies

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jamestrandung/go-component/sample/dto"
	"github.com/jamestrandung/go-component/sample/utils"
)

type TravelPlan struct {
	DistanceInKM      float64
	DurationInSeconds float64
}

type IMapService interface {
	FetchTravelPlan(ctx context.Context, vehicleTypeID int64, pickUpLocation, dropOffLocation dto.Location, pickUpTime time.Time) (TravelPlan, error)
}

type MapService struct{}

func NewMapService() *MapService {
	return &MapService{}
}

func (s *MapService) FetchTravelPlan(
	ctx context.Context,
	vehicleTypeID int64,
	pickUpLocation, dropOffLocation dto.Location,
	pickUpTime time.Time,
) (TravelPlan, error) {
	fmt.Printf("fetching travel plan to get from %v to %v using vehicle %v at %v\n", pickUpLocation, dropOffLocation, vehicleTypeID, pickUpTime)

	if utils.FlipCoin() {
		return TravelPlan{}, errors.New("map service is down")
	}

	return TravelPlan{
		DistanceInKM:      7.5,
		DurationInSeconds: 30 * 60,
	}, nil
}

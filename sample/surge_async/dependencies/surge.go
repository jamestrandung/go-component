package dependencies

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jamestrandung/go-component/sample/dto"
	"github.com/jamestrandung/go-component/sample/utils"
)

type ISurgeEngine interface {
	FetchSurge(ctx context.Context, vehicleTypeID int64, pickUpLocation dto.Location, pickUpTime time.Time) (float64, error)
}

type SurgeEngine struct{}

func NewSurgeEngine() *SurgeEngine {
	return &SurgeEngine{}
}

func (e *SurgeEngine) FetchSurge(ctx context.Context, vehicleTypeID int64, pickUpLocation dto.Location, pickUpTime time.Time) (float64, error) {
	fmt.Printf("fetching surge for vehicle %v at location %v on %v\n", vehicleTypeID, pickUpLocation, pickUpTime)

	if utils.FlipCoin() {
		return 0, errors.New("surge engine is down")
	}

	return 1.5, nil
}

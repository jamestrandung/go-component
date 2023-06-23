package routing

import (
	"time"

	"github.com/jamestrandung/go-component/sample/dto"
)

type Input interface {
	GetVehicleTypeID() int64
	GetPickUpLocation() dto.Location
	GetDropOffLocation() dto.Location
	GetPickUpTime() time.Time
}

type output struct {
	distanceInKM      float64
	durationInSeconds float64
}

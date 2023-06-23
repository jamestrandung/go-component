package request

import (
	"time"

	"github.com/jamestrandung/go-component/sample/dto"
)

type FareCalculationRequest struct {
	VehicleTypeID int64
	PickUp        dto.Location
	DropOff       dto.Location
	PickUpTime    time.Time
}

func (r FareCalculationRequest) GetVehicleTypeID() int64 {
	return r.VehicleTypeID
}

func (r FareCalculationRequest) GetPickUpLocation() dto.Location {
	return r.PickUp
}

func (r FareCalculationRequest) GetDropOffLocation() dto.Location {
	return r.DropOff
}

func (r FareCalculationRequest) GetPickUpTime() time.Time {
	return r.PickUpTime
}

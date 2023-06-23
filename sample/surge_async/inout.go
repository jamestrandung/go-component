package surge

import (
	"time"

	"github.com/jamestrandung/go-component/sample/dto"
)

type Input interface {
	GetVehicleTypeID() int64
	GetPickUpLocation() dto.Location
	GetPickUpTime() time.Time
}

type output struct {
	surge float64
}

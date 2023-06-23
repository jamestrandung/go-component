package fare

import "github.com/jamestrandung/go-component/sample/dto"

type Input interface {
	GetVehicleTypeID() int64
	GetSurge() float64
	GetDistanceInKM() float64
	GetDurationInSeconds() float64
	GetRunningFare() *dto.Fare
}

type Metadata struct {
	AppliedStartingFare  float64
	AppliedPerKMFare     float64
	AppliedPerMinuteFare float64
}

type output struct {
	metadata Metadata
}

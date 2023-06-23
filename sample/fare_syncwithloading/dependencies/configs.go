package dependencies

import (
	"context"
	"errors"
	"fmt"

	"github.com/jamestrandung/go-component/sample/utils"
)

type Configs struct {
	StartingFare  float64
	PerKMFare     float64
	PerMinuteFare float64
}

type IConfigStore interface {
	FetchConfigs(ctx context.Context, vehicleTypeID int64) (Configs, error)
}

type ConfigStore struct{}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{}
}

func (s *ConfigStore) FetchConfigs(ctx context.Context, vehicleTypeID int64) (Configs, error) {
	fmt.Printf("fetching configurations for vehicle %v\n", vehicleTypeID)

	if utils.FlipCoin() {
		return Configs{}, errors.New("config store is down")
	}

	return Configs{
		StartingFare:  1.5,
		PerKMFare:     2,
		PerMinuteFare: 3,
	}, nil
}

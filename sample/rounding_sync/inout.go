package rounding

import "github.com/jamestrandung/go-component/sample/dto"

type Input interface {
	GetRunningFare() *dto.Fare
}

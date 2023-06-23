package routing

import "github.com/jamestrandung/go-concurrency/v2/async"

type RoutingFuture interface {
	GetDistanceInKM() float64
	GetDurationInSeconds() float64
}

type future struct {
	task async.Task[output]
}

func (f future) GetDistanceInKM() float64 {
	r := f.task.ResultOrDefault(output{})
	return r.distanceInKM
}

func (f future) GetDurationInSeconds() float64 {
	r := f.task.ResultOrDefault(output{})
	return r.durationInSeconds
}

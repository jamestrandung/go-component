package surge

import "github.com/jamestrandung/go-concurrency/v2/async"

type SurgeFuture interface {
	GetSurge() float64
}

type future struct {
	task async.Task[output]
}

func (f future) GetSurge() float64 {
	r := f.task.ResultOrDefault(output{})
	return r.surge
}

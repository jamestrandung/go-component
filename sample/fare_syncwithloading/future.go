package fare

import "github.com/jamestrandung/go-concurrency/v2/async"

type FareFuture interface {
	GetMetadata() Metadata
}

type future struct {
	task async.Task[output]
}

func (f future) GetMetadata() Metadata {
	r := f.task.ResultOrDefault(output{})
	return r.metadata
}

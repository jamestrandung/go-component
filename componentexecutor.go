package component

import (
	"context"

	"github.com/jamestrandung/go-concurrency/v2/async"
)

//go:generate mockery --name IComponentExecutor --case underscore --inpackage
type IComponentExecutor interface {
	invokeSyncTask(ctx context.Context) error
	canBeInvokedAsync() bool
	invokeAsyncTask(ctx context.Context) error
	cancel(err error)
	InvokeExecutingTask(ctx context.Context) error
}

// ComponentExecutor encapsulates the tasks that need to be
// executed to carry out the business logic of a component.
type ComponentExecutor[V any, T any] struct {
	loadingTask        async.Task[V]
	executingSyncTask  async.Task[T]
	executingAsyncTask async.Task[T]
}

func (p ComponentExecutor[V, T]) invokeSyncTask(ctx context.Context) error {
	if p.executingSyncTask != nil {
		return p.executingSyncTask.ExecuteSync(ctx).Error()
	}

	return nil
}

func (p ComponentExecutor[V, T]) canBeInvokedAsync() bool {
	return p.loadingTask != nil || p.executingAsyncTask != nil
}

func (p ComponentExecutor[V, T]) invokeAsyncTask(ctx context.Context) error {
	// Errors from loading task will be handled by sync components
	if p.loadingTask != nil {
		p.loadingTask.ExecuteSync(ctx)

		return nil
	}

	// Errors from async tasks will stop the entire flow
	if p.executingAsyncTask != nil {
		return p.executingAsyncTask.ExecuteSync(ctx).Error()
	}

	return nil
}

func (p ComponentExecutor[V, T]) cancel(err error) {
	if p.loadingTask != nil {
		p.loadingTask.CancelWithReason(err)
	}

	if p.executingSyncTask != nil {
		p.executingSyncTask.CancelWithReason(err)
	}

	if p.executingAsyncTask != nil {
		p.executingAsyncTask.CancelWithReason(err)
	}
}

func (p ComponentExecutor[V, T]) InvokeExecutingTask(ctx context.Context) error {
	return p.GetExecutingTask().ExecuteSync(ctx).Error()
}

func (p ComponentExecutor[V, T]) GetExecutingTask() async.Task[T] {
	// ComponentExecutor must contain either executingSyncTask
	// or executingAsyncTask but not both at the same time.
	if p.executingSyncTask != nil {
		return p.executingSyncTask
	}

	return p.executingAsyncTask
}

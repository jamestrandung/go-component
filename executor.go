package component

import (
	"context"

	"github.com/jamestrandung/go-concurrency/v2/async"
)

//go:generate mockery --name IExecutor --case underscore --inpackage
type IExecutor interface {
	invokeSyncTask(ctx context.Context) error
	canBeInvokedAsync() bool
	invokeAsyncTask(ctx context.Context) error
	cancel(err error)
	InvokeExecutingTask(ctx context.Context) error
}

// ExecutorWithLoading encapsulates the tasks that need to be executed to carry
// out the business logic of a synchronous component with loading logic.
type ExecutorWithLoading[V any, T any] struct {
	loadingTask       async.Task[V]
	executingSyncTask async.Task[T]
}

func (e ExecutorWithLoading[V, T]) invokeSyncTask(ctx context.Context) error {
	if e.executingSyncTask != nil {
		return e.executingSyncTask.ExecuteSync(ctx).Error()
	}

	return nil
}

func (e ExecutorWithLoading[V, T]) canBeInvokedAsync() bool {
	return e.loadingTask != nil
}

func (e ExecutorWithLoading[V, T]) invokeAsyncTask(ctx context.Context) error {
	// Errors from loading task will be handled by sync components
	if e.loadingTask != nil {
		e.loadingTask.ExecuteSync(ctx)

		return nil
	}

	return nil
}

func (e ExecutorWithLoading[V, T]) cancel(err error) {
	if e.loadingTask != nil {
		e.loadingTask.CancelWithReason(err)
	}

	if e.executingSyncTask != nil {
		e.executingSyncTask.CancelWithReason(err)
	}
}

func (e ExecutorWithLoading[V, T]) InvokeExecutingTask(ctx context.Context) error {
	return e.GetExecutingTask().ExecuteSync(ctx).Error()
}

func (e ExecutorWithLoading[V, T]) GetExecutingTask() async.Task[T] {
	return e.executingSyncTask
}

// Executor encapsulates the tasks that need to be executed to carry
// out the business logic of a component without loading logic.
type Executor[T any] struct {
	executingSyncTask  async.Task[T]
	executingAsyncTask async.Task[T]
}

func (e Executor[T]) invokeSyncTask(ctx context.Context) error {
	if e.executingSyncTask != nil {
		return e.executingSyncTask.ExecuteSync(ctx).Error()
	}

	return nil
}

func (e Executor[T]) canBeInvokedAsync() bool {
	return e.executingAsyncTask != nil
}

func (e Executor[T]) invokeAsyncTask(ctx context.Context) error {
	// Errors from async tasks will stop the entire flow
	if e.executingAsyncTask != nil {
		return e.executingAsyncTask.ExecuteSync(ctx).Error()
	}

	return nil
}

func (e Executor[T]) cancel(err error) {
	if e.executingSyncTask != nil {
		e.executingSyncTask.CancelWithReason(err)
	}

	if e.executingAsyncTask != nil {
		e.executingAsyncTask.CancelWithReason(err)
	}
}

func (e Executor[T]) InvokeExecutingTask(ctx context.Context) error {
	return e.GetExecutingTask().ExecuteSync(ctx).Error()
}

func (e Executor[T]) GetExecutingTask() async.Task[T] {
	// ExecutorWithLoading must contain either executingSyncTask or
	// executingAsyncTask but not both at the same time.
	if e.executingSyncTask != nil {
		return e.executingSyncTask
	}

	return e.executingAsyncTask
}

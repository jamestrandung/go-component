package component

import (
	"context"

	"github.com/jamestrandung/go-concurrency/v2/async"
)

// AsyncComponent represents those components that can be executed
// concurrently together with other AsyncComponent. They can also
// get executed sequentially with other SyncComponent if engineers
// choose to do so.
//
//go:generate mockery --name AsyncComponent --case underscore --inpackage
type AsyncComponent[T any] interface {
	Execute(ctx context.Context) (T, error)
}

// SyncComponent represents those components that must be executed
// sequentially together with other SyncComponent. These components
// must not be executed concurrently with other AsyncComponent.
//
// Engineers should use SyncComponentWithLoading instead if they
// need to perform some loading before executing the main logic.
//
//go:generate mockery --name SyncComponent --case underscore --inpackage
type SyncComponent[T any] interface {
	ExecuteSync(ctx context.Context) (T, error)
}

// LoadData contains the data and/or the error that was returned
// by the loading task. SyncComponentWithLoading is responsible
// for handling this error in its executing logic.
type LoadData[V any] struct {
	Data V
	Err  error
}

// SyncComponentWithLoading represents those components that must be
// executed sequentially together with other SyncComponent. They must
// not be executed concurrently with other AsyncComponent.
//
// Engineers should use SyncComponent instead if they don't need to
// perform any loading before executing the main logic.
//
//go:generate mockery --name SyncComponentWithLoading --case underscore --inpackage
type SyncComponentWithLoading[V any, T any] interface {
	Load(ctx context.Context) (V, error)
	ExecuteSync(ctx context.Context, data LoadData[V]) (T, error)
}

// CreateSyncExecutor returns an Executor encapsulating the executing
// task that would be handled by the given SyncComponent.
func CreateSyncExecutor[T any](c SyncComponent[T]) Executor[T] {
	return Executor[T]{
		executingSyncTask: async.NewTask[T](
			func(ctx context.Context) (T, error) {
				return c.ExecuteSync(ctx)
			},
		),
	}
}

// CreateAsyncExecutor returns an Executor encapsulating the executing
// task that would be handled by the given AsyncComponent.
func CreateAsyncExecutor[T any](c AsyncComponent[T]) Executor[T] {
	return Executor[T]{
		executingAsyncTask: async.NewTask[T](
			func(ctx context.Context) (T, error) {
				return c.Execute(ctx)
			},
		),
	}
}

// CreateSyncExecutorWithLoading returns an ExecutorWithLoading encapsulating the
// loading & executing tasks that would be handled by the given component.
func CreateSyncExecutorWithLoading[V any, T any](c SyncComponentWithLoading[V, T]) ExecutorWithLoading[V, T] {
	loadingTask := async.NewTask[V](
		func(ctx context.Context) (V, error) {
			return c.Load(ctx)
		},
	)

	executingSyncTask := async.NewTask[T](
		func(ctx context.Context) (T, error) {
			// Block & wait
			data, err := loadingTask.Outcome()

			return c.ExecuteSync(
				ctx,
				LoadData[V]{
					Data: data,
					Err:  err,
				},
			)
		},
	)

	return ExecutorWithLoading[V, T]{
		loadingTask:       loadingTask,
		executingSyncTask: executingSyncTask,
	}
}

// CreateSyncOrchestratingExecutor returns a component that is meant for orchestrating
// some logic without returning any values beside throwing an error if necessary.
func CreateSyncOrchestratingExecutor(doFn func(ctx context.Context) error) Executor[any] {
	return Executor[any]{
		executingSyncTask: async.NewTask[any](
			func(ctx context.Context) (interface{}, error) {
				return nil, doFn(ctx)
			},
		),
	}
}

// CreateSyncOrchestratingExecutorWithResult returns a component that is meant
// for orchestrating some logic that returns some values and throws an error
// if necessary.
func CreateSyncOrchestratingExecutorWithResult[T any](doFn func(ctx context.Context) (T, error)) (Executor[T], async.Task[T]) {
	t := async.NewTask[T](
		func(ctx context.Context) (T, error) {
			return doFn(ctx)
		},
	)

	return Executor[T]{
		executingSyncTask: t,
	}, t
}

package component

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/jamestrandung/go-concurrency/v2/async"

	"github.com/stretchr/testify/assert"
)

func TestForkJoinFailingFast(t *testing.T) {
	defer func(original func(flow ExecutionFlow, currentLayerIdx int, err error)) {
		cancelTasks = original
	}(cancelTasks)

	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "empty flow",
			test: func(t *testing.T) {
				actual := ForkJoinFailingFast(context.Background(), ExecutionFlow{})
				assert.Nil(t, actual)
			},
		},
		{
			desc: "flow with nil Executors",
			test: func(t *testing.T) {
				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: nil,
					},
				)

				assert.Nil(t, actual)
			},
		},
		{
			desc: "executing task of sync component returns an error",
			test: func(t *testing.T) {
				e := ExecutorWithLoading[int, int]{
					loadingTask:       async.Completed(1, nil),
					executingSyncTask: async.Completed(0, assert.AnError),
				}

				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: [][]IExecutor{
							{e},
						},
					},
				)

				assert.Equal(t, assert.AnError, actual)
			},
		},
		{
			desc: "executing task of async component returns an error",
			test: func(t *testing.T) {
				e := Executor[int]{
					executingAsyncTask: async.Completed(0, assert.AnError),
				}

				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: [][]IExecutor{
							{e},
						},
					},
				)

				assert.Equal(t, assert.AnError, actual)
			},
		},
		{
			desc: "executing task of sync component returns no error",
			test: func(t *testing.T) {
				e := ExecutorWithLoading[int, int]{
					loadingTask:       async.Completed(0, assert.AnError),
					executingSyncTask: async.Completed(1, nil),
				}

				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: [][]IExecutor{
							{e},
						},
					},
				)

				assert.Nil(t, actual)
			},
		},
		{
			desc: "executing task of async component returns no error",
			test: func(t *testing.T) {
				tp1 := Executor[int]{
					executingAsyncTask: async.Completed(1, nil),
				}

				tp2 := Executor[int]{
					executingAsyncTask: async.Completed(2, nil),
				}

				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: [][]IExecutor{
							{tp1, tp2},
						},
					},
				)

				assert.Nil(t, actual)
			},
		},
		{
			desc: "executing tasks of sync component are executed sequentially",
			test: func(t *testing.T) {
				var val int

				tp1 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							val = 1
							return 1, nil
						},
					),
				}

				tp2 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							val = 2
							return 2, nil
						},
					),
				}

				var wg sync.WaitGroup
				wg.Add(1)

				var anotherVal int
				tp3 := Executor[int]{
					executingAsyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							defer wg.Done()

							anotherVal = 99
							return 1, nil
						},
					),
				}

				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: [][]IExecutor{
							{tp1, tp2, tp3},
						},
					},
				)

				wg.Wait()

				assert.Nil(t, actual)
				assert.Equal(t, 2, val, "Val must carry value assigned by the 2nd mock")
				assert.Equal(t, 99, anotherVal, "Val must carry value assigned by the 3rd mock")
			},
		},
		{
			desc: "one failing task from sync component will cancel all tasks",
			test: func(t *testing.T) {
				var isCancelTasksCalled bool
				doCancelTasks := cancelTasks
				cancelTasks = func(flow ExecutionFlow, currentLayerIdx int, err error) {
					isCancelTasksCalled = true
					doCancelTasks(flow, currentLayerIdx, err)
				}

				var val int

				tp1 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							val = 1
							return 1, nil
						},
					),
				}

				tp2 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							val = 2
							return 2, errors.New("error from sync task")
						},
					),
				}

				tp3 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							val = 3
							return 3, nil
						},
					),
				}

				tp4 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							val = 4
							return 4, nil
						},
					),
				}

				tp5 := Executor[int]{
					executingAsyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							<-time.After(1 * time.Second)
							return 5, errors.New("error from async task")
						},
					),
				}

				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: [][]IExecutor{
							{tp1, tp2, tp3, tp4, tp5},
						},
					},
				)

				assert.Equal(t, "error from sync task", actual.Error())
				assert.Equal(t, 2, val, "Val must carry value assigned by the 2nd mock right before returning an error")
				assert.True(t, isCancelTasksCalled)
			},
		},
		{
			desc: "one failing task from async component will cancel all tasks",
			test: func(t *testing.T) {
				var isCancelTasksCalled bool
				doCancelTasks := cancelTasks
				cancelTasks = func(flow ExecutionFlow, currentLayerIdx int, err error) {
					isCancelTasksCalled = true
					doCancelTasks(flow, currentLayerIdx, err)
				}

				var groupCtx context.Context

				tp1 := Executor[int]{
					executingAsyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							groupCtx = ctx
							return 1, errors.New("error from async task")
						},
					),
				}

				tp2 := Executor[int]{
					executingAsyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							<-time.After(1 * time.Second)
							return 2, nil
						},
					),
				}

				tp3 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							<-time.After(1 * time.Second)
							return 3, errors.New("error from sync task")
						},
					),
				}

				tp4 := ExecutorWithLoading[int, int]{
					loadingTask: async.Completed(0, assert.AnError),
					executingSyncTask: async.NewTask(
						func(ctx context.Context) (int, error) {
							<-time.After(1 * time.Second)
							return 4, errors.New("error from sync task")
						},
					),
				}

				actual := ForkJoinFailingFast(
					context.Background(),
					ExecutionFlow{
						Executors: [][]IExecutor{
							{tp1, tp2, tp3, tp4},
						},
					},
				)

				assert.Equal(t, "error from async task", actual.Error())
				assert.True(t, isCancelTasksCalled)
				assert.Equal(t, context.Canceled, groupCtx.Err(), "when one task fails, the context sent into each task should have been cancelled")
			},
		},
	}

	for _, scenario := range scenarios {
		sc := scenario
		t.Run(sc.desc, sc.test)
	}
}

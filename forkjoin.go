package component

import (
	"context"
	"strings"
	"sync"
)

// ForkJoinFailingFast invokes the executors in the given ExecutionFlow based on its type. If an executor comes from
// an async component, it will be executed asynchronously. If an executor comes from a sync component, its loading
// task will be executed asynchronously while its executing task will be executed synchronously based on the order
// of the given tasks.
//
// If any of the executing tasks of async or sync components returns an error, the function will stop immediately
// and return this error to the caller.
var ForkJoinFailingFast = func(ctx context.Context, flow ExecutionFlow) error {
	if len(flow.Executors) == 0 {
		return nil
	}

	if len(flow.Executors) == 1 {
		return doForkJoinFailingFast(ctx, flow, 0)
	}

	// Buffer of len(flow.Executors) to take at most 1 error from each layer of executors
	errChan := make(chan error, len(flow.Executors))

	for i := 0; i < len(flow.Executors); i++ {
		go func(idx int) {
			errChan <- doForkJoinFailingFast(ctx, flow, idx)
		}(i)
	}

	done := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			// If any of the goroutines returns an error, the
			// entire flow stops immediately.
			if err != nil {
				return err
			}

			// If err is nil, mark 1 layer as complete until all
			// layer completes successfully.
			done = done + 1
			if done == len(flow.Executors) {
				return nil
			}
		}
	}
}

var doForkJoinFailingFast = func(ctx context.Context, flow ExecutionFlow, currentLayerIdx int) error {
	executors := flow.Executors[currentLayerIdx]

	var errOnce sync.Once
	errChan := make(chan error, 1)

	var wg sync.WaitGroup

	// Execute loading + async components asynchronously
	for _, executor := range executors {
		if !executor.canBeInvokedAsync() {
			continue
		}

		wg.Add(1)

		e := executor
		go func() {
			defer wg.Done()

			err := e.invokeAsyncTask(ctx)

			// When err is async.ErrCancelled, it means this task is
			// being actively cancelled by the sync goroutine. We can
			// swallow this error and let the other goroutine return
			// an error to the caller.
			if err == nil || strings.Contains(err.Error(), "task cancelled with reason") {
				return
			}

			errOnce.Do(
				func() {
					// Release the main thread first before cancelling tasks
					errChan <- err
				},
			)

			cancelTasks(flow, currentLayerIdx, err)
		}()
	}

	wg.Add(1)

	// Execute sync components sequentially
	go func() {
		defer wg.Done()

		for _, executor := range executors {
			// Block & wait for error before executing the next component
			if err := executor.invokeSyncTask(ctx); err != nil {
				// When err is async.ErrCancelled, it means this task is
				// being actively cancelled by the async goroutine. We
				// must stop execution and let the other goroutine return
				// an error to the caller.
				if strings.Contains(err.Error(), "task cancelled with reason") {
					break
				}

				errOnce.Do(
					func() {
						// Release the main thread first before cancelling tasks
						errChan <- err
					},
				)

				cancelTasks(flow, currentLayerIdx, err)

				return
			}
		}
	}()

	// Wait & close when ALL goroutines have returned.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		// If any of the goroutines returns an error, the
		// entire flow stops immediately.
		return err
	case <-done:
		return nil
	}
}

var cancelTasks = func(flow ExecutionFlow, currentLayerIdx int, err error) {
	flow.Cancel(currentLayerIdx, err)
}

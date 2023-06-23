package component

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSyncExecutingTask(t *testing.T) {
	mockSyncComponent := &MockSyncComponent[int]{}
	mockSyncComponent.On("ExecuteSync", mock.Anything).
		Return(1, assert.AnError).
		Once()

	actual := CreateSyncExecutor[int](mockSyncComponent)
	assert.NotNil(t, actual.executingSyncTask)

	err := actual.invokeSyncTask(context.Background())
	assert.Equal(t, assert.AnError, err)

	result, err := actual.executingSyncTask.Outcome()
	assert.Equal(t, 1, result)
	assert.Equal(t, assert.AnError, err)
}

func TestCreateAsyncExecutingTask(t *testing.T) {
	mockAsyncComponent := &MockAsyncComponent[int]{}
	mockAsyncComponent.On("Execute", mock.Anything).
		Return(1, assert.AnError).
		Once()

	actual := CreateAsyncExecutor[int](mockAsyncComponent)
	assert.NotNil(t, actual.executingAsyncTask)

	err := actual.invokeAsyncTask(context.Background())
	assert.Equal(t, assert.AnError, err)

	result, err := actual.executingAsyncTask.Outcome()
	assert.Equal(t, 1, result)
	assert.Equal(t, assert.AnError, err)
}

func TestCreateSyncLoadingExecutingTask(t *testing.T) {
	mockSyncComponentWithLoading := &MockSyncComponentWithLoading[int, int]{}
	mockSyncComponentWithLoading.On("Load", mock.Anything).
		Return(1, assert.AnError).
		Once()
	mockSyncComponentWithLoading.On("ExecuteSync", mock.Anything, LoadData[int]{Data: 1, Err: assert.AnError}).
		Return(2, assert.AnError).
		Once()

	actual := CreateSyncExecutorWithLoading[int, int](mockSyncComponentWithLoading)
	assert.NotNil(t, actual.loadingTask)
	assert.NotNil(t, actual.executingSyncTask)

	actual.invokeAsyncTask(context.Background())

	err := actual.invokeSyncTask(context.Background())
	assert.Equal(t, assert.AnError, err)

	result, err := actual.executingSyncTask.Outcome()
	assert.Equal(t, 2, result)
	assert.Equal(t, assert.AnError, err)

	mock.AssertExpectationsForObjects(t, mockSyncComponentWithLoading)
}

func TestCreateOrchestratingTask(t *testing.T) {
	doFn := func(ctx context.Context) error {
		return assert.AnError
	}

	actual := CreateSyncOrchestratingExecutor(doFn)
	assert.NotNil(t, actual.executingSyncTask)

	err := actual.invokeSyncTask(context.Background())
	assert.Equal(t, assert.AnError, err)
}

func TestCreateOrchestratingTaskWithResult(t *testing.T) {
	scenarios := []struct {
		desc string
		test func(t *testing.T)
	}{
		{
			desc: "doFn returns an error",
			test: func(t *testing.T) {
				doFn := func(ctx context.Context) (interface{}, error) {
					return 1, assert.AnError
				}

				pair, future := CreateSyncOrchestratingExecutorWithResult(doFn)
				assert.NotNil(t, pair.executingSyncTask)

				err := pair.invokeSyncTask(context.Background())
				assert.Equal(t, assert.AnError, err)
				assert.Equal(t, 2, future.ResultOrDefault(2))
			},
		},
		{
			desc: "doFn returns no error",
			test: func(t *testing.T) {
				doFn := func(ctx context.Context) (interface{}, error) {
					return 1, nil
				}

				pair, future := CreateSyncOrchestratingExecutorWithResult(doFn)
				assert.NotNil(t, pair.executingSyncTask)

				err := pair.invokeSyncTask(context.Background())
				assert.Nil(t, err)
				assert.Equal(t, 1, future.ResultOrDefault(2))
			},
		},
	}

	for _, scenario := range scenarios {
		sc := scenario
		t.Run(sc.desc, sc.test)
	}
}

package fare

import (
	"github.com/jamestrandung/go-component"
	"github.com/jamestrandung/go-component/sample/fare_syncwithloading/dependencies"
)

type factory struct {
	configStore dependencies.IConfigStore
}

var f factory

func InitializeFactory(mapService dependencies.IConfigStore) {
	f = factory{
		configStore: mapService,
	}
}

func GetExecutorFuture(input Input) (component.ExecutorWithLoading[dependencies.Configs, output], FareFuture) {
	c := Component{
		configStore: f.configStore,
		input:       input,
	}

	e := component.CreateSyncExecutorWithLoading[dependencies.Configs, output](c)

	return e, future{
		task: e.GetExecutingTask(),
	}
}

package surge

import (
	"github.com/jamestrandung/go-component"
	"github.com/jamestrandung/go-component/sample/surge_async/dependencies"
)

type factory struct {
	surgeEngine dependencies.ISurgeEngine
}

var f factory

func InitializeFactory(surgeEngine dependencies.ISurgeEngine) {
	f = factory{
		surgeEngine: surgeEngine,
	}
}

func GetExecutorFuture(input Input) (component.Executor[output], SurgeFuture) {
	c := Component{
		surgeEngine: f.surgeEngine,
		input:       input,
	}

	e := component.CreateAsyncExecutor[output](c)

	return e, future{
		task: e.GetExecutingTask(),
	}
}

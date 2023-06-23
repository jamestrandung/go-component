package routing

import (
	"github.com/jamestrandung/go-component"
	"github.com/jamestrandung/go-component/sample/routing_async/dependencies"
)

type factory struct {
	mapService dependencies.IMapService
}

var f factory

func InitializeFactory(mapService dependencies.IMapService) {
	f = factory{
		mapService: mapService,
	}
}

func GetExecutorFuture(input Input) (component.Executor[output], RoutingFuture) {
	c := Component{
		mapService: f.mapService,
		input:      input,
	}

	e := component.CreateAsyncExecutor[output](c)

	return e, future{
		task: e.GetExecutingTask(),
	}
}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jamestrandung/go-component"
	"github.com/jamestrandung/go-component/sample/dto"
	fare "github.com/jamestrandung/go-component/sample/fare_syncwithloading"
	fareDep "github.com/jamestrandung/go-component/sample/fare_syncwithloading/dependencies"
	"github.com/jamestrandung/go-component/sample/request"
	rounding "github.com/jamestrandung/go-component/sample/rounding_sync"
	routing "github.com/jamestrandung/go-component/sample/routing_async"
	routingDep "github.com/jamestrandung/go-component/sample/routing_async/dependencies"
	surge "github.com/jamestrandung/go-component/sample/surge_async"
	surgeDep "github.com/jamestrandung/go-component/sample/surge_async/dependencies"
)

func init() {
	routing.InitializeFactory(routingDep.NewMapService())
	surge.InitializeFactory(surgeDep.NewSurgeEngine())
	fare.InitializeFactory(fareDep.NewConfigStore())
}

func main() {
	req := request.FareCalculationRequest{
		VehicleTypeID: 123,
		PickUp: dto.Location{
			Lat: 1.3148639,
			Lng: 103.7589081,
		},
		DropOff: dto.Location{
			Lat: 1.3544924,
			Lng: 103.9837306,
		},
		PickUpTime: time.Now(),
	}

	runningFare := dto.MakeRunningFare()

	routingExecutor, routingFuture := routing.GetExecutorFuture(req)
	surgeExecutor, surgeFuture := surge.GetExecutorFuture(req)

	// Wire the components together to use the outputs of a
	// component as the inputs of another component.
	//
	// Each component will automatically block & wait if the
	// component it depends on has not completed yet.
	fareExecutor, fareFuture := fare.GetExecutorFuture(
		struct {
			request.FareCalculationRequest
			routing.RoutingFuture
			surge.SurgeFuture
			dto.RunningFare
		}{
			req,
			routingFuture,
			surgeFuture,
			runningFare,
		},
	)
	roundingExecutor := rounding.GetExecutor(runningFare)

	// The order of appending is the same order that
	// synchronous components will get executed.
	executionFlow := component.NewExecutionFlowBuilder().
		Append(
			routingExecutor,
			surgeExecutor,
			fareExecutor,
			roundingExecutor,
		).Get()

	// ForkJoin will execute all async components and loading executors in parallel to
	// maximize performance. At the same time, it will execute synchronous components
	// 1-by-1 in the exact order in which they were appended to the execution flow.
	//
	// The very first error thrown by any executor will end ForkJoin immediately.
	if err := component.ForkJoinFailingFast(context.Background(), executionFlow); err != nil {
		fmt.Printf("ending execution flow early due to error: %v \n", err.Error())

		return
	}

	fmt.Printf("calculated fare: %v\n", runningFare.GetRunningFare().Amount)
	fmt.Printf("applied fare configs: %v\n", fareFuture.GetMetadata())
}

# go-component

This package provides simple yet powerful abstractions for creating synchronous and asynchronous components that can be 
plugged and played to realize a truly loosely-coupled architecture. You can read the code & run `main.go` in the `sample` 
package for a demonstration.

The main properties of the proposed model are:
- Structural - applications are described as components, their inputs and their outputs encapsulated as futures.
- Asynchronous/synchronous - you can decide the type of a component by implementing the corresponding interface. The order
of execution for synchronous components can be determined when building an execution flow.
- Isolated - each component describes what it needs for its logic via an Input interface, which can then be provided by
1 or more components via their outputs wrapped as futures. State is not shared across components.
- Concurrent - execution is greedy, all asynchronous logic will get started in a goroutine immediately when an execution
flow gets triggered. Synchronous logic will still get executed in the desired sequence. Components will automatically block 
& wait for each other when they access methods of futures that have not resolved yet.
- Fail fast - each component must handle its own errors (e.g. by using some default values as output or by logging the error
& then ignoring it to return early in case some logic can be bypassed). If an error is returned by any components, the entire 
execution flow will stop immediately and this error will be used as the final result of this execution.
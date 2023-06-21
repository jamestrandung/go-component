package component

// ExecutionFlow ...
type ExecutionFlow struct {
	Executors [][]IComponentExecutor
}

// Cancel cancels all executors from the given layer down.
func (f ExecutionFlow) Cancel(firstLayerIdx int, err error) {
	for i := firstLayerIdx; i < len(f.Executors); i++ {
		for j := 0; j < len(f.Executors[i]); j++ {
			f.Executors[i][j].cancel(err)
		}
	}
}

// ExecutionFlowBuilder ...
type ExecutionFlowBuilder struct {
	executorLayers [][]IComponentExecutor
}

// NewExecutionFlowBuilder ...
func NewExecutionFlowBuilder() *ExecutionFlowBuilder {
	return &ExecutionFlowBuilder{
		executorLayers: make([][]IComponentExecutor, 1),
	}
}

// Append appends the given executors to the end of the current executor layer.
func (b *ExecutionFlowBuilder) Append(executors ...IComponentExecutor) *ExecutionFlowBuilder {
	currentIdx := len(b.executorLayers) - 1
	b.executorLayers[currentIdx] = append(b.executorLayers[currentIdx], executors...)

	return b
}

// NextLayer moves the builder to the next layer of executors, effectively finalize the current layer.
func (b *ExecutionFlowBuilder) NextLayer() *ExecutionFlowBuilder {
	b.executorLayers = append(b.executorLayers, []IComponentExecutor{})

	return b
}

// Get returns the current flow.
func (b *ExecutionFlowBuilder) Get() ExecutionFlow {
	return ExecutionFlow{
		Executors: b.executorLayers,
	}
}

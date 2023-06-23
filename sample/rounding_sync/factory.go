package rounding

import (
	"github.com/jamestrandung/go-component"
)

func GetExecutor(input Input) component.Executor[any] {
	c := Component{
		input: input,
	}

	return component.CreateSyncExecutor[any](c)
}

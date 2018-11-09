package node

import "context"

type Node interface {
	GetContract() Contract
	GetPossibleOutputs() []string
	GetDescription() string
	Run(context.Context) (string, context.Context, error)
}

type Contract struct {
	Inputs  map[string]string
	Outputs map[string]string
}

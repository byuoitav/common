package flow

import (
	"context"
	"time"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/troubleshootingflow/node"
)

type FlowDefinition struct {
	//Essential Initial Definition
	FlowChart  map[string]FlowState
	Triggers   []string
	StartState string
	//Time to Rerun
	TTR time.Duration
	//Time to wait before marking alert as flapping.
	FlappingInterval time.Duration
}

type Flow struct {
	FlowDefinition
	//Filled during process

	StateHistroy []Step
	CurrentState string
	Ctx          context.Context //Context

	//Flow-wide Flag(s)
	//If set to true, the alert clearing during run of flow will not abort flow
	FlappingFlag    bool
	TriggeringAlert Alert
}

type FlowState struct {
	ID          string
	node        node.Node
	Transitions map[string]string
}

type Step struct {
	ID     string
	Result string
}

/*
StartFlow assumes that a Flow has been initialized and that every FlowState has been linked to its associated node
*/
func StartFlow(flow Flow, alert Alert, ctx context.Context) error {
	flow.CurrentState = flow.StartState
	//TODO -> flow.Ctx = ctx
	var curr FlowState
	var result string
	var err error
	for {
		curr = flow.FlowChart[flow.CurrentState]
		result, flow.Ctx, err = curr.node.Run(flow.Ctx)
		if err != nil {
			log.L.Errorf("Flow Aborting: %v", err.Error())
			//TODO: Mark as Cancel
			return err
		}
		flow.StateHistory = append(flow.StateHistory, Step{curr.ID, result})
		if len(curr.Transitions) < 1 {
			return nil
			//TODO: END
		}
		flow.CurrentState = curr.Transitions[result]
	}
	return nil
}

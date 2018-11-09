package eventnode

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/v2/events"
	"github.com/byuoitav/messenger"
)

// EventNode is a node in the event mesh of a room
type EventNode struct {
	Node *messenger.Node
}

// NewEventNode returns a new event node with the given data.
// filters: an array of strings to filter events recieved by
// routerAddress: address of the router to connect to
// name: name of event node
func NewEventNode(name, routerAddress string, filters []string) *EventNode {
	n := &EventNode{
		Node: messenger.NewNode(name, filters),
	}

	n.Node.ConnectToRouter(routerAddress)
	return n
}

// PublishEvent publishes an event with the given tag
func (n *EventNode) PublishEvent(tag string, event events.Event) *nerr.E {
	log.L.Debugf("Sending an event with tag: %v", tag)

	// turn event into bytes
	bytes, err := json.Marshal(event)
	if err != nil {
		return nerr.Translate(err).Addf("failed to marshal event")
	}

	n.Node.Write(messenger.Message{Header: tag, Body: bytes})
	return nil
}

func (n *EventNode) Read() (events.Event, *nerr.E) {
	var toReturn events.Event
	msg := n.Node.Read()

	err := json.Unmarshal(msg.Body, &toReturn)
	if err != nil {
		return toReturn, nerr.Create(fmt.Sprintf("unable to unmarshal message: %s", msg.Body), reflect.TypeOf("").String())
	}

	return toReturn, nil
}

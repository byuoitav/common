package events

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/byuoitav/messenger"
)

type EventNode struct {
	Node *messenger.Node
}

// filters: an array of strings to filter events recieved by
// addrs: addresses of subscriber to subscribe to
// name: name of event node
func NewEventNode(name, address string, filters []string) *EventNode {
	n := &EventNode{
		Node: messenger.NewNode(name, filters),
	}

	n.Node.ConnectToRouter(address)

	return n
}

func (n *EventNode) PublishEvent(eventType string, event Event) error {
	// turn event into bytes
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	n.Node.Write(messenger.Message{Header: eventType, Body: bytes})
	return nil
}

func (n *EventNode) Read() (Event, error) {
	var toReturn Event
	msg := n.Node.Read()

	err := json.Unmarshal(msg.Body, &toReturn)
	if err != nil {
		return toReturn, errors.New(fmt.Sprintf("unable to unmarshal message: %s", msg.Body))
	}

	return toReturn, nil
}

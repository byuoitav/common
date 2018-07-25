package events

import (
	"log"

	"github.com/byuoitav/messenger"
	"github.com/fatih/color"
)

func NewRouter(routingTable map[string][]string, addrs []string) (*messenger.Router, error) {
	r := messenger.NewRouter()

	go r.StartRouter(routingTable)

	err := r.ConnectToRouters(addrs, routingTable)
	if err != nil {
		log.Printf(color.HiRedString("failed to connect to peers: %s", err))
		return r, err
	}

	return r, nil
}

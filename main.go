package main

import (
	"fmt"
	"github.com/hashicorp/mdns"
)

func main() {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	mdns.Lookup("_http._tcp", entriesCh)
	close(entriesCh)
}

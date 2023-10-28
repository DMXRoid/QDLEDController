package discovery

import (
	"fmt"
	"github.com/DMXRoid/QDLEDController/v2/led"
	"github.com/hashicorp/mdns"
	"time"
)

var mdnsResults chan *mdns.ServiceEntry

func Init() {
	mdnsResults = make(chan *mdns.ServiceEntry, 50)
	go func() {
		for entry := range mdnsResults {
			err := led.Register(entry.AddrV4.String())
			if err != nil {
				fmt.Println(err)
			}

		}
	}()
	go func() {
		for {
			discover()
			time.Sleep(10 * time.Second)
		}

	}()
}

func discover() {
	// Start the lookup
	p := mdns.DefaultParams("_qdled._tcp")
	p.DisableIPv6 = true
	p.Entries = mdnsResults
	mdns.Lookup("_qdled._tcp", mdnsResults)
}

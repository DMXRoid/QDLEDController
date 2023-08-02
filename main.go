package main

import (
	"github.com/DMXRoid/QDLEDController/v2/discovery"
	"github.com/DMXRoid/QDLEDController/v2/led"
	"github.com/DMXRoid/QDLEDController/v2/service"
)

var wait chan bool

func main() {
	led.Init()
	discovery.Init()
	service.Init()
	<-wait

}

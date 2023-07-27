package led

import (
	"encoding/json"
	"fmt"
	pb "github.com/DMXRoid/QDLEDController/proto"
	"net/http"
)

var c = &http.Client{}

func Init() {

}

func (led *pb.LED) Update(doRestart bool) error {
	var err error

	return err
}

func (led *pb.LED) Refresh() error {
	var err error
	r := map[string]interface{}{
		"action": "get-config",
	}
	err = json.Unmarshal(c.Post("http://"+led.GetIpAddress())+"/act", led)
	return err
}

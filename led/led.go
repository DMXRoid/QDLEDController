package led

import (
	"encoding/json"
	"fmt"
	pb "github.com/DMXRoid/QDLEDController/v2/proto"
	pbjson "google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"strings"
)

type LED struct {
	*pb.LED
}

var RegisteredLEDs map[string]*LED

var c = &http.Client{}

func Init() {
	RegisteredLEDs = make(map[string]*LED)
}

func Register(ipAddress string) error {
	var err error
	if _, ok := RegisteredLEDs[ipAddress]; !ok {
		fmt.Println("Adding", ipAddress)
		l := &LED{
			LED: &pb.LED{
				IpAddress: ipAddress,
			},
		}
		err = l.Probe()
		if err == nil {
			err = l.Refresh()
			if err == nil {
				RegisteredLEDs[ipAddress] = l
			}
		}
	}
	return err
}

func (led *LED) UpdateConfig() error {
	var err error
	m := pbjson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}
	j, err := m.Marshal(led.LED)
	if err == nil {
		err = led.Update("config", j)
	}
	return err
}

func (led *LED) Update(configType string, p []byte) error {
	var err error

	payload := fmt.Sprintf("{\"action\": \"set-%s\", \"%s\": %s}", configType, configType, p)
	fmt.Println(payload)
	_, err = c.Post(led.GetActionURL(), "application/json", strings.NewReader(payload))
	if err != nil {
		fmt.Println("error:::", err)
	}
	return err
}

func (led *LED) UpdateColors() error {
	var err error
	m := pbjson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}

	j, err := m.Marshal(led.LED.Color)

	if err == nil {
		err = led.Update("color-config", j)
	}
	return err
}

func (led *LED) UpdateLights() error {
	var err error
	m := pbjson.MarshalOptions{
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}

	j, err := m.Marshal(led.LED.Lights)

	if err == nil {
		err = led.Update("light-config", j)
	}
	return err
}

func (led *LED) Probe() error {
	_, err := c.Get(led.GetProbeURL())
	return err
}

func (led *LED) GetProbeURL() string {
	return fmt.Sprintf("http://%s/ok", led.GetIpAddress())
}

func (led *LED) GetActionURL() string {
	return fmt.Sprintf("http://%s/act", led.GetIpAddress())
}

func (led *LED) Refresh() error {
	var err error

	/*r := map[string]interface{}{
		"action": "get-config",
	}*/

	b, err := c.Post(led.GetActionURL(), "application/json", strings.NewReader(`{"action":"get-config"}`))
	if err == nil {
		defer b.Body.Close()
		bodyJSON, err := io.ReadAll(b.Body)
		fmt.Println(string(bodyJSON))
		r := make(map[string]interface{})
		err = json.Unmarshal(bodyJSON, &r)
		if err != nil {
			fmt.Println("error", err)
		} else {
			rr, err := json.Marshal(r["message"])
			if err != nil {
				fmt.Println("error: ", err)
			} else {
				err = pbjson.Unmarshal(rr, led.LED)
				if err != nil {
					fmt.Println("error:", err)
				} else {
					fmt.Println(fmt.Sprintf("%+v", led.LED))
				}
			}
		}
	}

	return err
}

func Update(l *pb.LED) error {
	var err error
	if _, ok := RegisteredLEDs[l.GetIpAddress()]; ok {
		RegisteredLEDs[l.GetIpAddress()].LED = l
		err = RegisteredLEDs[l.GetIpAddress()].UpdateConfig()
	}
	return err
}

func JSONToLED(j string, l *LED) error {
	var err error

	return err
}

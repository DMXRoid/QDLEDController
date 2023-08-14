package led

import (
	"encoding/json"
	"fmt"
	pb "github.com/DMXRoid/QDLEDController/v2/proto"
	pbjson "google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LED struct {
	*pb.LED
	ProbeFailCount int
}

var RegisteredLEDs map[string]*LED

var c = &http.Client{}

var mutex *sync.RWMutex

var failedProbes chan string

func Init() {
	failedProbes = make(chan string, 20)
	RegisteredLEDs = make(map[string]*LED)
	mutex = &sync.RWMutex{}
	time.AfterFunc(60*time.Second, FreshenUp)
}

func Register(ipAddress string) error {
	var err error
	mutex.Lock()
	defer mutex.Unlock()
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
	mutex.RLock()
	defer mutex.RUnlock()
	var err error
	if _, ok := RegisteredLEDs[l.GetIpAddress()]; ok {
		if l.GetMdnsName() == RegisteredLEDs[l.GetIpAddress()].GetMdnsName() && l.GetDataPin() == RegisteredLEDs[l.GetIpAddress()].GetDataPin() {
			RegisteredLEDs[l.GetIpAddress()].LED.Color = l.Color
			RegisteredLEDs[l.GetIpAddress()].LED.Lights = l.Lights
			err = RegisteredLEDs[l.GetIpAddress()].UpdateColors()
			if err == nil {
				err = RegisteredLEDs[l.GetIpAddress()].UpdateLights()
			}
		} else {
			RegisteredLEDs[l.GetIpAddress()].LED = l
			err = RegisteredLEDs[l.GetIpAddress()].UpdateConfig()
		}

	}
	return err
}

func SelfRegister(l *pb.LED) error {
	var err error
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := RegisteredLEDs[l.GetIpAddress()]; !ok {
		RegisteredLEDs[l.GetIpAddress()] = &LED{
			LED: l,
		}
	}
	return err
}

func Sync(sourceIdentifier, targetIdentifier string) error {
	var err error
	mutex.Lock()
	defer mutex.Unlock()

	if s, ok := RegisteredLEDs[sourceIdentifier]; ok {
		if t, ok := RegisteredLEDs[targetIdentifier]; ok {
			t.Color = s.Color
			t.Lights = s.Lights
			t.UpdateColors()
			t.UpdateLights()
		} else {
			err = fmt.Errorf(fmt.Sprintf("Unregistered target: %s", targetIdentifier))
		}

	} else {
		err = fmt.Errorf(fmt.Sprintf("Unregistered source: %s", sourceIdentifier))
	}

	return err

}

func FreshenUp() {
	var err error
	mutex.Lock()
	defer mutex.Unlock()

	for k, led := range RegisteredLEDs {
		go func() {
			var err error
			err = led.Probe()
			if err != nil {
				failedProbes <- led.GetIpAddress()
			} else {
				err = led.Refresh()
			}
		}()
		if err == nil {
			err = led.Refresh()
		}

		if err != nil {
			led.ProbeFailCount++
			if led.ProbeFailCount > 3 {
				delete(RegisteredLEDs, k)
			}
		} else {
			led.ProbeFailCount = 0
		}

	}

	time.AfterFunc(60*time.Second, FreshenUp)
}

package led

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DMXRoid/QDLEDController/v2/db"
	pb "github.com/DMXRoid/QDLEDController/v2/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	pbjson "google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
	"strings"
	"time"
)

type LED struct {
	*pb.LED
	ProbeFailCount int
}

type LEDs map[string]*LED

var c = &http.Client{}

func Init() {
	time.AfterFunc(300*time.Second, FreshenUp)
}

func Register(ipAddress string) error {
	var err error

	if _, err := GetLED(ipAddress); err != nil {
		l := &LED{
			LED: &pb.LED{
				IpAddress: ipAddress,
			},
		}
		err = l.Probe()
		if err == nil {
			err = l.Refresh()
			if err == nil {
				err = l.Save()
			}
		}
	} else {
		err = nil
	}

	return err
}

func (led *LED) Save() error {
	var err error
	c := db.GetCollection("led")
	f := bson.D{{
		"led.ipaddress",
		led.GetIpAddress(),
	}}
	_, err = c.ReplaceOne(context.TODO(), f, led, options.Replace().SetUpsert(true))
	if err != nil {
		fmt.Println(fmt.Sprintf(":::SAVE ERROR::: %s", err))
	}
	return err
}

func (led *LED) Delete() error {
	var err error
	c := db.GetCollection("led")
	f := bson.D{{
		"led.ipaddress",
		led.GetIpAddress(),
	}}
	_, err = c.DeleteOne(context.TODO(), f)
	if err != nil {
		fmt.Println(fmt.Sprintf(":::SAVE ERROR::: %s", err))
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
	fmt.Println(fmt.Sprintf("updating %s", led.LED.GetIpAddress()))
	fmt.Println(":::::::UPDATING:::::::")
	fmt.Println(payload)
	_, err = c.Post(led.GetActionURL(), "application/json", strings.NewReader(payload))
	if err != nil {
		fmt.Println("error:::", err)
	}
	led.Save()
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

func (led *LED) Toggle() error {
	led.LED.Lights.IsEnabled = !led.LED.Lights.IsEnabled
	return led.UpdateLights()
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
					err = led.Save()
					//fmt.Println(fmt.Sprintf("%+v", led.LED))
				}
			}
		}
	}

	return err
}

func GetLED(ipAddress string) (*LED, error) {
	var led *LED
	var err error
	led = &LED{}
	err = db.GetCollection("led").FindOne(context.TODO(), bson.D{{"led.ipaddress", ipAddress}}).Decode(led)
	return led, err
}

func GetAllLEDs() ([]*LED, error) {
	var leds []*LED
	var err error
	//opts := options.Find().SetSort(bson.D{{"led.friendlyName", 1}})
	opts := bson.D{}
	c, err := db.GetCollection("led").Find(context.TODO(), opts)
	if err == nil {
		err = c.All(context.TODO(), &leds)
	}
	if err != nil {
		fmt.Println(err)
	}

	return leds, err
}

func Update(l *pb.LED) error {
	var led *LED
	var err error
	led, err = GetLED(l.GetIpAddress())
	if err == nil {
		if l.GetMdnsName() == led.GetMdnsName() && l.GetDataPin() == led.GetDataPin() {
			led.LED.Color = l.Color
			led.LED.Lights = l.Lights
			err = led.UpdateColors()
			if err == nil {
				err = led.UpdateLights()
			}
		} else {

			fmt.Println(fmt.Sprintf("Full config change for %+v", l))
			led.LED = l
			err = led.UpdateConfig()
		}

	}
	return err
}

func SelfRegister(l *pb.LED) error {
	var err error
	err = Register(l.GetIpAddress())
	return err
}

func Sync(sourceIdentifier string, targetIdentifier []string) error {
	var err error
	fmt.Println("SYNCING STARTS")
	fmt.Println(fmt.Sprintf("%+v", targetIdentifier))
	if s, err := GetLED(sourceIdentifier); err == nil {
		for _, tid := range targetIdentifier {
			if t, err := GetLED(tid); err == nil {
				t.Color = s.Color
				t.Lights = s.Lights
				t.UpdateColors()
				t.UpdateLights()
				fmt.Printf(fmt.Sprintf("syncing %s to %s", sourceIdentifier, tid))
			} else {
				err = fmt.Errorf(fmt.Sprintf("Unregistered target: %s", tid))
			}
		}

	} else {
		err = fmt.Errorf(fmt.Sprintf("Unregistered source: %s", sourceIdentifier))
	}

	return err

}

func Toggle(ipAddress string) error {
	var err error
	var l *LED

	l, err = GetLED(ipAddress)
	if err == nil {
		err = l.Toggle()
	}
	return err
}

func FreshenUp() {
	var err error
	var leds []*LED
	leds, err = GetAllLEDs()
	if err == nil {
		for _, l := range leds {
			err = l.Probe()
			if err != nil {
				l.ProbeFailCount++
				if l.ProbeFailCount >= 3 {
					l.Delete()
				} else {
					l.Save()
				}
			} else {
				l.ProbeFailCount = 0
				l.Save()
			}
		}
	}
	time.AfterFunc(300*time.Second, FreshenUp)
}

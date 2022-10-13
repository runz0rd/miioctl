package zhimiairpmb4a

import (
	"encoding/json"
	"fmt"

	"github.com/runz0rd/miioctl/miio"
)

// https://home.miot-spec.com/spec/zhimi.airp.mb4a
type Device struct {
	*miio.Client
	IsOn             bool    `json:"is_on,omitempty"`
	Fault            string  `json:"fault,omitempty"` // device error
	Mode             float64 `json:"mode,omitempty"`
	PM25             float64 `json:"pm25,omitempty"`         // ug/m3
	FilterUsage      float64 `json:"filter_usage,omitempty"` // percent
	FilterUsageHours float64 `json:"filter_usage_hours,omitempty"`
	RPM              float64 `json:"rpm,omitempty"`
}

func New(addr, token string) (*Device, error) {
	dev := &Device{Client: miio.New(addr, token)}
	if err := dev.Query(); err != nil {
		return nil, err
	}
	return dev, nil
}

func (_ Device) DeviceId() string {
	return "zhimi.airp.mb4a"
}

func (d *Device) Query() error {
	resp, err := d.GetProperties(
		// https://home.miot-spec.com/spec/zhimi.airp.mb4a
		[]map[string]interface{}{
			{"siid": 2, "piid": 1},
			{"siid": 2, "piid": 2},
			{"siid": 2, "piid": 4},
			{"siid": 3, "piid": 4},
			{"siid": 4, "piid": 1},
			{"siid": 4, "piid": 3},
			{"siid": 9, "piid": 1},
		})
	if err != nil {
		return err
	}
	for i, result := range resp.Results {
		switch result.Siid {
		case 2:
			switch result.Piid {
			case 1:
				d.IsOn = resp.Results[i].Value.(bool)
			case 2:
				d.Fault = resp.Results[i].Value.(string)
			case 4:
				d.Mode = resp.Results[i].Value.(float64)
			}
		case 3:
			switch result.Piid {
			case 4:
				d.PM25 = resp.Results[i].Value.(float64)
			}
		case 4:
			switch result.Piid {
			case 1:
				d.FilterUsage = resp.Results[i].Value.(float64)
			case 3:
				d.FilterUsageHours = resp.Results[i].Value.(float64)
			}
		case 9:
			switch result.Piid {
			case 1:
				d.RPM = resp.Results[i].Value.(float64)
			}
		}
	}
	return nil
}

func (d *Device) SetPower(state bool) error {
	if _, err := d.Send("set_properties", []map[string]interface{}{{"siid": 2, "piid": 1, "value": state}}); err != nil {
		return err
	}
	d.IsOn = state
	return nil
}

func (d Device) TogglePower() error {
	state := true
	if d.IsOn {
		state = false
	}
	return d.SetPower(state)
}

func (d Device) ToString(field string) string {
	in, _ := json.Marshal(d)
	if field == "all" {
		return string(in)
	}
	var out map[string]interface{}
	json.Unmarshal(in, &out)
	return fmt.Sprint(out[field])
}

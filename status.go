package miioctl

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

type Status struct {
	Powered bool   `json:"powered,omitempty"`
	Aqi     int    `json:"aqi,omitempty"`    // μg/m³
	Mode    string `json:"mode,omitempty"`   // OperationMode.Favorite
	Filter  int    `json:"filter,omitempty"` // %
	Speed   int    `json:"speed,omitempty"`  // rpm
}

const (
	statusLines      = 28
	statusDebugLines = 54
)

func NewStatus(output string, debug bool) (*Status, error) {
	lines := strings.Split(output, "\n")
	if debug {
		if len(lines) != statusDebugLines {
			log.Debug(lines)
			return nil, errors.New("wrong output format")
		}
		lines = lines[statusDebugLines-statusLines:]
	} else {
		strings.Contains(lines[0], "WARNING")
		if len(lines) != statusLines {
			log.Debug(lines)
			return nil, errors.New("wrong output format")
		}
	}
	mapping, err := linesToMap(lines, ": ")
	if err != nil {
		return nil, err
	}
	powered := strings.Contains(mapping["Power"], "on")
	aqi, err := strconv.Atoi(strings.Replace(mapping["AQI"], " μg/m³", "", -1))
	if err != nil {
		return nil, err
	}
	mode := mapping["Mode"]
	filter, err := strconv.Atoi(strings.Replace(mapping["Filter life remaining"], " %", "", -1))
	if err != nil {
		return nil, err
	}
	speed, err := strconv.Atoi(strings.Replace(mapping["Motor speed"], " rpm", "", -1))
	if err != nil {
		return nil, err
	}
	return &Status{powered, aqi, mode, filter, speed}, nil
}

func (s Status) Get(field string) interface{} {
	if field == "all" {
		return spew.Sdump(s)
	}
	in, _ := json.Marshal(s)
	var out map[string]interface{}
	json.Unmarshal(in, &out)
	return out[field]
}


func linesToMap(lines []string, delim string) (map[string]string, error) {
	mapping := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		pieces := strings.Split(line, delim)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("cannot parse line %q", line)
		}
		mapping[pieces[0]] = pieces[1]
	}
	return mapping, nil
}

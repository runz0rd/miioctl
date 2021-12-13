package miioctl

import (
	"fmt"
	"strconv"
	"strings"
)

type Status struct {
	Powered bool
	Aqi     int    // μg/m³
	Mode    string // OperationMode.Favorite
	Filter  int    // %
	Speed   int    // rpm
}

const (
	statusLines      = 11
	statusDebugLines = 38
)

func NewStatus(output string, debug bool) (*Status, error) {
	lines := strings.Split(output, "\n")
	if debug {
		if len(lines) != statusDebugLines {
			return nil, fmt.Errorf("wrong output format")
		}
		lines = lines[statusDebugLines-statusLines:]
	} else {
		strings.Contains(lines[0], "WARNING")
		lines = lines[1:]
		if len(lines) != statusLines {
			return nil, fmt.Errorf("wrong output format")
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

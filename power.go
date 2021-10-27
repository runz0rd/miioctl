package miioctl

import "fmt"

type PowerCommand int

func NewPowerCommand(input string) (PowerCommand, error) {
	var pc PowerCommand
	switch input {
	case "on":
		pc = PowerOn
	case "off":
		pc = PowerOff
	case "toggle":
		pc = PowerToggle
	default:
		return pc, fmt.Errorf("invalid input %q for power command", input)
	}
	return pc, nil
}

const (
	PowerOn PowerCommand = iota
	PowerOff
	PowerToggle
)

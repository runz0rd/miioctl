package miioctl

import (
	"context"

	"os/exec"

	"github.com/pkg/errors"
)

type MiioCmd struct {
	name  string
	args  []string
	debug bool
}

func NewMiioCmd(device, ip, token string, debug bool) (*MiioCmd, error) {
	// check for miiocli
	_, err := exec.LookPath("miiocli")
	if err != nil {
		return nil, err
	}
	c := &MiioCmd{debug: debug}
	c.name = "miiocli"
	if debug {
		c.args = append(c.args, "--debug")
	}
	c.args = append(c.args, device, "--ip", ip, "--token", token)

	return c, nil
}

func (c MiioCmd) Info(ctx context.Context) *exec.Cmd {
	c.args = append(c.args, "info")
	return exec.CommandContext(ctx, c.name, c.args...)
}

func (c MiioCmd) Status(ctx context.Context) (*Status, error) {
	c.args = append(c.args, "status")
	out, err := exec.CommandContext(ctx, c.name, c.args...).CombinedOutput()
	if err != nil {
		return nil, errors.WithMessage(err, string(out))
	}
	return NewStatus(string(out), c.debug)
}

func (c MiioCmd) Power(ctx context.Context, pc PowerCommand) error {
	onCmd := exec.CommandContext(ctx, c.name, append(c.args, "on")...)
	offCmd := exec.CommandContext(ctx, c.name, append(c.args, "off")...)
	switch pc {
	case PowerOn:
		if out, err := onCmd.CombinedOutput(); err != nil {
			return errors.WithMessage(err, string(out))
		}
	case PowerOff:
		if out, err := offCmd.CombinedOutput(); err != nil {
			return errors.WithMessage(err, string(out))
		}
	case PowerToggle:
		status, err := c.Status(ctx)
		if err != nil {
			return err
		}
		if status.Powered {
			if out, err := onCmd.CombinedOutput(); err != nil {
				return errors.WithMessage(err, string(out))
			}
		} else {
			if out, err := offCmd.CombinedOutput(); err != nil {
				return errors.WithMessage(err, string(out))
			}
		}
	}
	return nil
}

func (c MiioCmd) IsPowered(ctx context.Context) (bool, error) {
	status, err := c.Status(ctx)
	if err != nil {
		return false, err
	}
	return status.Powered, nil
}

package miioctl

import (
	"context"

	"os/exec"

	cmd "github.com/foomo/gograpple/exec"
	"github.com/pkg/errors"
)

type MiioCmd struct {
	cmd.Cmd
}

func NewMiioCommand(device, ip, token string) (*MiioCmd, error) {
	// check for miiocli
	_, err := exec.LookPath("miiocli")
	if err != nil {
		return nil, err
	}
	return &MiioCmd{*cmd.NewCommand("miiocli").Args(device, "--ip", ip, "--token", token)}, nil
}

func (c MiioCmd) Info() *cmd.Cmd {
	return c.Args("info")
}

func (c MiioCmd) Status(ctx context.Context) (*Status, error) {
	out, err := c.Args("status").Run(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, out)
	}
	return NewStatus(out)
}

func (c MiioCmd) Power(ctx context.Context, pc PowerCommand) error {
	status, err := c.Status(ctx)
	if err != nil {
		return err
	}
	switch pc {
	case PowerOn:
		if !status.Powered {
			if out, err := c.Args("on").Run(ctx); err != nil {
				return errors.WithMessage(err, out)
			}
		}
	case PowerOff:
		if status.Powered {
			if out, err := c.Args("off").Run(ctx); err != nil {
				return errors.WithMessage(err, out)
			}
		}
	case PowerToggle:
		if status.Powered {
			if out, err := c.Args("off").Run(ctx); err != nil {
				return errors.WithMessage(err, out)
			}
		} else {
			if out, err := c.Args("on").Run(ctx); err != nil {
				return errors.WithMessage(err, out)
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

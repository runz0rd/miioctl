package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/runz0rd/miioctl"
)

func main() {
	var flagIp, flagToken, flagPower string
	var flagAqi, flagDebug bool
	flag.StringVar(&flagIp, "ip", "", "ip for miio device")
	flag.StringVar(&flagToken, "token", "", "token for miio device")
	flag.StringVar(&flagPower, "power", "", "power on/off/toggle")
	flag.BoolVar(&flagAqi, "aqi", false, "return current aqi readout")
	flag.BoolVar(&flagDebug, "debug", false, "debug mode")
	flag.Parse()

	if flagDebug {
		start := time.Now()
		defer func() {
			log.Printf("elapsed %v", time.Since(start))
		}()
	}
	if err := run(flagIp, flagToken, flagAqi, flagPower); err != nil {
		log.Fatal(err)
	}
}

func run(ip, token string, aqi bool, power string) error {
	ctx := context.Background()
	apctl, err := miio.NewMiioCommand("airpurifiermb4", ip, token)
	if err != nil {
		return err
	}
	if aqi {
		status, err := apctl.Status(ctx)
		if err != nil {
			return err
		}
		fmt.Print(status.Aqi)
		return nil
	}

	if power != "" {
		pc, err := miio.NewPowerCommand(power)
		if err != nil {
			return err
		}
		return apctl.Power(ctx, pc)
	}
	return nil
}

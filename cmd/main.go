package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/runz0rd/miioctl"
	"gopkg.in/yaml.v3"
)

func main() {
	var flagConfig, flagPower string
	var flagAqi, flagDebug bool
	flag.StringVar(&flagConfig, "config", "config.yaml", "config file path")
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
	if err := run(flagConfig, flagAqi, flagPower); err != nil {
		log.Fatal(err)
	}
}

func run(config string, aqi bool, power string) error {
	c, err := NewConfig(config)
	if err != nil {
		return err
	}
	ctx := context.Background()
	apctl, err := miioctl.NewMiioCmd("airpurifiermb4", c.Ip, c.Token, c.Debug)
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
		pc, err := miioctl.NewPowerCommand(power)
		if err != nil {
			return err
		}
		return apctl.Power(ctx, pc)
	}
	return nil
}

type Config struct {
	Ip    string `yaml:"ip,omitempty"`
	Token string `yaml:"token,omitempty"`
	Debug bool   `yaml:"debug,omitempty"`
}

func NewConfig(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := Config{}
	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

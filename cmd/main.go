package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/runz0rd/miioctl"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	var flagConfig, flagPower, flagStatus, flagServeAddr string
	var flagDebug bool
	var flagMin, flagMax int
	flag.StringVar(&flagConfig, "config", "config.yaml", "config file path")
	flag.StringVar(&flagPower, "power", "", "power on/off/toggle")
	flag.StringVar(&flagStatus, "status", "", "return current status readout")
	flag.StringVar(&flagServeAddr, "serve", "", "serve mode addr")
	flag.BoolVar(&flagDebug, "debug", false, "debug mode")
	flag.IntVar(&flagMin, "min", 0, "threshold to turn off")
	flag.IntVar(&flagMax, "max", 0, "threshold to turn on")
	flag.Parse()

	if flagDebug {
		log.SetLevel(log.DebugLevel)
	}
	start := time.Now()
	defer func() {
		log.Debugf("elapsed %v", time.Since(start))
	}()

	if err := run(flagConfig, flagStatus, flagServeAddr, flagPower, flagMin, flagMax); err != nil {
		log.Fatal(err)
	}
}

func run(config, status, serveAddr string, power string, tmin, tmax int) error {
	c, err := NewConfig(config)
	if err != nil {
		return err
	}
	ctx := context.Background()
	apctl, err := miioctl.NewMiioCmd("airpurifiermb4", c.Ip, c.Token, c.Debug)
	if err != nil {
		return err
	}
	if serveAddr != "" {
		r := mux.NewRouter()
		r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			s, err := apctl.Status(ctx)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error: %v", err)))
				return
			}
			w.Write([]byte(fmt.Sprint(s.Get("all"))))
		})
		r.HandleFunc("/power/{power}", func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			pc, err := miioctl.NewPowerCommand(vars["power"])
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error: %v", err)))
				return
			}
			if err := apctl.Power(ctx, pc); err != nil {
				w.Write([]byte(fmt.Sprintf("error: %v", err)))
				return
			}
		})
		log.Infof("serving on %v", serveAddr)
		log.Fatal(http.ListenAndServe(serveAddr, r))
	}
	if status != "" {
		log.Debug("status called")
		s, err := apctl.Status(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("%v", s.Get(status))
		return nil
	}

	if power != "" {
		log.Debug("power called")
		pc, err := miioctl.NewPowerCommand(power)
		if err != nil {
			return err
		}
		return apctl.Power(ctx, pc)
	}

	if tmin == 0 && tmax == 0 {
		return nil
	}
	if tmin >= tmax {
		return fmt.Errorf("min must be less than max")
	}
	if tmin > 0 || tmax > 0 {
		log.Debug("tmin or mtax called")
		status, err := apctl.Status(ctx)
		if err != nil {
			return err
		}
		if tmax > 0 && status.Aqi >= tmax {
			return apctl.Power(ctx, miioctl.PowerOn)
		}
		if tmin > 0 && status.Aqi <= tmin {
			return apctl.Power(ctx, miioctl.PowerOff)
		}
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/runz0rd/miioctl/miio/zhimiairpmb4a"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func main() {
	var flagConfig, flagPower, flagStatus, flagServeAddr string
	var flagDebug bool
	flag.StringVar(&flagConfig, "config", "config.yaml", "config file path")
	flag.StringVar(&flagPower, "power", "", "power on/off/toggle")
	flag.StringVar(&flagStatus, "status", "", "return current status readout")
	flag.StringVar(&flagServeAddr, "serve", "", "serve mode addr")
	flag.BoolVar(&flagDebug, "debug", false, "debug mode")
	flag.Parse()

	if flagDebug {
		log.SetLevel(log.DebugLevel)
	}
	start := time.Now()
	defer func() {
		log.Debugf("elapsed %v", time.Since(start))
	}()

	if err := run(flagConfig, flagStatus, flagServeAddr, flagPower); err != nil {
		log.Fatal(err)
	}
}

func run(config, status, serveAddr string, power string) error {
	c, err := NewConfig(config)
	if err != nil {
		return err
	}
	gatherer := zhimiairpmb4a.NewGatherer(c.Ip, c.Token, prometheus.NewRegistry())
	if serveAddr != "" {
		r := mux.NewRouter()
		r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			statusHandler(c, w, r)
		})
		r.HandleFunc("/power/{power}", func(w http.ResponseWriter, r *http.Request) {
			powerHandler(c, w, r)
		})
		r.Handle("/metrics", promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))

		log.Infof("serving on %v", serveAddr)
		log.Fatal(http.ListenAndServe(serveAddr, r))
	}

	device, err := zhimiairpmb4a.New(c.Ip, c.Token)
	if err != nil {
		return err
	}
	defer device.Close()
	if status != "" {
		log.Debug("status called")
		fmt.Printf("%v", device.ToString(status))
		return nil
	}

	if power != "" {
		log.Debug("power called")
		switch power {
		case "on":
			err = device.SetPower(true)
		case "off":
			err = device.SetPower(false)
		case "toggle":
			err = device.TogglePower()
		default:
			return fmt.Errorf("bad param: %v", power)
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

func statusHandler(c *Config, w http.ResponseWriter, r *http.Request) {
	device, err := zhimiairpmb4a.New(c.Ip, c.Token)
	if err != nil {
		panic(err)
	}
	defer device.Close()

	w.Write([]byte(device.ToString("all")))
}

func powerHandler(c *Config, w http.ResponseWriter, r *http.Request) {
	device, err := zhimiairpmb4a.New(c.Ip, c.Token)
	if err != nil {
		panic(err)
	}
	defer device.Close()

	switch mux.Vars(r)["power"] {
	case "on":
		err = device.SetPower(true)
	case "off":
		err = device.SetPower(false)
	case "toggle":
		err = device.TogglePower()
	default:
		w.Write([]byte(fmt.Sprintf("bad param: %v", mux.Vars(r)["power"])))
		return
	}
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error: %v", err)))
		return
	}
}

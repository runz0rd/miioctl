package zhimiairpmb4a

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	"github.com/runz0rd/miioctl/miio"
)

var powered = promauto.NewGauge(prometheus.GaugeOpts{
	Name:      "powered",
	Namespace: "miio_zhimiairpmb4a",
	Help:      "indicates powered state",
})

var pm25 = promauto.NewGauge(prometheus.GaugeOpts{
	Name:      "pm25",
	Namespace: "miio_zhimiairpmb4a",
	Help:      "indicates pm25 measurement",
})

var filter = promauto.NewGauge(prometheus.GaugeOpts{
	Name:      "filter",
	Namespace: "miio_zhimiairpmb4a",
	Help:      "indicates filter state",
})

type Gatherer struct {
	ip    string
	token string
	r     *prometheus.Registry
}

func NewGatherer(ip, token string, r *prometheus.Registry) *Gatherer {
	r.MustRegister(powered, pm25, filter)
	return &Gatherer{ip, token, r}
}

func (g *Gatherer) Gather() ([]*dto.MetricFamily, error) {
	client := miio.New(g.ip, g.token)
	defer client.Close()
	device, err := New(client)
	if err != nil {
		panic(err)
	}
	if device.IsOn {
		powered.Set(1)
	} else {
		powered.Set(0)
	}
	pm25.Set(device.PM25)
	filter.Set(device.FilterUsage)
	return g.r.Gather()
}

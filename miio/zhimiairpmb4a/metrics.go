package zhimiairpmb4a

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
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
	d *Device
	r *prometheus.Registry
}

func NewGatherer(d *Device, r *prometheus.Registry) *Gatherer {
	r.MustRegister(powered, pm25, filter)
	return &Gatherer{d, r}
}

func (g *Gatherer) Gather() ([]*dto.MetricFamily, error) {
	if err := g.d.Query(); err != nil {
		return nil, err
	}
	if g.d.IsOn {
		powered.Set(1)
	} else {
		powered.Set(0)
	}
	pm25.Set(g.d.PM25)
	filter.Set(g.d.FilterUsage)
	return g.r.Gather()
}

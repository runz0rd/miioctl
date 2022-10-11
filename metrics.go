package miioctl

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/version"
)

var powered = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "powered",
	Help: "indicates powered state",
})

var aqi = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "aqi",
	Help: "indicates aqi measurement",
})

var filter = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "filter",
	Help: "indicates filter state",
})

type StatusGatherer struct {
	miioCmd *MiioCmd
	r       *prometheus.Registry
}

func NewStatusGatherer(miioCmd *MiioCmd, r *prometheus.Registry) *StatusGatherer {
	r.MustRegister(version.NewCollector("airpurifiermiot_exporter"))
	r.MustRegister(powered, aqi, filter)
	return &StatusGatherer{miioCmd, r}
}

func (g *StatusGatherer) Gather() ([]*dto.MetricFamily, error) {
	s, err := g.miioCmd.Status(context.Background())
	if err != nil {
		return nil, err
	}
	if s.Powered {
		powered.Set(1)
	} else {
		powered.Set(0)
	}
	aqi.Set(float64(s.Aqi))
	filter.Set(float64(s.Filter))
	return g.r.Gather()
}

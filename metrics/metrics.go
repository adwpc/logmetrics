package metrics

import (
	"sync"

	"github.com/adwpc/logmetrics/model"
	"github.com/adwpc/logmetrics/prometheus/client_golang/prometheus"
	"github.com/adwpc/logmetrics/zlog"
)

var (
	log = zlog.Log
)

type PromInterface interface {
	Deal(float64, string)
}

type Counter struct {
	c           *prometheus.CounterVec
	firstReport bool
}

func (c *Counter) Deal(v float64, alert string) {
	if c.c != nil {
		c.c.With(prometheus.Labels{"alert": alert}).Add(v)
	}
}

type Gauge struct {
	g *prometheus.GaugeVec
}

func (g *Gauge) Deal(v float64, alert string) {
	if g.g != nil {
		g.g.With(prometheus.Labels{"alert": alert}).Set(v)
	}
}

type Histogram struct {
	h *prometheus.HistogramVec
}

func (h *Histogram) Deal(v float64, alert string) {
	if h.h != nil {
		h.h.With(prometheus.Labels{"alert": alert}).Observe(v)
	}
}

var (
	// pool  map[string]interface{}
	pool  map[string]PromInterface
	mutex sync.Mutex
)

func Get(key string, tp string, alert string) PromInterface {
	mutex.Lock()
	defer mutex.Unlock()
	if pool == nil {
		pool = make(map[string]PromInterface)
	}
	if _, ok := pool[key]; !ok {
		switch tp {
		case model.METRIC_COUNTER:
			pool[key] = NewCounter(key)
		case model.METRIC_GAUGE:
			pool[key] = NewGauge(key)
		case model.METRIC_HISTOGRAM:
			pool[key] = NewHistogram(key)
		default:
			log.Error().Msg("Get default")

		}
	}
	return pool[key]
}

func NewCounter(name string) *Counter {
	m := make(map[string]string)
	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        name,
			Help:        "CounterVec",
			ConstLabels: m,
		},
		[]string{"alert"},
	)
	if c == nil {
		log.Error().Msg("c == nil")
	}
	prometheus.MustRegister(c)

	return &Counter{
		c: c,
	}
}

func NewGauge(name string) *Gauge {
	m := make(map[string]string)
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        name,
			Help:        "GaugeVec",
			ConstLabels: m,
		},
		[]string{"alert"},
	)
	if g == nil {
		log.Error().Msg("c == nil")
	}
	prometheus.MustRegister(g)

	return &Gauge{
		g: g,
	}
}

func NewHistogram(name string) *Histogram {
	m := make(map[string]string)
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        name,
			Help:        "HistogramVec",
			ConstLabels: m,
		},
		[]string{"alert"},
	)
	if h == nil {
		log.Error().Msg("c == nil")
	}
	prometheus.MustRegister(h)

	return &Histogram{
		h: h,
	}
}

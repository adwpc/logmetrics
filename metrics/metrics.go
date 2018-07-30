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
	Deal(float64)
	SetAlert(string)
	SetFirstReport()
	SetFirstReportNext()
}

type Counter struct {
	c           prometheus.Counter
	firstReport bool
	val         []float64
}

func (c *Counter) SetFirstReport() {
	c.firstReport = true
	if len(c.val) > 0 {
		c.c.Add(c.val[0])
	}
}

func (c *Counter) SetFirstReportNext() {
	for i := 1; i < len(c.val); i++ {
		c.c.Add(c.val[i])
	}
}

func (c *Counter) Deal(v float64) {
	if c.c != nil {
		if c.firstReport {
			c.c.Add(v)
		} else {
			c.val = append(c.val, v)
		}

	}
}

func (c *Counter) SetAlert(a string) {
	if c.c != nil {
		c.c.SetLabel("alert", a)
	}
}

type Gauge struct {
	g prometheus.Gauge
}

func (g *Gauge) SetFirstReport() {
}

func (g *Gauge) SetFirstReportNext() {
}

func (g *Gauge) Deal(v float64) {
	if g.g != nil {
		g.g.Set(v)
	}
}

func (g *Gauge) SetAlert(a string) {
	if g.g != nil {
		g.g.SetLabel("alert", a)
	}
}

type Histogram struct {
	h prometheus.Histogram
}

func (h *Histogram) SetFirstReport() {
}

func (h *Histogram) SetFirstReportNext() {
}

func (h *Histogram) SetAlert(a string) {
	if h.h != nil {
		h.h.SetLabel("alert", a)
	}
}

func (h *Histogram) Deal(v float64) {
	if h.h != nil {
		h.h.Observe(v)
	}
}

var (
	// pool  map[string]interface{}
	pool  map[string]PromInterface
	mutex sync.Mutex
)

func FirstReport() {
	mutex.Lock()
	defer mutex.Unlock()
	for _, v := range pool {
		v.SetFirstReport()
	}
	log.Info().Msg("metrics.FirstReport")
}

func FirstReportNext() {
	mutex.Lock()
	defer mutex.Unlock()
	for _, v := range pool {
		v.SetFirstReportNext()
	}
	log.Info().Msg("metrics.FirstReportNext")
}

func Get(key string, tp string, alert string) PromInterface {
	mutex.Lock()
	defer mutex.Unlock()
	if pool == nil {
		pool = make(map[string]PromInterface)
	}
	if _, ok := pool[key]; !ok {
		switch tp {
		case model.METRIC_COUNTER:
			pool[key] = NewCounter(key, alert)
		case model.METRIC_GAUGE:
			pool[key] = NewGauge(key, alert)
		case model.METRIC_HISTOGRAM:
			pool[key] = NewHistogram(key, alert)
		default:
			log.Error().Msg("Get default")

		}
	}
	pool[key].SetAlert(alert)
	return pool[key]
}

func NewCounter(name string, alert string) *Counter {
	m := make(map[string]string)
	if alert != "" {
		m["alert"] = alert
	}
	c := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:        name,
			Help:        "help",
			ConstLabels: m,
		})
	if c == nil {
		log.Error().Msg("c == nil")
	}
	prometheus.MustRegister(c)

	return &Counter{
		c: c,
	}
}

func NewGauge(name string, alert string) *Gauge {
	m := make(map[string]string)
	if alert != "" {
		m["alert"] = alert
	}
	c := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name:        name,
			Help:        "help",
			ConstLabels: m,
		})
	if c == nil {
		log.Error().Msg("c == nil")
	}
	prometheus.MustRegister(c)

	return &Gauge{
		g: c,
	}
}

func NewHistogram(name string, alert string) *Histogram {
	m := make(map[string]string)
	if alert != "" {
		m["alert"] = alert
	}
	c := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:        name,
			Help:        "help",
			ConstLabels: m,
		})
	if c == nil {
		log.Error().Msg("c == nil")
	}
	prometheus.MustRegister(c)

	return &Histogram{
		h: c,
	}
}

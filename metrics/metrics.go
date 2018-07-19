package metrics

import (
	"sync"

	"github.com/adwpc/logmetrics/zlog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	log = zlog.Log
)

type PromInterface interface {
	Deal(float64)
}

type Counter struct {
	c prometheus.Counter
}

func (c *Counter) Deal(v float64) {
	if c.c != nil {
		c.c.Add(v)
	}
}

type Gauge struct {
	g prometheus.Gauge
}

func (g *Gauge) Deal(v float64) {
	if g.g != nil {
		g.g.Set(v)
	}
}

type Histogram struct {
	h prometheus.Histogram
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

func Get(key string, alert string) PromInterface {
	mutex.Lock()
	defer mutex.Unlock()
	if pool == nil {
		pool = make(map[string]PromInterface)
	}
	if _, ok := pool[key]; !ok {
		pool[key] = NewCounter(key, alert)
	}
	return pool[key]
}

// func NewCounter(namespace string, subSystem string, name string, help string) *prometheus.Counter {
func NewCounter(name string, alert string) *Counter {
	m := make(map[string]string)
	m["alert"] = alert
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

// func NewHistogram(histogram *prometheus.Histogram, namespace string,
// subSystem string, name string, help string) {

// *histogram = prometheus.NewHistogram(
// prometheus.HistogramOpts{
// Namespace: namespace,
// Subsystem: subSystem,
// Name:      name,
// Help:      help,
// Buckets:   DefBuckets,
// })
// prometheus.MustRegister(*histogram)
// }
type Metrics struct {
}

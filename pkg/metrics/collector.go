package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace      = "tp-link_router"
	labelClient    = "client"
	labelOperation = "operation"
)

type Collector struct {
	failure *prometheus.CounterVec
	rxLAN   *prometheus.GaugeVec
	txLAN   *prometheus.GaugeVec
	rxWAN   *prometheus.GaugeVec
	txWAN   *prometheus.GaugeVec
}

func New() *Collector {
	return &Collector{
		failure: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "failure_total",
				Help:      "Total number of failed operations",
			},
			[]string{labelOperation},
		),
		rxWAN: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "rx_lan_traffic_total",
				Help:      "",
			},
			[]string{labelClient},
		),
		txWAN: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "tx_lan_traffic_total",
				Help:      "",
			},
			[]string{labelClient},
		),
	}
}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, descs)
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	c.failure.Collect(metrics)
	c.rxWAN.Collect(metrics)
	c.txWAN.Collect(metrics)
}

func (c *Collector) Failure(operation string) {
	c.failure.With(prometheus.Labels{labelOperation: operation}).Inc()
}

func (c *Collector) RxLAN(client string, value float64) {
	c.rxWAN.With(prometheus.Labels{labelClient: client}).Set(value)
}

func (c *Collector) TxLAN(client string, value float64) {
	c.txWAN.With(prometheus.Labels{labelClient: client}).Set(value)
}

func (c *Collector) RxWAN(client string, value float64) {
	c.rxWAN.With(prometheus.Labels{labelClient: client}).Set(value)
}

func (c *Collector) TxWAN(client string, value float64) {
	c.txWAN.With(prometheus.Labels{labelClient: client}).Set(value)
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	ProductsCreated prometheus.Counter
	ProductsDeleted prometheus.Counter
}

// New registers product metrics on the given registerer.
func New(reg prometheus.Registerer) *Metrics {
	factory := promauto.With(reg)

	return &Metrics{
		ProductsCreated: factory.NewCounter(prometheus.CounterOpts{
			Namespace: "products",
			Name:      "created_total",
			Help:      "Total number of products created.",
		}),
		ProductsDeleted: factory.NewCounter(prometheus.CounterOpts{
			Namespace: "products",
			Name:      "deleted_total",
			Help:      "Total number of products deleted.",
		}),
	}
}

func NewNoop() *Metrics {
	return New(prometheus.NewRegistry())
}

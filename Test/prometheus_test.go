package Test

import (
	"github.com/peifengll/casbin_demo/rbac/cacheadapter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"testing"
	"time"
)

func TestLoadFrom(t *testing.T) {
	registry := prometheus.NewRegistry()

	registry.MustRegister(cacheadapter.GetCollectors()...)
	//registry.MustRegister(cacheadapter.LoadTimeHistogram)
	//cacheadapter.LoadFormCounter.WithLabelValues("mysql").
	go func() {
		for i := 0; i < 50; i++ {
			time.Sleep(1 * time.Second)
			TestLoad(t)
		}
	}()
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Registry: registry,
	}))
	http.ListenAndServe(":8080", nil)

}

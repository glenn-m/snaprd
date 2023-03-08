package main

import (
	"net/http"

	"github.com/alecthomas/kingpin/v2"

	"github.com/glenn-m/snaprd/internal/snaprd"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	configFile = kingpin.Flag("config",
		"The path to the snaprd configuration file.").Short('c').Default("snaprd.yaml").Envar("SNAPRD_CONFIG_FILE").String()
	metricsPort = kingpin.Flag("metricsPort",
		"The port the Prometheus metrics will be exposed on.").Short('p').Default("8080").Envar("SNAPRD_METRICS_PORT").String()
	metricsPath = kingpin.Flag("metricsPath",
		"The path the Prometheus metrics will be exposed on.").Short('m').Default("/metrics").Envar("SNARPD_METRICS_PATH").String()
)

func main() {
	// Parse in the flags
	kingpin.Parse()

	// Set timestamp in log format
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	srd, err := snaprd.New(*configFile)
	if err != nil {
		log.WithError(err).Fatal("error creating new snaprd instance...")
	}

	srd.Run()

	// Expose Prometheus Metrics
	http.Handle(*metricsPath, promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	// Setup listener for metrics
	log.Info("snaprd listening on " + *metricsPort)
	if err := http.ListenAndServe(":" + *metricsPort, nil); err != nil {
		log.WithError(err).Fatal("failed to listen on " + *metricsPort)
	}
}

package snaprd

import (
	"os"
	"strconv"
	"time"

	"github.com/glenn-m/snaprd/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	runFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "snaprd_run_failures",
		Help: "The number of failed snaprd runs",
	}, []string{"command"})
	runTime = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "snaprd_job_duration_seconds",
		Help: "The time taken to run a snaprd job",
	})
	numberRuns = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snaprd_run_total",
		Help: "The number of snaprd runs",
	})
	touchedFiles = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "snaprd_touched_file_count",
		Help: "The number of touched files in the current run",
	})
)

// Snaprd contains the configuration for snaprd
type Snaprd struct {
	Config   *config.Config
	LogFiles []string
}

// New creates a new instance of Snaprd
func New(configFile string) (*Snaprd, error) {
	config, err := config.Parse(configFile)
	if err != nil {
		return nil, err
	}
	snaprd := Snaprd{
		Config: config,
	}

	return &snaprd, nil
}

func (s *Snaprd) cleanup() {
	log.Info("Running cleanup...")
	for _, file := range s.LogFiles {
		log.Info(file)
		if err := os.Remove(file); err != nil {
			log.WithError(err).Info("error whilst running cleanup")
		}
	}
}

// Run starts the cron schedule
func (s *Snaprd) Run() {
	// Setup Cron scheduler
	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	c.AddFunc(s.Config.Schedule, func() {
		log.Info("Running scheduled snapraid")
		defer s.cleanup()
		st := time.Now()

		if s.Config.Snapraid.Touch {
			log.Info("Running touch...")
			touchOut, err := s.ExecCmd("touch")
			if err != nil {
				log.WithError(err).Error("error whilst running snapraid touch")
				runFailures.With(prometheus.Labels{"command": "touch"}).Inc()
				return
			}

			numTouched, err := s.ParseTouch(touchOut)
			if err != nil {
				log.WithError(err).Error("error whilst parsing touch logfile output")
				runFailures.With(prometheus.Labels{"command": "touch"}).Inc()
				return
			}

			touchedFiles.Set(numTouched)
			log.WithFields(log.Fields{"number": numTouched}).Info("Files touched...")
		}

		log.Info("Running diff...")
		diffOut, err := s.ExecCmd("diff")
		if err != nil {
			log.WithError(err).Error("error whilst running snapraid diff...")
			runFailures.With(prometheus.Labels{"command": "diff"}).Inc()
			return
		}

		log.Info("Checking if sync required...")
		syncRequired, err := s.ParseDiff(diffOut)
		if err != nil {
			log.WithError(err).Error("error whilst parsing touch logfile output...")
			runFailures.With(prometheus.Labels{"command": "diff"}).Inc()
			return
		}

		if syncRequired {
			log.Info("Running sync...")
			_, err := s.ExecCmd("sync")
			if err != nil {
				log.WithError(err).Error("error whilst running snapraid sync...")
				runFailures.With(prometheus.Labels{"command": "sync"}).Inc()
				return
			}
		}

		if s.Config.Scrub.Enabled {
			log.Info("Running scrub...")
			_, err := s.ExecCmd(
				"scrub",
				"-p",
				strconv.Itoa(s.Config.Scrub.Percentage),
				"-o",
				strconv.Itoa(s.Config.Scrub.OlderThan),
			)
			if err != nil {
				log.WithError(err).Error("error whilst running snapraid scrub...")
				runFailures.With(prometheus.Labels{"command": "scrub"}).Inc()
				return
			}
		}

		et := time.Since(st).Seconds()
		runTime.Set(et)
		numberRuns.Inc()
		log.WithFields(log.Fields{"Duration": et}).Info("Run complete")
	})

	c.Start()
}

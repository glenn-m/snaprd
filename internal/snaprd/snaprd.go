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

	s.LogFiles = []string{}
}

// Run starts the cron schedule
func (s *Snaprd) Run() {
	// Setup Cron scheduler
	c := cron.New(cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
	))

	_, err := c.AddFunc(s.Config.Schedule, func() {
		log.Info("Running scheduled snapraid")
		defer s.cleanup()
		st := time.Now()

		if s.Config.Snapraid.Touch {
			touchOut, err := s.runCmd("touch")
			if err != nil {
				return
			}

			numTouched, err := s.ParseTouch(touchOut)
			if err != nil {
				log.WithError(err).Error("error whilst parsing touch logfile output")
				// do you want to increment this as a runfailure again? maybe make a parseFailures metric
				// and use that instead
				// you could also move the logging and metrics to the parse methods directly to clean up the
				// logic in thi func.
				runFailures.With(prometheus.Labels{"command": "touch"}).Inc()
				return
			}

			touchedFiles.Set(numTouched)
			log.WithFields(log.Fields{"number": numTouched}).Info("Files touched...")
		}

		diffOut, err := s.runCmd("diff")
		if err != nil {
			return
		}

		log.Info("Checking if sync required...")
		diff, err := s.ParseDiff(diffOut)
		if err != nil {
			log.WithError(err).Error("error whilst parsing touch logfile output...")
			runFailures.With(prometheus.Labels{"command": "diff"}).Inc()
			return
		}

		if s.Config.Snapraid.DeleteThreshold > -1 {
			log.Info("Checking if diff exceeds delete threshold...")
			if diff.Removed >= float64(s.Config.Snapraid.DeleteThreshold) {
				log.WithFields(
					log.Fields{"removed": diff.Removed},
				).Error("aborting snapraid run, removed files would exceed delete threshold")
				runFailures.With(prometheus.Labels{"command": "diff"}).Inc()
				return
			}
		}

		if diff.SyncRequired {
			log.Info("Running sync...")
			_, err := s.runCmd("sync")
			if err != nil {
				return
			}
		}

		if s.Config.Scrub.Enabled {
			_, err := s.runCmd(
				"scrub",
				"-p",
				strconv.Itoa(s.Config.Scrub.Percentage),
				"-o",
				strconv.Itoa(s.Config.Scrub.OlderThan),
			)
			if err != nil {
				return
			}
		}

		et := time.Since(st).Seconds()
		runTime.Set(et)
		numberRuns.Inc()
		log.WithFields(log.Fields{"Duration": et}).Info("Run complete")
	})
	if err != nil {
		log.WithError(err).Error("error starting cron schedule")
	}

	c.Start()
}

// runCmd wraps ExecCmd to wrap logging and metrics
func (s *Snaprd) runCmd(cmd string, args ...string) (*os.File, error) {
	log.Infof("Running %s...", cmd)

	out, err := s.ExecCmd(cmd, args...)
	if err != nil {
		log.WithError(err).Errorf("error whilst running snapraid %s", cmd)
		runFailures.With(prometheus.Labels{"command": cmd}).Inc()
	}

	return out, err
}

package snaprd

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var diffStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "snaprd_diff_status",
	Help: "The number of files to be added, removed, updated on snapraid",
}, []string{"action"})

// Diff is a data struct from the diff log
type Diff struct {
	SyncRequired bool
	Equal        float64
	Added        float64
	Removed      float64
	Updated      float64
	Moved        float64
	Copied       float64
	Restored     float64
}

// ParseTouch parses the Touch cmd logfile output
func (s *Snaprd) ParseTouch(f *os.File) (float64, error) {
	var numTouched float64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "touch:") {
			numTouched++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0.0, err
	}

	return numTouched, nil
}

// ParseDiff parses the Diff cmd logfile output
func (s *Snaprd) ParseDiff(f *os.File) (*Diff, error) {
	var diff Diff
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "msg:status: There are differences!") {
			diff.SyncRequired = true
		}
		if strings.Contains(scanner.Text(), "summary:") {
			line := strings.Split(scanner.Text(), ":")
			if line[1] == "exit" {
				continue
			}

			// Convert to float64 for Prom metrics
			metricVal, err := strconv.ParseFloat(line[2], 64)
			if err != nil {
				return nil, err
			}

			switch line[1] {
			case "equal":
				diff.Equal = metricVal
				diffStatus.With(prometheus.Labels{"action": "equal"}).Set(metricVal)
			case "added":
				diff.Added = metricVal
				diffStatus.With(prometheus.Labels{"action": "added"}).Set(metricVal)
			case "removed":
				diff.Removed = metricVal
				diffStatus.With(prometheus.Labels{"action": "removed"}).Set(metricVal)
			case "updated":
				diff.Updated = metricVal
				diffStatus.With(prometheus.Labels{"action": "updated"}).Set(metricVal)
			case "moved":
				diff.Moved = metricVal
				diffStatus.With(prometheus.Labels{"action": "moved"}).Set(metricVal)
			case "copied":
				diff.Copied = metricVal
				diffStatus.With(prometheus.Labels{"action": "copied"}).Set(metricVal)
			case "restored":
				diff.Restored = metricVal
				diffStatus.With(prometheus.Labels{"action": "restored"}).Set(metricVal)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
  
	return &diff, nil
}

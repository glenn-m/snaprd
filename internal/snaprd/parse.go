package snaprd

import (
	"bufio"
	"os"
	"strings"
)

// ParseTouch parses the Touch cmd logfile output
func (s *Snaprd) ParseTouch(f *os.File) (float64, error) {
	scanner := bufio.NewScanner(f)
	numTouched := 0.0
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
func (s *Snaprd) ParseDiff(f *os.File) (bool, error) {
	scanner := bufio.NewScanner(f)
	var syncRequired bool
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "msg:status: There are differences!") {
			syncRequired = true
		}
	}
	return syncRequired, nil
}

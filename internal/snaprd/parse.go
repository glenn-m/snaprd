package snaprd

import (
	"os"
	"bufio"
	"strings"
)

// Parse the Touch cmd logfile output
func (s *Snaprd) ParseTouch(f *os.File) (int, error) {
	var numTouched int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "touch:") {
			numTouched++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return numTouched, nil
}

// Parse the Diff cmd logfile output
func (s *Snaprd) ParseDiff(f *os.File) (bool, error) {
	scanner := bufio.NewScanner(f)

	var syncRequired bool
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "msg:status: There are differences!") {
			syncRequired = true
		}
	}

	// I think you'll want an error
	err := scanner.Err()

	return syncRequired, err
}

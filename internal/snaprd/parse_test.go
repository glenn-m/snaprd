package snaprd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTouch(t *testing.T) {
	expected := 1.0
	f, err := os.Open("../../testdata/touch.log")
	if err != nil {
		t.Error(err)
	}
	s, err := New(configFile)
	if err != nil {
		t.Error(err)
	}
	actual, err := s.ParseTouch(f)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "They should be equal")
}

// TODO: Expand unit test to test metric values parsed from file
func TestParseDiff(t *testing.T) {
	expected := &Diff{
		SyncRequired: true,
		Equal:        4496,
		Added:        471,
		Removed:      157,
		Updated:      0,
		Moved:        0,
		Copied:       0,
		Restored:     0,
	}
	f, err := os.Open("../../testdata/diff.log")
	if err != nil {
		t.Error(err)
	}
	s, err := New(configFile)
	if err != nil {
		t.Error(err)
	}
	actual, err := s.ParseDiff(f)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "They should be equal")
}

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

func TestParseDiff(t *testing.T) {
	expected := true
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

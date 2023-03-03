package snaprd

import (
	"github.com/glenn-m/snaprd/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var configFile = "../../snaprd.yaml"

func TestNew(t *testing.T) {
	expected := Snaprd{
		Config: &config.Config{
			Schedule: "*/1 * * * *",
			Snapraid: config.Snapraid{
				Executable:      "snapraid",
				Config:          "snapraid.conf",
				DeleteThreshold: 40,
				Touch:           true,
			},
			Scrub: config.Scrub{
				Enabled:    false,
				Percentage: 12,
				OlderThan:  10,
			},
		},
	}
	actual, err := New(configFile)
	assert.Nil(t, err)
	assert.Equal(t, &expected, actual, "They should be equal")
}

func TestParseTouch(t *testing.T) {
	expected := 2.0
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

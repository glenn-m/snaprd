package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var configFile = "../../snaprd.example.yaml"

func TestParse(t *testing.T) {
	expected := Config{
		Schedule: "0 1 * * *",
		Snapraid: Snapraid{
			Executable:      "snapraid",
			ConfigPath:      "snapraid.conf",
			DeleteThreshold: 40,
			Touch:           true,
		},
		Scrub: Scrub{
			Enabled:    false,
			Percentage: 12,
			OlderThan:  10,
		},
	}
	actual, err := Parse(configFile)
	assert.Nil(t, err)
	assert.Equal(t, actual, &expected, "They should be equal")
}

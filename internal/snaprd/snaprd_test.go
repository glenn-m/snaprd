package snaprd

import (
	"testing"

	"github.com/glenn-m/snaprd/internal/config"
	"github.com/stretchr/testify/assert"
)

var configFile = "../../snaprd.example.yaml"

func TestNew(t *testing.T) {
	expected := Snaprd{
		Config: &config.Config{
			Schedule: "0 1 * * *",
			Snapraid: config.Snapraid{
				Executable:      "snapraid",
				ConfigPath:      "snapraid.conf",
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

package config_test

import (
	"testing"

	. "github.com/mfsyahrz/image_feed_api/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	t.Run("failed init config - file not found", func(t *testing.T) {
		newConfig, err := New("./random_path")
		assert.NotNil(t, err)
		assert.Nil(t, newConfig)
	})

	t.Run("failed init config - invalid file", func(t *testing.T) {
		newConfig, err := New("../../tests/config/env_invalid.test")
		assert.NotNil(t, err)
		assert.Nil(t, newConfig)
	})

	t.Run("success init config", func(t *testing.T) {
		wantConfig := &Config{
			Service: Service{
				Name: "imagefeed_api",
				Port: Port{
					REST: "8080",
				},
			},
		}

		newConfig, err := New("../../test/config/env.test")
		assert.Nil(t, err)
		assert.NotNil(t, wantConfig, newConfig)
	})
}

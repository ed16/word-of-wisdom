package config_test

import (
	"context"
	"os"
	"testing"

	"github.com/ed16/word-of-wisdom/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	ctx := context.Background()

	// Test with default values
	t.Run("default values", func(t *testing.T) {
		cfg, err := config.NewConfig(ctx, config.ServerConfig{})
		assert.NoError(t, err)
		assert.Equal(t, "0.0.0.0:8080", cfg.ListenAddr)
		assert.Equal(t, byte(5), cfg.Difficulty)
	})

	// Test with environment variables
	t.Run("environment variables", func(t *testing.T) {
		os.Setenv("SERVER_ADDR", "127.0.0.1:9090")
		os.Setenv("DIFFICULTY", "10")
		defer os.Unsetenv("SERVER_ADDR")
		defer os.Unsetenv("DIFFICULTY")

		cfg, err := config.NewConfig(ctx, config.ServerConfig{})
		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1:9090", cfg.ListenAddr)
		assert.Equal(t, byte(10), cfg.Difficulty)
	})

	// Test error handling
	t.Run("error handling", func(t *testing.T) {
		os.Setenv("DIFFICULTY", "asc") // intentionally breaking it
		defer os.Unsetenv("DIFFICULTY")

		var cfg config.ServerConfig
		_, err := config.NewConfig(ctx, cfg)
		require.Error(t, err)
	})
}

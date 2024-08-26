package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("should return an error when file not exits", func(t *testing.T) {
		// Action
		err := LoadConfig(".env")

		// Assert
		require.NotNil(t, err)
	})
	t.Run("should not return an error", func(t *testing.T) {
		// Arrange
		expected := "Go-Blog-App"

		// Action
		err := LoadConfig("../../.env")
		appName := Cfg.AppConfig.AppName

		// Assert
		require.Nil(t, err)
		require.Equal(t, expected, appName)
	})
}

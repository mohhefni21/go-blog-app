package database

import (
	"mohhefni/go-blog-app/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func init() {
	filepath := "../../.env"

	err := config.LoadConfig(filepath)
	if err != nil {
		panic(err)
	}
}

func TestDatabase(t *testing.T) {
	t.Run("should not return an error", func(t *testing.T) {
		// Action
		_, err := Connection(config.Cfg.DBconfig)

		// Assert
		require.Nil(t, err)
	})
}

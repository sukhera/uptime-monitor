package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sukhera/uptime-monitor/internal/api/routes"
	"github.com/sukhera/uptime-monitor/internal/config"
	"github.com/sukhera/uptime-monitor/internal/database"
	"github.com/sukhera/uptime-monitor/internal/models"
)

func TestAPIIntegration(t *testing.T) {
	// Setup test database
	cfg := &config.Config{
		MongoURI: "mongodb://localhost:27017",
		DBName:   "statuspage_test",
	}

	db, err := database.NewConnection(cfg.MongoURI, cfg.DBName)
	require.NoError(t, err)
	defer db.Close()

	// Setup test server
	router := routes.Setup(db)
	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("health check", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		assert.Equal(t, "healthy", result["status"])
	})

	t.Run("status endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/status")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result []models.ServiceStatus
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		// Should return empty array if no services configured
		assert.IsType(t, []models.ServiceStatus{}, result)
	})
}

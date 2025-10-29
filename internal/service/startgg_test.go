package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStartGGClient(t *testing.T) {
	// Set up environment variable for test
	originalToken := os.Getenv("START_GG_API_KEY")
	defer func() {
		if originalToken != "" {
			os.Setenv("START_GG_API_KEY", originalToken)
		} else {
			os.Unsetenv("START_GG_API_KEY")
		}
	}()

	os.Setenv("START_GG_API_KEY", "test-token")

	client, err := NewStartGGClient()
	require.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "test-token", client.apiToken)
	assert.NotNil(t, client.client)
	assert.Equal(t, 30*time.Second, client.client.Timeout)
}

func TestNewStartGGClient_NoAPIKey(t *testing.T) {
	// Remove environment variable
	originalToken := os.Getenv("START_GG_API_KEY")
	os.Unsetenv("START_GG_API_KEY")
	defer func() {
		if originalToken != "" {
			os.Setenv("START_GG_API_KEY", originalToken)
		}
	}()

	client, err := NewStartGGClient()
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.Contains(t, err.Error(), "START_GG_API_KEY environment variable not set")
}

func TestStartGGClient_GetUpcomingTournaments(t *testing.T) {
	// Set up environment variable for test
	originalToken := os.Getenv("START_GG_API_KEY")
	defer func() {
		if originalToken != "" {
			os.Setenv("START_GG_API_KEY", originalToken)
		} else {
			os.Unsetenv("START_GG_API_KEY")
		}
	}()
	os.Setenv("START_GG_API_KEY", "test-token")

	// Mock server response
	mockResponse := map[string]interface{}{
		"data": map[string]interface{}{
			"tournaments": map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"id":        "tournament-1",
						"name":      "Test Tournament 1",
						"slug":      "test-tournament-1",
						"startAt":   float64(1640995200),
						"endAt":     float64(1641081600),
						"venueName": "Test Venue",
						"city":      "Test City",
						"state":     "Test State",
						"country":   "Test Country",
						"events": []interface{}{
							map[string]interface{}{
								"id":          "event-1",
								"name":        "Singles",
								"slug":        "singles",
								"numEntrants": float64(64),
								"videogame": map[string]interface{}{
									"id":   "game-1",
									"name": "Super Smash Bros. Ultimate",
								},
							},
						},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := &StartGGClient{
		apiToken: "test-token",
		client:   server.Client(),
	}

	originalURL := StartGGAPIURL
	StartGGAPIURL = server.URL
	defer func() { StartGGAPIURL = originalURL }()

	tournaments, err := client.GetUpcomingTournaments(5)
	require.NoError(t, err)
	require.Len(t, tournaments, 1)

	t1 := tournaments[0]
	assert.Equal(t, "tournament-1", t1.ID)
	assert.Equal(t, "Test Tournament 1", t1.Name)
	assert.Equal(t, "test-tournament-1", t1.Slug)
	assert.Equal(t, int64(1640995200), t1.StartAt)
	assert.Equal(t, int64(1641081600), t1.EndAt)
	assert.Equal(t, "Test Venue", *t1.VenueName)
	assert.Equal(t, "Test City", *t1.City)
	assert.Equal(t, "Test State", *t1.State)
	assert.Equal(t, "Test Country", *t1.Country)

	require.Len(t, t1.Events, 1)
	event := t1.Events[0]
	assert.Equal(t, "event-1", event.ID)
	assert.Equal(t, "Singles", event.Name)
	assert.Equal(t, "singles", event.Slug)
	assert.Equal(t, 64, *event.NumEntrants)
	assert.Equal(t, "game-1", event.Videogame.ID)
	assert.Equal(t, "Super Smash Bros. Ultimate", event.Videogame.Name)
}

func TestStartGGClient_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid API token"}`))
	}))
	defer server.Close()

	client := &StartGGClient{
		apiToken: "invalid-token",
		client:   server.Client(),
	}

	originalURL := StartGGAPIURL
	StartGGAPIURL = server.URL
	defer func() { StartGGAPIURL = originalURL }()

	_, err := client.GetUpcomingTournaments(5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed with status 401")
}

func TestStartGGClient_GraphQLError(t *testing.T) {
	mockResponse := map[string]interface{}{
		"errors": []interface{}{
			map[string]interface{}{
				"message": "Invalid query",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := &StartGGClient{
		apiToken: "test-token",
		client:   server.Client(),
	}

	originalURL := StartGGAPIURL
	StartGGAPIURL = server.URL
	defer func() { StartGGAPIURL = originalURL }()

	_, err := client.GetUpcomingTournaments(5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GraphQL errors")
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type StartGGService struct {
	apiToken string
	endpoint string
	client   *http.Client
}

func NewStartGGService() *StartGGService {
	return &StartGGService{
		apiToken: os.Getenv("STARTGG_API_TOKEN"),
		endpoint: "https://api.start.gg/gql/alpha",
		client:   &http.Client{},
	}
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{}            `json:"data"`
	Errors []GraphQLError         `json:"errors,omitempty"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

// QueryTournaments fetches tournament data from Start.GG
func (s *StartGGService) QueryTournaments(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", s.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	var result GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", result.Errors)
	}

	return &result, nil
}

// Example query to get user's tournament placements
const UserTournamentsQuery = `
query UserTournaments($userId: ID!, $perPage: Int!) {
  user(id: $userId) {
    player {
      sets(perPage: $perPage) {
        nodes {
          event {
            tournament {
              name
              startAt
            }
            name
          }
          displayScore
          fullRoundText
        }
      }
    }
  }
}
`

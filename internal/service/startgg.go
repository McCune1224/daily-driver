package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var StartGGAPIURL = "https://api.start.gg/gql/alpha"

type StartGGClient struct {
	apiToken string
	client   *http.Client
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

type Tournament struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	StartAt   int64   `json:"startAt"`
	EndAt     int64   `json:"endAt"`
	VenueName *string `json:"venueName"`
	City      *string `json:"city"`
	State     *string `json:"state"`
	Country   *string `json:"country"`
	Events    []Event `json:"events"`
}

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	NumEntrants *int   `json:"numEntrants"`
	Videogame   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"videogame"`
}

type TournamentResult struct {
	Tournament Tournament `json:"tournament"`
	Placement  int        `json:"placement"`
	Seed       int        `json:"seed"`
	Entrants   int        `json:"entrants"`
	Sets       []Set      `json:"sets"`
}

type Set struct {
	ID           string `json:"id"`
	DisplayScore string `json:"displayScore"`
	WinnerID     string `json:"winnerId"`
	Slots        []Slot `json:"slots"`
}

type Slot struct {
	Entrant *Entrant `json:"entrant"`
}

type Entrant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID     string `json:"id"`
	Slug   string `json:"slug"`
	Player Player `json:"player"`
}

type Player struct {
	GamerTag string `json:"gamerTag"`
	Prefix   string `json:"prefix"`
}

func NewStartGGClient() (*StartGGClient, error) {
	apiToken := os.Getenv("START_GG_API_KEY")
	if apiToken == "" {
		return nil, fmt.Errorf("START_GG_API_KEY environment variable not set")
	}

	return &StartGGClient{
		apiToken: apiToken,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *StartGGClient) makeRequest(query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", StartGGAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", gqlResp.Errors)
	}

	return &gqlResp, nil
}

// GetUpcomingTournaments fetches upcoming tournaments
func (c *StartGGClient) GetUpcomingTournaments(limit int) ([]Tournament, error) {
	query := `
		query GetUpcomingTournaments($perPage: Int) {
			tournaments(query: {
				perPage: $perPage
				filter: {
					upcoming: true
				}
			}) {
				nodes {
					id
					name
					slug
					startAt
					endAt
					venueName
					city
					state
					country
					events {
						id
						name
						slug
						numEntrants
						videogame {
							id
							name
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"perPage": limit,
	}

	resp, err := c.makeRequest(query, variables)
	if err != nil {
		return nil, err
	}

	// Parse the response
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response structure")
	}

	tournamentsData, ok := data["tournaments"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("tournaments field not found")
	}

	nodes, ok := tournamentsData["nodes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("nodes field not found")
	}

	var tournaments []Tournament
	for _, node := range nodes {
		nodeData, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		tournament := Tournament{}
		if id, ok := nodeData["id"].(string); ok {
			tournament.ID = id
		}
		if name, ok := nodeData["name"].(string); ok {
			tournament.Name = name
		}
		if slug, ok := nodeData["slug"].(string); ok {
			tournament.Slug = slug
		}
		if startAt, ok := nodeData["startAt"].(float64); ok {
			tournament.StartAt = int64(startAt)
		}
		if endAt, ok := nodeData["endAt"].(float64); ok {
			tournament.EndAt = int64(endAt)
		}
		if venueName, ok := nodeData["venueName"].(string); ok {
			tournament.VenueName = &venueName
		}
		if city, ok := nodeData["city"].(string); ok {
			tournament.City = &city
		}
		if state, ok := nodeData["state"].(string); ok {
			tournament.State = &state
		}
		if country, ok := nodeData["country"].(string); ok {
			tournament.Country = &country
		}

		// Parse events
		if eventsData, ok := nodeData["events"].([]interface{}); ok {
			for _, eventData := range eventsData {
				eventMap, ok := eventData.(map[string]interface{})
				if !ok {
					continue
				}

				event := Event{}
				if id, ok := eventMap["id"].(string); ok {
					event.ID = id
				}
				if name, ok := eventMap["name"].(string); ok {
					event.Name = name
				}
				if slug, ok := eventMap["slug"].(string); ok {
					event.Slug = slug
				}
				if numEntrants, ok := eventMap["numEntrants"].(float64); ok {
					num := int(numEntrants)
					event.NumEntrants = &num
				}

				if videogameData, ok := eventMap["videogame"].(map[string]interface{}); ok {
					if id, ok := videogameData["id"].(string); ok {
						event.Videogame.ID = id
					}
					if name, ok := videogameData["name"].(string); ok {
						event.Videogame.Name = name
					}
				}

				tournament.Events = append(tournament.Events, event)
			}
		}

		tournaments = append(tournaments, tournament)
	}

	return tournaments, nil
}

// GetTournamentBySlug fetches a specific tournament by its slug
func (c *StartGGClient) GetTournamentBySlug(slug string) (*Tournament, error) {
	query := `
		query GetTournament($slug: String) {
			tournament(slug: $slug) {
				id
				name
				slug
				startAt
				endAt
				venueName
				city
				state
				country
				events {
					id
					name
					slug
					numEntrants
					videogame {
						id
						name
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"slug": slug,
	}

	resp, err := c.makeRequest(query, variables)
	if err != nil {
		return nil, err
	}

	// Parse the response
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response structure")
	}

	tournamentData, ok := data["tournament"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("tournament field not found")
	}

	tournament := &Tournament{}
	if id, ok := tournamentData["id"].(string); ok {
		tournament.ID = id
	}
	if name, ok := tournamentData["name"].(string); ok {
		tournament.Name = name
	}
	if slug, ok := tournamentData["slug"].(string); ok {
		tournament.Slug = slug
	}
	if startAt, ok := tournamentData["startAt"].(float64); ok {
		tournament.StartAt = int64(startAt)
	}
	if endAt, ok := tournamentData["endAt"].(float64); ok {
		tournament.EndAt = int64(endAt)
	}
	if venueName, ok := tournamentData["venueName"].(string); ok {
		tournament.VenueName = &venueName
	}
	if city, ok := tournamentData["city"].(string); ok {
		tournament.City = &city
	}
	if state, ok := tournamentData["state"].(string); ok {
		tournament.State = &state
	}
	if country, ok := tournamentData["country"].(string); ok {
		tournament.Country = &country
	}

	// Parse events (similar to above)
	if eventsData, ok := tournamentData["events"].([]interface{}); ok {
		for _, eventData := range eventsData {
			eventMap, ok := eventData.(map[string]interface{})
			if !ok {
				continue
			}

			event := Event{}
			if id, ok := eventMap["id"].(string); ok {
				event.ID = id
			}
			if name, ok := eventMap["name"].(string); ok {
				event.Name = name
			}
			if slug, ok := eventMap["slug"].(string); ok {
				event.Slug = slug
			}
			if numEntrants, ok := eventMap["numEntrants"].(float64); ok {
				num := int(numEntrants)
				event.NumEntrants = &num
			}

			if videogameData, ok := eventMap["videogame"].(map[string]interface{}); ok {
				if id, ok := videogameData["id"].(string); ok {
					event.Videogame.ID = id
				}
				if name, ok := videogameData["name"].(string); ok {
					event.Videogame.Name = name
				}
			}

			tournament.Events = append(tournament.Events, event)
		}
	}

	return tournament, nil
}

// GetUserByGamerTag searches for a user by their gamer tag
func (c *StartGGClient) GetUserByGamerTag(gamerTag string) (*User, error) {
	query := `
		query GetUserByGamerTag($gamerTag: String) {
			user(slug: $gamerTag) {
				id
				slug
				player {
					gamerTag
					prefix
				}
			}
		}
	`

	variables := map[string]interface{}{
		"gamerTag": gamerTag,
	}

	resp, err := c.makeRequest(query, variables)
	if err != nil {
		return nil, err
	}

	// Parse the response
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response structure")
	}

	userData, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("user field not found")
	}

	user := &User{}

	if id, ok := userData["id"].(string); ok {
		user.ID = id
	}
	if slug, ok := userData["slug"].(string); ok {
		user.Slug = slug
	}

	if playerData, ok := userData["player"].(map[string]interface{}); ok {
		player := Player{}
		if gamerTag, ok := playerData["gamerTag"].(string); ok {
			player.GamerTag = gamerTag
		}
		if prefix, ok := playerData["prefix"].(string); ok {
			player.Prefix = prefix
		}
		user.Player = player
	}

	return user, nil
}

// SearchUsers searches for users by gamer tag (returns multiple results)
func (c *StartGGClient) SearchUsers(gamerTag string, limit int) ([]User, error) {
	query := `
		query SearchUsers($query: String, $perPage: Int) {
			users(query: {
				search: {
					fields: GAMER_TAG
					query: $query
				}
			}, perPage: $perPage) {
				nodes {
					id
					slug
					player {
						gamerTag
						prefix
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"query":   gamerTag,
		"perPage": limit,
	}

	resp, err := c.makeRequest(query, variables)
	if err != nil {
		return nil, err
	}

	// Parse the response
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response structure")
	}

	usersData, ok := data["users"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("users field not found")
	}

	nodes, ok := usersData["nodes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("nodes field not found")
	}

	var users []User
	for _, node := range nodes {
		nodeData, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		user := User{}

		if id, ok := nodeData["id"].(string); ok {
			user.ID = id
		}
		if slug, ok := nodeData["slug"].(string); ok {
			user.Slug = slug
		}

		if playerData, ok := nodeData["player"].(map[string]interface{}); ok {
			player := Player{}
			if gamerTag, ok := playerData["gamerTag"].(string); ok {
				player.GamerTag = gamerTag
			}
			if prefix, ok := playerData["prefix"].(string); ok {
				player.Prefix = prefix
			}
			user.Player = player
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserTournamentHistory fetches tournament history for a specific user
func (c *StartGGClient) GetUserTournamentHistory(userSlug string, limit int) ([]TournamentResult, error) {
	query := `
		query GetUserTournamentHistory($slug: String, $perPage: Int) {
			user(slug: $slug) {
				player {
					gamerTag
					prefix
				}
				recentStandings(perPage: $perPage) {
					nodes {
						entrant {
							id
							name
							participants {
								gamerTag
								prefix
							}
						}
						standing {
							placement
							seed
						}
						tournament {
							id
							name
							slug
							startAt
							endAt
							venueAddress
							city
							addrState
							countryCode
							events {
								id
								name
								slug
								numEntrants
								videogame {
									id
									name
								}
								sets(perPage: 10) {
									nodes {
										id
										displayScore
										winnerId
										slots {
											entrant {
												id
												name
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"slug":    userSlug,
		"perPage": limit,
	}

	resp, err := c.makeRequest(query, variables)
	if err != nil {
		return nil, err
	}

	// Parse the response
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response structure")
	}

	userData, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("user field not found")
	}

	recentStandings, ok := userData["recentStandings"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("recentStandings field not found")
	}

	nodes, ok := recentStandings["nodes"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("nodes field not found")
	}

	var results []TournamentResult
	for _, node := range nodes {
		nodeData, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		result := TournamentResult{}

		// Parse standing information
		if standingData, ok := nodeData["standing"].(map[string]interface{}); ok {
			if placement, ok := standingData["placement"].(float64); ok {
				result.Placement = int(placement)
			}
			if seed, ok := standingData["seed"].(float64); ok {
				result.Seed = int(seed)
			}
		}

		// Parse tournament information
		if tournamentData, ok := nodeData["tournament"].(map[string]interface{}); ok {
			tournament := Tournament{}

			if id, ok := tournamentData["id"].(string); ok {
				tournament.ID = id
			}
			if name, ok := tournamentData["name"].(string); ok {
				tournament.Name = name
			}
			if slug, ok := tournamentData["slug"].(string); ok {
				tournament.Slug = slug
			}
			if startAt, ok := tournamentData["startAt"].(float64); ok {
				tournament.StartAt = int64(startAt)
			}
			if endAt, ok := tournamentData["endAt"].(float64); ok {
				tournament.EndAt = int64(endAt)
			}
			if venueAddress, ok := tournamentData["venueAddress"].(string); ok {
				tournament.VenueName = &venueAddress
			}
			if city, ok := tournamentData["city"].(string); ok {
				tournament.City = &city
			}
			if addrState, ok := tournamentData["addrState"].(string); ok {
				tournament.State = &addrState
			}
			if countryCode, ok := tournamentData["countryCode"].(string); ok {
				tournament.Country = &countryCode
			}

			// Parse events
			if eventsData, ok := tournamentData["events"].([]interface{}); ok {
				for _, eventData := range eventsData {
					eventMap, ok := eventData.(map[string]interface{})
					if !ok {
						continue
					}

					event := Event{}
					if id, ok := eventMap["id"].(string); ok {
						event.ID = id
					}
					if name, ok := eventMap["name"].(string); ok {
						event.Name = name
					}
					if slug, ok := eventMap["slug"].(string); ok {
						event.Slug = slug
					}
					if numEntrants, ok := eventMap["numEntrants"].(float64); ok {
						num := int(numEntrants)
						event.NumEntrants = &num
						result.Entrants = num // Use the largest entrant count
					}

					if videogameData, ok := eventMap["videogame"].(map[string]interface{}); ok {
						if id, ok := videogameData["id"].(string); ok {
							event.Videogame.ID = id
						}
						if name, ok := videogameData["name"].(string); ok {
							event.Videogame.Name = name
						}
					}

					// Parse sets for this event
					if setsData, ok := eventMap["sets"].(map[string]interface{}); ok {
						if setNodes, ok := setsData["nodes"].([]interface{}); ok {
							for _, setNode := range setNodes {
								setMap, ok := setNode.(map[string]interface{})
								if !ok {
									continue
								}

								set := Set{}
								if id, ok := setMap["id"].(string); ok {
									set.ID = id
								}
								if displayScore, ok := setMap["displayScore"].(string); ok {
									set.DisplayScore = displayScore
								}
								if winnerId, ok := setMap["winnerId"].(string); ok {
									set.WinnerID = winnerId
								}

								// Parse slots
								if slotsData, ok := setMap["slots"].([]interface{}); ok {
									for _, slotData := range slotsData {
										slotMap, ok := slotData.(map[string]interface{})
										if !ok {
											continue
										}

										slot := Slot{}
										if entrantData, ok := slotMap["entrant"].(map[string]interface{}); ok {
											entrant := Entrant{}
											if id, ok := entrantData["id"].(string); ok {
												entrant.ID = id
											}
											if name, ok := entrantData["name"].(string); ok {
												entrant.Name = name
											}
											slot.Entrant = &entrant
										}
										set.Slots = append(set.Slots, slot)
									}
								}

								result.Sets = append(result.Sets, set)
							}
						}
					}

					tournament.Events = append(tournament.Events, event)
				}
			}

			result.Tournament = tournament
		}

		results = append(results, result)
	}

	return results, nil
}

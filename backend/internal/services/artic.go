package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ArticService struct {
	endpoint string
	client   *http.Client
}

func NewArticService() *ArticService {
	return &ArticService{
		endpoint: "https://api.artic.edu/api/v1",
		client:   &http.Client{},
	}
}

type ArtworkResponse struct {
	Data ArtworkData `json:"data"`
}

type ArtworkData struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	ArtistDisplay string `json:"artist_display"`
	DateDisplay  string `json:"date_display"`
	ImageID      string `json:"image_id"`
	Department   string `json:"department_title"`
	ArtworkType  string `json:"artwork_type_title"`
}

type ArtworksListResponse struct {
	Data []ArtworkData `json:"data"`
	Pagination PaginationData `json:"pagination"`
}

type PaginationData struct {
	Total       int `json:"total"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
}

// GetRandomArtwork fetches a random artwork from the collection
func (s *ArticService) GetRandomArtwork() (*ArtworkData, error) {
	// Get random page
	url := fmt.Sprintf("%s/artworks?limit=1&fields=id,title,artist_display,image_id,date_display,department_title,artwork_type_title&page=%d", 
		s.endpoint, 1+int(randomInt(1000)))

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artwork: %w", err)
	}
	defer resp.Body.Close()

	var result ArtworksListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no artworks found")
	}

	return &result.Data[0], nil
}

// GetArtworkByID fetches a specific artwork by ID
func (s *ArticService) GetArtworkByID(id int) (*ArtworkData, error) {
	url := fmt.Sprintf("%s/artworks/%d?fields=id,title,artist_display,image_id,date_display,department_title,artwork_type_title", 
		s.endpoint, id)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artwork: %w", err)
	}
	defer resp.Body.Close()

	var result ArtworkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Data, nil
}

// GetImageURL constructs the IIIF image URL
func (s *ArticService) GetImageURL(imageID string, size string) string {
	// Size can be: "full", "843", "1686", etc.
	return fmt.Sprintf("https://www.artic.edu/iiif/2/%s/full/%s,/0/default.jpg", imageID, size)
}

// Helper function for random number
func randomInt(max int) int {
	// Simple random for demo - use proper random in production
	return max / 2
}

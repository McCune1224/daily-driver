package art

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

const ChicagoArtInstituteAPIURL = "https://api.artic.edu/api/v1/artworks"

type ChicagoArtAPIHandler struct {
	BaseURL string
	Client  *http.Client
}

func NewChicagoAPIClient() *ChicagoArtAPIHandler {
	return &ChicagoArtAPIHandler{
		BaseURL: ChicagoArtInstituteAPIURL,
		Client: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
	}
}

// Structs for API response
type ArtworkResponse struct {
	Data       []Artwork  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type ArtworkDetailResponse struct {
	Data Artwork `json:"data"`
}

type Artwork struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	ArtistDisplay string `json:"artist_display"`
	DateDisplay   string `json:"date_display"`
	MediumDisplay string `json:"medium_display"`
	ImageID       string `json:"image_id"`
	Thumbnail     struct {
		AltText string `json:"alt_text"`
	} `json:"thumbnail"`
}

func (a *Artwork) ImageURL() string {
	if a.ImageID == "" {
		return ""
	}
	return fmt.Sprintf("https://www.artic.edu/iiif/2/%s/full/843,/0/default.jpg", a.ImageID)
}

type Pagination struct {
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
}

func (api *ChicagoArtAPIHandler) GetRandomArtwork() (*Artwork, error) {
	// Method 1: Get total count and fetch random page
	totalCount, err := api.GetTotalArtworkCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Calculate random page (assuming 100 items per page)
	itemsPerPage := 100
	totalPages := (totalCount + itemsPerPage - 1) / itemsPerPage
	randomPage := rand.Intn(totalPages) + 1

	// Fetch random page
	url := fmt.Sprintf("%s?page=%d&limit=%d", ChicagoArtInstituteAPIURL, randomPage, itemsPerPage)
	artworks, err := api.FetchArtworks(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artworks: %w", err)
	}

	if len(artworks) == 0 {
		return nil, fmt.Errorf("no artworks found")
	}

	// Pick random artwork from the page
	randomIndex := rand.Intn(len(artworks))
	selectedArtwork := artworks[randomIndex]

	// Get detailed information
	detailedArtwork, err := api.GetArtworkDetails(selectedArtwork.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get artwork details: %w", err)
	}

	return detailedArtwork, nil
}

func (api *ChicagoArtAPIHandler) GetTotalArtworkCount() (int, error) {
	url := fmt.Sprintf("%s?limit=1", ChicagoArtInstituteAPIURL)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var artworkResp ArtworkResponse
	if err := json.Unmarshal(body, &artworkResp); err != nil {
		return 0, err
	}

	return artworkResp.Pagination.Total, nil
}

func (api *ChicagoArtAPIHandler) FetchArtworks(url string) ([]Artwork, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artworkResp ArtworkResponse
	if err := json.Unmarshal(body, &artworkResp); err != nil {
		return nil, err
	}

	return artworkResp.Data, nil
}

func (api *ChicagoArtAPIHandler) GetArtworkDetails(artworkID int) (*Artwork, error) {
	fields := "id,title,artist_display,date_display,medium_display,image_id,thumbnail"
	url := fmt.Sprintf("%s/%d?fields=%s", ChicagoArtInstituteAPIURL, artworkID, fields)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var detailResp ArtworkDetailResponse
	if err := json.Unmarshal(body, &detailResp); err != nil {
		return nil, err
	}

	return &detailResp.Data, nil
}

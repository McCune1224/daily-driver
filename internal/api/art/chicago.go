package art

import (
	"encoding/json"
	"net/http"
)

type ChicagoArtworkResponse struct {
	Data []struct {
		Title     string `json:"title"`
		Artist    string `json:"artist_display"`
		Thumbnail struct {
			URL string `json:"lqip"`
		} `json:"thumbnail"`
		ImageID string `json:"image_id"`
		IIIFURL string `json:"iiif_url"`
	} `json:"data"`
}

type ChicagoArtInstituteAPI struct {
	BaseURL string
	Client  *http.Client
}

func NewChicagoAPIClient() *ChicagoArtInstituteAPI {
	return &ChicagoArtInstituteAPI{
		BaseURL: "https://api.artic.edu/api/v1",
		Client:  &http.Client{},
	}
}

func (api *ChicagoArtInstituteAPI) GetArtworks() (*http.Response, error) {
	req, err := http.NewRequest("GET", api.BaseURL+"/artworks", nil)
	if err != nil {
		return nil, err
	}
	return api.Client.Do(req)
}

func (api *ChicagoArtInstituteAPI) GetRandomArtwork() (*ChicagoArtworkResponse, error) {
	url := api.BaseURL + "/artworks"

	// Fetching one random artwork by adding random parameters or logic
	url = url + "?limit=1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artworkResp ChicagoArtworkResponse
	if err = json.NewDecoder(resp.Body).Decode(&artworkResp); err != nil {
		return nil, err
	}

	// Construct the true image URL
	if len(artworkResp.Data) > 0 {
		artwork := &artworkResp.Data[0]
		artwork.IIIFURL = artwork.IIIFURL + "/full/full/0/default.jpg"
	}

	return &artworkResp, nil
}

package test

import (
	"daily-driver/internal/api/art"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChicagoAPIClient_ErrorHandling(t *testing.T) {
	// Test case 1: Simulate a server error
	serverErrorHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(serverErrorHandler))
	defer server.Close()

	client := art.ChicagoArtInstituteAPI{
		BaseURL: server.URL,
		Client:  &http.Client{},
	}

	_, err := client.GetRandomArtwork()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	// Test case 2: Simulate a malformed JSON response

	// Test case 4: Verify true image URL is constructed
	validResponseHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{
			"data": [{
				"title": "Test Artwork",
				"artist_display": "Test Artist",
				"thumbnail": {"lqip": "test_thumbnail_url"},
				"image_id": "test_image_id",
				"iiif_url": "https://example.com/iiif/2/test_image_id"
			}]
		}`
		w.Write([]byte(response))
	}

	server = httptest.NewServer(http.HandlerFunc(validResponseHandler))
	defer server.Close()

	client.BaseURL = server.URL
	result, err := client.GetRandomArtwork()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Data) == 0 {
		t.Fatalf("expected at least one artwork, got none")
	}

	artwork := result.Data[0]
	if artwork.IIIFURL != "https://example.com/iiif/2/test_image_id/full/full/0/default.jpg" {
		t.Fatalf("expected true image URL, got %s", artwork.IIIFURL)
	}
	malformedJSONHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{")) // Incomplete JSON
	}

	server = httptest.NewServer(http.HandlerFunc(malformedJSONHandler))
	defer server.Close()

	client.BaseURL = server.URL

	_, err = client.GetRandomArtwork()
	if err == nil {
		t.Fatalf("expected an error due to malformed JSON, got nil")
	}

	// Test case 3: Invalid API URL
	client.BaseURL = "http://invalid-url"

	_, err = client.GetRandomArtwork()
	if err == nil {
		t.Fatalf("expected an error due to invalid URL, got %v", err)
	}
}

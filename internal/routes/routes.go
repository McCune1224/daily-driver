package routes

// URL constants for the daily-driver application
const (

	//WARNING: Should this be here? No! But I am not sure where else to put it right now and makign a new package for one const seems silly.
	PanelRotationIntervalSeconds int = 5
	// Root
	Root = "/"

	// Admin routes
	GarminBase   = "/garmin"
	GarminUpload = GarminBase + "/upload"

	// Panel routes
	PanelBase  = "/panel"
	PanelIndex = PanelBase + "/index"

	// Art routes
	ArtBase      = "/art"
	ArtRandomAPI = ArtBase + "/random"

	// API routes
	APIBase    = "/api"
	WeatherAPI = APIBase + "/weather"
)

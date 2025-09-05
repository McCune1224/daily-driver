# Smart Apartment Dashboard

A modern, full-screen dashboard application designed for TV display in apartments. Built with Go (Echo framework) for the backend and Alpine.js + HTMX for the frontend.

## Features

- **Real-time Data**: WebSocket connection for live updates
- **Rotating Panels**: Multiple data panels that auto-rotate
- **TV Optimized**: Large fonts, high contrast, full-screen design
- **Modern Stack**: Go + Echo + Alpine.js + HTMX + DaisyUI
- **No Node.js**: Uses CDN for all frontend dependencies

## Data Panels

1. **Weather**: Current conditions and 3-day forecast
2. **Garmin Stats**: Running data, weekly goals, progress
3. **System Info**: CPU, memory, network status
4. **Date & Time**: Current time and date display
5. **News**: Latest headlines and updates
6. **Custom**: Configurable panel for your data

## Quick Start

### Prerequisites

- Go 1.21 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd apartment-dashboard
```

2. Install Go dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run .
```

4. Open your browser to `http://localhost:8080`

### Project Structure

```
apartment-dashboard/
├── main.go           # Application entry point
├── websocket.go      # WebSocket handling
├── handlers.go       # API route handlers
├── go.mod           # Go module definition
├── static/
│   └── index.html   # Frontend application
└── README.md        # This file
```

## Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)

### API Integrations

To integrate real data sources, replace the mock implementations in `handlers.go`:

#### Weather API
Replace the mock weather data in `GetWeather()` with calls to:
- OpenWeatherMap API
- Weather.com API  
- Visual Crossing API

#### Garmin API
Replace the mock Garmin data in `GetGarminData()` with:
- Garmin Connect Developer API
- Terra API for health data
- Direct device integration

#### System Monitoring
Replace the mock system data in `GetSystemInfo()` with:
- psutil for system metrics
- Docker stats API
- Custom monitoring solution

## WebSocket Protocol

The application uses WebSocket for real-time updates. Message format:

```json
{
  "type": "data_type",
  "data": {
    // Panel-specific data
  }
}
```

Supported message types:
- `weather_update`
- `system_update` 
- `garmin_update`
- `welcome`

## Frontend Technologies

- **Alpine.js**: Reactive state management
- **HTMX**: WebSocket integration and AJAX
- **DaisyUI**: UI components (via CDN)
- **Tailwind CSS**: Styling utilities

## Deployment

### Docker

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./main"]
```

### Systemd Service

Create `/etc/systemd/system/apartment-dashboard.service`:

```ini
[Unit]
Description=Smart Apartment Dashboard
After=network.target

[Service]
Type=simple
User=dashboard
WorkingDirectory=/opt/apartment-dashboard
ExecStart=/opt/apartment-dashboard/main
Restart=always

[Install]
WantedBy=multi-user.target
```

## TV Setup

For optimal TV display:

1. **Browser Settings**:
   - Enable full-screen mode (F11)
   - Disable screensaver/sleep mode
   - Set homepage to dashboard URL

2. **TV Settings**:
   - Adjust overscan/picture size
   - Set appropriate refresh rate
   - Configure auto-power on

3. **Network**:
   - Use wired connection if possible
   - Configure static IP for reliability

## Development

### Adding New Panels

1. Add panel data structure to `dashboard()` function in `index.html`
2. Create panel template in the rotating display section
3. Add WebSocket message handler if needed
4. Update backend API if external data required

### API Development

1. Add new handler function in `handlers.go`
2. Register route in `main.go`
3. Update frontend to consume new endpoint

### WebSocket Events

Add new message types in `websocket.go` and corresponding handlers in the frontend.

## Troubleshooting

### WebSocket Connection Issues
- Check firewall settings
- Verify port 8080 is accessible
- Check browser console for errors

### Static File Issues
- Ensure `static/` directory exists
- Check file permissions
- Verify CDN links are accessible

### Performance Issues
- Monitor system resources
- Check WebSocket message frequency
- Optimize panel rotation timing

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues and questions:
- Create a GitHub issue
- Check existing documentation
- Review the troubleshooting guide

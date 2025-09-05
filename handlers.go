package main

import (
    "net/http"
    "time"
    "math/rand"

    "github.com/labstack/echo/v4"
)

type WeatherData struct {
    Temperature float64 `json:"temperature"`
    Condition   string  `json:"condition"`
    Humidity    int     `json:"humidity"`
    WindSpeed   int     `json:"wind_speed"`
    Forecast    []ForecastDay `json:"forecast"`
}

type ForecastDay struct {
    Day       string `json:"day"`
    High      int    `json:"high"`
    Low       int    `json:"low"`
    Condition string `json:"condition"`
}

type GarminData struct {
    LastRun      LastRun      `json:"last_run"`
    WeeklyStats  WeeklyStats  `json:"weekly_stats"`
    Goals        Goals        `json:"goals"`
}

type LastRun struct {
    Distance float64 `json:"distance"`
    Duration string  `json:"duration"`
    Pace     string  `json:"pace"`
    Date     string  `json:"date"`
}

type WeeklyStats struct {
    TotalDistance float64 `json:"total_distance"`
    TotalRuns     int     `json:"total_runs"`
    AvgHeartRate  int     `json:"avg_heart_rate"`
}

type Goals struct {
    WeeklyDistance float64 `json:"weekly_distance"`
    WeeklyRuns     int     `json:"weekly_runs"`
}

type SystemInfo struct {
    CPU         int    `json:"cpu"`
    Memory      int    `json:"memory"`
    Network     string `json:"network"`
    Uptime      string `json:"uptime"`
    Temperature int    `json:"temperature"`
}

func GetWeather(c echo.Context) error {
    // Mock weather data - replace with actual API calls
    weather := WeatherData{
        Temperature: 72 + float64(rand.Intn(10)-5),
        Condition:   "Partly Cloudy",
        Humidity:    60 + rand.Intn(20),
        WindSpeed:   5 + rand.Intn(10),
        Forecast: []ForecastDay{
            {Day: "Today", High: 78, Low: 65, Condition: "Sunny"},
            {Day: "Tomorrow", High: 75, Low: 62, Condition: "Cloudy"},
            {Day: "Friday", High: 80, Low: 68, Condition: "Rain"},
        },
    }

    return c.JSON(http.StatusOK, weather)
}

func GetGarminData(c echo.Context) error {
    // Mock Garmin data - replace with actual Garmin API calls
    garmin := GarminData{
        LastRun: LastRun{
            Distance: 5.2,
            Duration: "28:45",
            Pace:     "5:32",
            Date:     time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
        },
        WeeklyStats: WeeklyStats{
            TotalDistance: 18.6,
            TotalRuns:     4,
            AvgHeartRate:  158,
        },
        Goals: Goals{
            WeeklyDistance: 25,
            WeeklyRuns:     5,
        },
    }

    return c.JSON(http.StatusOK, garmin)
}

func GetSystemInfo(c echo.Context) error {
    // Mock system info - replace with actual system monitoring
    system := SystemInfo{
        CPU:         40 + rand.Intn(30),
        Memory:      60 + rand.Intn(25),
        Network:     "Connected",
        Uptime:      "5 days, 12 hours",
        Temperature: 40 + rand.Intn(15),
    }

    return c.JSON(http.StatusOK, system)
}

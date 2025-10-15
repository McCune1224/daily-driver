package service

import (
	"bytes"
	"daily-driver/internal/db"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

type GarminActivity struct {
	Activity *filedef.Activity
	Summary  *filedef.ActivitySummary
}

func DecodeGarminActivity(file *db.GarminFitFile) (*GarminActivity, error) {
	reader := bytes.NewReader(file.Data)
	garminDecoder := decoder.New(reader)
	fit, err := garminDecoder.Decode()
	if err != nil {
		return nil, err
	}

	activity := filedef.NewActivity(fit.Messages...)
	summary := filedef.NewActivitySummary(fit.Messages...)
	return &GarminActivity{activity, summary}, nil
}

// GetTotalMiles returns the total distance of the activity in miles.
func (g *GarminActivity) GetTotalMiles() float64 {
	if g == nil || g.Summary == nil || len(g.Summary.Sessions) == 0 {
		return 0
	}
	session := g.Summary.Sessions[0]
	if session == nil || session.TotalDistance == 0 {
		return 0
	}
	// TotalDistance is in meters, convert to miles
	return float64(session.TotalDistance) / 1609.34
}

func (g *GarminActivity) GetTotalKilometers() float64 {
	if g == nil || g.Summary == nil || len(g.Summary.Sessions) == 0 {
		return 0
	}
	session := g.Summary.Sessions[0]
	if session == nil || session.TotalDistance == 0 {
		return 0
	}
	// TotalDistance is in meters, convert to kilometers
	return float64(session.TotalDistance) / 1000
}

// GetTotalTime returns the total elapsed time of the activity in seconds.
func (g *GarminActivity) GetTotalTime() float64 {
	if g == nil || g.Summary == nil || len(g.Summary.Sessions) == 0 {
		return 0
	}
	session := g.Summary.Sessions[0]
	if session == nil || session.TotalElapsedTime == 0 {
		return 0
	}
	return float64(session.TotalElapsedTime)
}

// GetAveragePace returns the average pace in minutes per mile.
func (g *GarminActivity) GetAveragePace() float64 {
	totalMiles := g.GetTotalMiles()
	totalTime := g.GetTotalTime() // in seconds
	if totalMiles == 0 {
		return 0
	}
	paceSecondsPerMile := totalTime / totalMiles
	return paceSecondsPerMile / 60 // minutes per mile
}

// GetFastestPace returns the fastest (minimum) pace in minutes per mile from all laps in the activity.
func (g *GarminActivity) GetFastestPace() float64 {
	if g == nil || g.Activity == nil || len(g.Activity.Laps) == 0 {
		return 0
	}
	minPace := 1e9 // Large initial value
	for _, lap := range g.Activity.Laps {
		if lap == nil || lap.TotalElapsedTime == 0 || lap.TotalDistance == 0 {
			continue
		}
		distanceMiles := float64(lap.TotalDistance) / 1609.34
		if distanceMiles <= 0 {
			continue
		}
		paceSecondsPerMile := float64(lap.TotalElapsedTime) / distanceMiles
		paceMinutesPerMile := paceSecondsPerMile / 60
		if paceMinutesPerMile < minPace {
			minPace = paceMinutesPerMile
		}
	}
	if minPace == 1e9 {
		return 0
	}
	return minPace
}

// GetAverageHR returns the average heart rate in bpm for the activity.
func (g *GarminActivity) GetAverageHR() float64 {
	// Prefer summary session's AvgHeartRate if available
	if g == nil || g.Summary == nil || len(g.Summary.Sessions) == 0 {
		return 0
	}
	session := g.Summary.Sessions[0]
	if session != nil && session.AvgHeartRate != 0 {
		return float64(session.AvgHeartRate)
	}
	// Fallback to calculating average from records if available
	if g.Activity == nil || len(g.Activity.Records) == 0 {
		return 0
	}
	total := 0.0
	count := 0
	for _, rec := range g.Activity.Records {
		if rec == nil || rec.HeartRate == 0 {
			continue
		}
		total += float64(rec.HeartRate)
		count++
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// GetHighestHR returns the highest heart rate in bpm for the activity.
func (g *GarminActivity) GetHighestHR() float64 {
	// Prefer summary session's MaxHeartRate if available
	if g == nil || g.Summary == nil || len(g.Summary.Sessions) == 0 {
		return 0
	}
	session := g.Summary.Sessions[0]
	if session != nil && session.MaxHeartRate != 0 {
		return float64(session.MaxHeartRate)
	}
	// Fallback to calculating max from records if available
	if g.Activity == nil || len(g.Activity.Records) == 0 {
		return 0
	}
	maxHR := 0.0
	for _, rec := range g.Activity.Records {
		if rec == nil || rec.HeartRate == 0 {
			continue
		}
		hr := float64(rec.HeartRate)
		if hr > maxHR {
			maxHR = hr
		}
	}
	return maxHR
}

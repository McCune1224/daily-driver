package service

import (
	"bytes"
	"daily-driver/internal/db"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

func DecodeGarminActivity(file *db.GarminFitFile) (*filedef.Activity, *filedef.ActivitySummary, error) {
	reader := bytes.NewReader(file.Data)
	garminDecoder := decoder.New(reader)
	fit, err := garminDecoder.Decode()
	if err != nil {
		return nil, nil, err
	}

	activity := filedef.NewActivity(fit.Messages...)
	summary := filedef.NewActivitySummary(fit.Messages...)
	return activity, summary, nil

}

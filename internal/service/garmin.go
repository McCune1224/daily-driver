package service

import (
	"os"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/proto"
)

func DecodeGarminActivity(filepath string) (*proto.FIT, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	garminDecoder := decoder.New(f)
	fit, err := garminDecoder.Decode()
	if err != nil {
		return nil, err
	}
	return fit, nil
}

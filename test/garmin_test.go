package test

import (
	"os"
	"testing"

	"github.com/muktihari/fit/decoder"
	"github.com/muktihari/fit/profile/filedef"
)

const GarminFolderPath = "./"

func TestFoobar(t *testing.T) {
	t.Logf("Reading files from %s", GarminFolderPath)
	//glob read current directory files with .fit:
	files, err := os.ReadDir(GarminFolderPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".fit" {
			t.Logf("Found file: %s", file.Name())
		}
	}

	f, err := os.Open(files[0].Name())
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := decoder.New(f)

	decFit, err := dec.Decode()
	if err != nil {
		t.Fatal(err)
	}

	activity := filedef.NewActivity(decFit.Messages...)
	if activity == nil {
		t.Fatal("Failed to create activity from FIT messages")
	}
	// t.Logf("Activity: %+v", activity)

	// service.DecodeGarminActivity(GarminFolderPath + "2023-06-01-07-30-00.fit")
}

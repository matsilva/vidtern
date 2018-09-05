package videoassets_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/matsilva/vidtern/videoassets"
)

func TestDownloadFile(t *testing.T) {
	tmpdir := os.TempDir()
	cases := []struct {
		name             string
		filename         string
		url              string
		expectedFileName string
	}{
		{
			"download image",
			"peacock-feathers-3617474_1280.jpg",
			"https://cdn.pixabay.com/photo/2018/08/19/19/56/peacock-feathers-3617474_1280.jpg",
			path.Join(tmpdir, "peacock-feathers-3617474_1280.jpg"),
		},
		{
			"download video",
			"video-15138_tiny.mp4",
			"https://pixabay.com/en/videos/download/video-15138_tiny.mp4",
			path.Join(tmpdir, "video-15138_tiny.mp4"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := videoassets.DownloadFile(tc.url, tmpdir)
			if err != nil {
				t.Fatalf("could now download file: %v", err)
			}
			_, err = os.Stat(tc.expectedFileName)
			if err != nil {
				t.Fatalf("expected file does not exist %s", tc.expectedFileName)
			}
		})
		err := os.Remove(tc.expectedFileName)
		if err != nil {
			log.Fatalf("could not remove file %s", tc.expectedFileName)
		}
	}

	//TODO: Add test to make sure bad url will fail
}

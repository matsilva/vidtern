package videoconfig_test

import (
	"testing"

	"github.com/matsilva/vidtern/videoconfig"
)

func TestFromFile(t *testing.T) {
	t.Run("config scenes", func(t *testing.T) {
		v, err := videoconfig.FromFile("../testdata/sample-video-images.json")
		if err != nil {
			t.Fatalf("could not get video config from json %v", err)
		}
		if len(v.Scenes) != 3 {
			t.Errorf("expected 3 scenes, got %d", len(v.Scenes))
		}
		firstScene := v.Scenes[0]
		if firstScene.Duration != 10000 {
			t.Errorf("expected first scene duration to be 10000, got %d", firstScene.Duration)
		}
		if firstScene.Media != "https://cdn.pixabay.com/photo/2015/06/11/21/55/lake-erie-806259_1280.jpg" {
			t.Errorf("expected first scene media to be https://media.licdn.com/media/gcrc/dms/image/C4D12AQHkqA7nk15a2g/article-cover_image-shrink_423_752/0?e=1541030400&v=beta&t=EZ9Zp8Licz9X1YVB9uO7hFWkCwL5nyWCHYDrzoPcQvc, got %s", firstScene.Media)
		}
		if firstScene.Text != "2 Ways to Get Involved with the Erie Entrepreneurial Community" {
			t.Errorf("expected first scene text to be \"2 Ways to Get Involved with the Erie Entrepreneurial Community\", got \"%s\"", firstScene.Text)
		}
	})

	t.Run("config other properties", func(t *testing.T) {
		v, err := videoconfig.FromFile("../testdata/sample-video-images.json")
		if err != nil {
			t.Fatalf("could not get video config from json %v", err)
		}
		v.JobDir = "hello_world"
		if v.JobDir != "hello_world" {
			t.Fatalf("expected JobDir to be \"hello_world\"; got: %v", v.JobDir)
		}

		if v.VideoName != "vidtern_finished_video" {
			t.Fatalf("expected VideoName to be \"vidtern_finished_video\"; got: %v", v.VideoName)
		}
	})

	t.Run("bad config", func(t *testing.T) {
		v, err := videoconfig.FromFile("./does-not-exist.json")
		if err == nil {
			t.Errorf("expected an error from videoconfig.FromFile, got %v", v)
		}
	})
}

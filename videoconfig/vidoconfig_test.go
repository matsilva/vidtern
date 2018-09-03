package videoconfig_test

import (
	"testing"

	"github.com/matsilva/vidtern/videoconfig"
)

func TestFromFile(t *testing.T) {
	videoConfig, err := videoconfig.FromFile("./test-video.json")
	if err != nil {
		t.Fatalf("could not get video config from json %v", err)
	}
	if len(videoConfig.Scenes) != 3 {
		t.Errorf("expected 3 scenes, got %d", len(videoConfig.Scenes))
	}
	firstScene := videoConfig.Scenes[0]
	if firstScene.Duration != 10000 {
		t.Errorf("expected first scene duration to be 10000, got %d", firstScene.Duration)
	}
	if firstScene.Media != "https://media.licdn.com/media/gcrc/dms/image/C4D12AQHkqA7nk15a2g/article-cover_image-shrink_423_752/0?e=1541030400&v=beta&t=EZ9Zp8Licz9X1YVB9uO7hFWkCwL5nyWCHYDrzoPcQvc" {
		t.Errorf("expected first scene media to be https://media.licdn.com/media/gcrc/dms/image/C4D12AQHkqA7nk15a2g/article-cover_image-shrink_423_752/0?e=1541030400&v=beta&t=EZ9Zp8Licz9X1YVB9uO7hFWkCwL5nyWCHYDrzoPcQvc, got %s", firstScene.Media)
	}
	if firstScene.Text != "2 Ways to Get Involved with the Erie Entrepreneurial Community" {
		t.Errorf("expected first scene text to be \"2 Ways to Get Involved with the Erie Entrepreneurial Community\", got \"%s\"", firstScene.Text)
	}

	videoConfig, err = videoconfig.FromFile("./does-not-exist.json")
	if err == nil {
		t.Errorf("expected an error from videoconfig.FromFile, got %v", videoConfig)
	}
}

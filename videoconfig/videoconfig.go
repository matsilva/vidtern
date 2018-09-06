package videoconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//Scene contains information about an individual
//scene in the video
type Scene struct {
	Text     string `json:"text"`
	Media    string `json:"media"`
	Duration int    `json:"duration"`

	MediaInfo struct {
		FilePath string
		Type     string
	}
}

//VideoConfig contains the configuration for the video to be created
type VideoConfig struct {
	Scenes    []Scene
	VideoName string `json:"name"`
	JobDir    string
}

//FromFile returns a configuration parsed from the given file.
func FromFile(file string) (*VideoConfig, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s; err: %v", file, err)
	}

	return FromJSON(b)
}

//FromJSON returns a configuration parsed from a given json buffer
//good for parsing from a http request
func FromJSON(data []byte) (*VideoConfig, error) {
	var videoConfig VideoConfig
	if err := json.Unmarshal(data, &videoConfig); err != nil {
		return nil, fmt.Errorf("could not unmarshal data %v; err: %v", data, err)
	}
	if videoConfig.VideoName == "" {
		videoConfig.VideoName = "vidtern_finished_video"
	}
	return &videoConfig, nil
}

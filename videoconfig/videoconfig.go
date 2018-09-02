package videoconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//VideoConfig contains the configuration for the video to be created
type VideoConfig struct {
	Scenes []struct {
		Text     string `json:"text"`
		Media    string `json:"media"`
		Duration int    `json:"duration"`
	}
}

//FromFile returns a configuration parsed from the given file.
func FromFile(file string) (*VideoConfig, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s", file)
	}

	return FromJSON(b)
}

//FromJSON returns a configuration parsed from a given json buffer
//good for parsing from a http request
func FromJSON(data []byte) (*VideoConfig, error) {
	var videoConfig VideoConfig
	if err := json.Unmarshal(data, &videoConfig); err != nil {
		return nil, fmt.Errorf("could not unmarshal data %v", data)
	}
	return &videoConfig, nil
}

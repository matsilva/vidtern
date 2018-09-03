package vidtern

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/matsilva/vidtern/videoassets"
	"github.com/matsilva/vidtern/videoconfig"
	filetype "gopkg.in/h2non/filetype.v1"
)

//Create makes a video given a video config
func Create(videoConfig *videoconfig.VideoConfig) error {

	err := getVideoAssets(videoConfig)
	if err != nil {
		return err
	}
	//take the videoConfig and call ffmpeg to make the video!
	//create videos per scene in go routines

	//stich them all together
	return nil
}

//Gets all of the assets needed to make the video
func getVideoAssets(videoConfig *videoconfig.VideoConfig) error {

	var MediaTypes [2]string
	MediaTypes[0] = "image"
	MediaTypes[1] = "video"

	for _, scene := range videoConfig.Scenes {

		err := videoassets.DownloadFile(scene.Media, videoConfig.JobDir)
		if err != nil {
			return err
		}

		//add filepath to video config
		filepath := path.Join(videoConfig.JobDir, path.Base(scene.Media))
		scene.MediaInfo.FilePath = filepath

		//get the filetype and add to video config
		buf, err := ioutil.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("could not read file %s to determine type; err %v", filepath, err)
		}

		if filetype.IsImage(buf) {
			scene.MediaInfo.Type = MediaTypes[0]
			//TODO: Get dimensions
		}

		if filetype.IsVideo(buf) {
			scene.MediaInfo.Type = MediaTypes[1]
			//TODO: Get dimensions, fps, duration
		}

		if scene.MediaInfo.Type == "" {
			kind, _ := filetype.Match(buf)
			return fmt.Errorf("supported file types are image and video, %s is %s", path.Base(filepath), kind)
		}
	}
	return nil
}

//Creates an individual video from scene
func createScene(videoConfig *videoconfig.VideoConfig, sceneIndex int) error {

	//example text param
	//drawtext="fontfile=/usr/share/fonts/truetype/freefont/FreeSerif.ttf: text='Test Text':\
	//x=100: y=50: fontsize=24: fontcolor=yellow@0.2: box=1: boxcolor=red@0.2"
	return nil
}

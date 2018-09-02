package vidtern

import (
	"github.com/matsilva/vidtern/videoconfig"
)

//Create makes a video given a video config
func Create(videoConfig *videoconfig.VideoConfig) error {

	//take the videoConfig and call ffmpeg to make the video!

	//create videos per scene in go routines
	//stich them all together
	return nil
}

//Gets all of the assets needed to make the video
func defineVideoAssets(videoConfig *videoconfig.VideoConfig) error {
	//download assets
	//get more information about the assets
	// media type: [image, video]

	//if image
	// dimensions

	//if video
	// dimensions
	//fps
	//duration
	return nil
}

func downloadVideoAssets(videoConfig *videoconfig.VideoConfig) error {
	waitCount := len(videoConfig.Scenes)
	for scene := range videoConfig.Scenes {
		filename := urlParts[len(urlParts)-1]
	}
}

//Creates an individual video from scene
func createScene(videoConfig *videoconfig.VideoConfig, sceneIndex int) error {

	//example text param
	//drawtext="fontfile=/usr/share/fonts/truetype/freefont/FreeSerif.ttf: text='Test Text':\
	//x=100: y=50: fontsize=24: fontcolor=yellow@0.2: box=1: boxcolor=red@0.2"
	return nil
}

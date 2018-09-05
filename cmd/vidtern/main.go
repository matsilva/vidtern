package main

import (
	"flag"
	"log"
	"os"

	"github.com/matsilva/vidtern/videoconfig"
	"github.com/matsilva/vidtern/vidtern"
)

func main() {
	videoName := flag.String("videoname", "vidtern_video", "name for the finished video")
	configPath := flag.String("config", "", "path to the video config json file")
	jobDir := flag.String("jobdir", os.TempDir(), "directory to use for the video job")
	flag.Parse()

	videoConfig, err := videoconfig.FromFile(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	videoConfig.VideoName = *videoName + ".mp4"
	videoConfig.JobDir = *jobDir

	err = vidtern.Create(videoConfig)

	if err != nil {
		log.Fatal(err)
	}
}

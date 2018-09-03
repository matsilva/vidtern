package main

import (
	"flag"
	"log"
	"os"

	"github.com/matsilva/vidtern/videoconfig"
	"github.com/matsilva/vidtern/vidtern"
)

func main() {
	configPath := flag.String("config", "", "path to the video config json file")
	jobDir := flag.String("jobdir", os.TempDir(), "directory to use for the video job")
	flag.Parse()

	videoConfig, err := videoconfig.FromFile(*configPath)
	videoConfig.JobDir = *jobDir

	if err != nil {
		log.Fatal(err)
	}

	err = vidtern.Create(videoConfig)

	if err != nil {
		log.Fatal(err)
	}
}

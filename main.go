package main

import (
	"log"

	"github.com/matsilva/vidtern/videoconfig"
	"github.com/matsilva/vidtern/vidtern"
)

func main() {

	videoConfig, err := videoconfig.FromFile("./test-video.json")
	if err != nil {
		log.Fatal(err)
	}

	err = vidtern.Create(videoConfig)

	if err != nil {
		log.Fatal(err)
	}

}

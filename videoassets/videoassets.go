package videoassets

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/matsilva/vidtern/videoconfig"
	filetype "gopkg.in/h2non/filetype.v1"
)

func showProgress(filepath string, total int, stop chan bool) {
	start := time.Now()
	for {
		select {
		case <-stop:
			elapsed := time.Since(start)
			log.Printf("%s => completed in %s\n", filepath, elapsed)
			return
		default:

			f, err := os.Open(filepath)
			if err != nil {
				log.Fatalf("could not open %s to show progress; err %v", filepath, err)
			}

			fstat, err := f.Stat()
			if err != nil {
				log.Fatalf("could not get stats for %s; err %v", filepath, err)
			}

			size := fstat.Size()
			//prevent failed percent equation
			if size == 0 {
				size = 1
			}

			percent := float64(size) / float64(total) * 100

			log.Printf("%.0f%%\n", percent)
			//print every second
			time.Sleep(time.Second)

		}
	}
}

//DownloadFile downloads a file from a url
func DownloadFile(url, dest string) error {
	filename := path.Base(url)
	log.Printf("downloading file: %s from %s\n", filename, url)
	filepath := path.Join(dest, filename)

	// create writable file
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("could not create file to write %s; err %v", filepath, err)
	}

	defer f.Close()

	//request the file
	res, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("could not GET %s; err %v", url, err)
	}

	defer res.Body.Close()

	//get the total bytes of the file
	total, err := strconv.Atoi(res.Header.Get("Content-Length"))

	if err != nil {
		return fmt.Errorf("could not convert Content-Length to int %v", err)
	}

	//create a signal to tell showProgress to stop
	stop := make(chan bool)
	go showProgress(filepath, total, stop)

	//write file
	_, err = io.Copy(f, res.Body)

	if err != nil {
		return fmt.Errorf("could not write to file %s; err %v", filepath, err)
	}
	stop <- true
	return nil
}

//GetVideoAssets will download all of the media assets needed for the video
//and further add media information to the video config
func GetVideoAssets(videoConfig *videoconfig.VideoConfig) error {

	var MediaTypes [2]string
	MediaTypes[0] = "image"
	MediaTypes[1] = "video"

	for _, scene := range videoConfig.Scenes {

		err := DownloadFile(scene.Media, videoConfig.JobDir)
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

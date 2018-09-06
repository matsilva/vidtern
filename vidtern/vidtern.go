package vidtern

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/matsilva/vidtern/videoassets"
	"github.com/matsilva/vidtern/videoconfig"
)

//Create makes a video given a video config
func Create(videoConfig *videoconfig.VideoConfig) error {

	//get assets syncronously &
	//try to avoid overloading network requests
	err := videoassets.GetVideoAssets(videoConfig)
	if err != nil {
		return err
	}

	//TODO: Add concurrency
	//allow up to 3 jobs at a time, configurable by the user
	for index := range videoConfig.Scenes {
		err := createScene(videoConfig, index)
		if err != nil {
			return err
		}
	}

	//concat all videos together into a single video
	err = createVideoFromScenes(videoConfig)
	if err != nil {
		return err
	}
	return nil
}

//Creates an individual video from scene
func createScene(videoConfig *videoconfig.VideoConfig, index int) error {

	scene := videoConfig.Scenes[index]
	//https://golang.org/pkg/os/exec/#Command
	var cmdOpts []string

	fmt.Printf("media type: %s for %s", scene.MediaInfo.Type, scene.MediaInfo.FilePath)
	switch scene.MediaInfo.Type {
	case "image":
		//add image and loop it for the configured duration
		duration := strconv.Itoa(int(scene.Duration / 1000)) //fix this
		cmdOpts = append(cmdOpts, "-loop 1 -i "+scene.MediaInfo.FilePath+" -c:v libx264 -t "+duration)
	case "video":
		cmdOpts = append(cmdOpts, "-i "+scene.MediaInfo.FilePath)
	}

	if scene.Text != "" {
		//https://www.ffmpeg.org/ffmpeg-filters.html#drawtext-1
		drawtext := "text='" + scene.Text + "': x=100: y=50: fontsize=24: fontcolor=yellow@0.2: box=1: boxcolor=red@0.2"
		cmdOpts = append(cmdOpts, "-vf drawtext=\""+drawtext+"\"")
	}

	filename := "vidtern__scene_" + strconv.Itoa(index+1) + ".mp4"

	//add out filename
	cmdOpts = append(cmdOpts, path.Join(videoConfig.JobDir, filename))
	fmt.Printf("running: ffmpeg %s", strings.Join(cmdOpts, " "))
	cmd := exec.Command("ffmpeg", cmdOpts...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	log.Printf("ffmpeg stdout: %s", stdout.String())
	log.Printf("ffmpeg stderr: %s", stderr.String())
	if err != nil {
		return fmt.Errorf("could not create %s; err: %v", filename, err)
	}
	return nil
}

func createVideoFromScenes(videoConfig *videoconfig.VideoConfig) error {

	tmpfile, err := makeSceneFileList(videoConfig)
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	cmdOpts := []string{
		"-f concat -i",
		path.Join(videoConfig.JobDir, tmpfile.Name()),
		"-c copy",
		path.Join(videoConfig.JobDir, videoConfig.VideoName+".mp4"),
	}

	cmd := exec.Command("ffmpeg", cmdOpts...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not create %s; err: %v", videoConfig.VideoName, err)
	}
	log.Printf("ffmpeg: %s", out.String())
	return nil
}

func makeSceneFileList(videoConfig *videoconfig.VideoConfig) (*os.File, error) {

	var content strings.Builder
	for index := range videoConfig.Scenes {
		filename := "vidtern__scene_" + strconv.Itoa(index+1) + ".mp4"
		content.WriteString("file '" + path.Join(videoConfig.JobDir, filename) + "'\n")

	}
	tmpfile, err := ioutil.TempFile(videoConfig.JobDir, "vidtern_scenes.txt")
	if err != nil {
		return nil, fmt.Errorf("could not create new temp file %v", err)
	}

	//Make sure to to remove the file
	//defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(content.String())); err != nil {
		return nil, fmt.Errorf("could not write to temp file %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		return nil, fmt.Errorf("could not close temp file %v", err)
	}
	return tmpfile, nil
}

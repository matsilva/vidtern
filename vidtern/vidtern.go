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
	var args []string
	args = append(args, "-y") //overwrite

	switch scene.MediaInfo.Type {
	case "image":
		//add image and loop it for the configured duration
		duration := strconv.Itoa(int(scene.Duration / 1000))
		args = append(args, "-loop")
		args = append(args, "1")
		args = append(args, "-i")
		args = append(args, ""+scene.MediaInfo.FilePath+"")
		args = append(args, "-c:v")
		args = append(args, "libx264")
		args = append(args, "-t")
		args = append(args, duration)
	case "video":
		args = append(args, "-i")
		args = append(args, scene.MediaInfo.FilePath)
	}

	// if scene.Text != "" {
	// 	//https://www.ffmpeg.org/ffmpeg-filters.html#drawtext-1
	// 	//TODO: Fix paths to work on any os, make font a predictable path
	// 	var drawtext strings.Builder
	// 	drawtext.WriteString("fontfile=/Users/evsilva22/go/src/github.com/matsilva/vidtern/testdata/Lato-Black.ttf")
	// 	drawtext.WriteString(": text='" + scene.Text + "'")
	// 	drawtext.WriteString(": x=(w-tw)/2")
	// 	drawtext.WriteString(": y=(h/PHI)+th")
	// 	drawtext.WriteString(": fontsize=48")
	// 	drawtext.WriteString(": fontcolor=yellow@0.2")
	// 	drawtext.WriteString(": box=1")
	// 	drawtext.WriteString(": boxcolor=red@0.2")
	// 	args = append(args, "-vf")
	// 	args = append(args, "drawtext=\""+drawtext.String()+"\"")
	// }

	filename := "vidtern__scene_" + strconv.Itoa(index+1) + ".mp4"

	args = append(args, "-vf")
	args = append(args, "scale=w=1920:-2:force_original_aspect_ratio=decrease")
	//add out filename
	args = append(args, path.Join(videoConfig.JobDir, filename))
	fmt.Printf("running: ffmpeg %s\n", strings.Join(args, " "))
	// cmd := exec.Command("ffmpeg", cmdArgs.String())
	cmd := exec.Command("ffmpeg", args...)
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

	// defer os.Remove(tmpfile.Name()) // clean up

	args := []string{
		"-y",
		"-f",
		"concat",
		"-safe", // -safe 0 helps with unsafe filename error https://stackoverflow.com/questions/38996925/ffmpeg-concat-unsafe-file-name
		"0",
		"-i",
		tmpfile.Name(),
		"-c",
		"copy",
		path.Join(videoConfig.JobDir, videoConfig.VideoName+".mp4"),
	}

	fmt.Printf("running: ffmpeg %s\n", strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
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
	content.WriteString("ffconcat version 1.0\n")
	for index, scene := range videoConfig.Scenes {
		filename := "vidtern__scene_" + strconv.Itoa(index+1) + ".mp4"
		content.WriteString("file '" + path.Join(videoConfig.JobDir, filename) + "'\n")
		if scene.Duration > 0 {
			duration := strconv.Itoa(int(scene.Duration / 1000))
			content.WriteString("duration " + duration + "\n")
		}
	}
	tmpfile, err := ioutil.TempFile(videoConfig.JobDir, "vidtern_scenes")
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

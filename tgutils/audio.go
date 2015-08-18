// Package tgutils provides extra functions to make certain tasks easier.
package tgutils

import (
	"github.com/syfaro/telegram-bot-api"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextSuffix() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

// this function ripped from ioutils.TempFile, except with a suffix, instead of prefix.
func tempFileWithSuffix(dir, suffix string) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, nextSuffix()+suffix)
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

// EncodeAudio takes a file and attempts to convert it to a .ogg for Telegram.
// It then updates the path to the audio file in the AudioConfig.
//
// This function requires ffmpeg and opusenc to be installed on the system!
func EncodeAudio(audio *tgbotapi.AudioConfig) error {
	f, err := tempFileWithSuffix(os.TempDir(), "_tgutils.ogg")
	if err != nil {
		return err
	}
	defer f.Close()

	ffmpegArgs := []string{
		"-i",
		audio.FilePath,
		"-f",
		"wav",
		"-",
	}

	opusArgs := []string{
		"--bitrate",
		"256",
		"-",
		f.Name(),
	}

	c1 := exec.Command("ffmpeg", ffmpegArgs...)
	c2 := exec.Command("opusenc", opusArgs...)

	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout
	c2.Start()
	c1.Run()
	c2.Wait()

	return nil
}

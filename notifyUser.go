package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/beeep"
)

func notify() {
	log.Println("Time to take a break !")
	workingDir, _ := os.Getwd()
	err := beeep.Notify("Time to take a break !", "It's important to take a break to relieve eye strain", filepath.Join(workingDir, "eye.png"))
	if err != nil {
		log.Fatalln(err)
	}
	go playSound()
}

func playSound() {
	readCloser := ioutil.NopCloser(bytes.NewReader(resourceNotifMp3.StaticContent))
	streamer, format, err := mp3.Decode(readCloser)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
	time.Sleep(1 * time.Second)
}

package util

import (
	"DouYin/config"
	"log"
)

type Ffmsg struct {
	VideoName string
	ImageName string
}

var Ffchan chan Ffmsg

func InitFfmpeg() {
	Ffchan = make(chan Ffmsg, config.MaxMsgCount) // Initialize communication channel for FFmpeg processing
	go dispatcher()                               // Start dispatcher goroutine to process FFmpeg tasks
}

func dispatcher() {
	// Dispatcher for FFmpeg processing tasks
	for ffmsg := range Ffchan {
		go func(f Ffmsg) {
			err := Ffmpeg(f.VideoName, f.ImageName) // Execute FFmpeg processing
			if err != nil {
				Ffchan <- f                   // Retry if processing fails
				log.Fatal("dispatcher error") // Log error if processing fails
			}
		}(ffmsg)
	}
}

func Ffmpeg(videoName string, imageName string) error {
	// Execute FFmpeg command via SSH
	session, err := SshClient.NewSession()
	if err != nil {
		log.Fatal("bulid ssh session error", err) // Log error if session creation fails
	}
	defer session.Close()
	combo, err := session.CombinedOutput("/home/zhouyx/ffmpeg/bin/ffmpeg -ss 00:00:01 -i /home/zhouyx/ftpuser/video/" + videoName + ".mp4 -vframes 1 /home/zhouyx/ftpuser/images/" + imageName + ".jpg")
	if err != nil {
		log.Fatal("command:", string(combo)) // Log error and output if FFmpeg command fails
		return err
	}
	return nil
}

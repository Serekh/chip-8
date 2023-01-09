package utils

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"log"
)

const (
	sampleRate = 44100
)

type Beeper struct {
	player *audio.Player
}

func NewBeeper(beepSound []byte) *Beeper {
	return &Beeper{
		player: audio.NewContext(sampleRate).NewPlayerFromBytes(beepSound),
	}
}

func (beeper *Beeper) Play() {
	if beeper.player.IsPlaying() {
		log.Fatalf("Error trying to play sound the beeper already playing")
	}

	_ = beeper.player.Rewind()
	beeper.player.Play()
}

func (beeper *Beeper) Close() {
	_ = beeper.player.Close()
}

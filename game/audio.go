package game

import "github.com/hajimehoshi/ebiten/v2/audio"

type Audio struct {
	audioCtx   *audio.Context
	jumpPlayer *audio.Player
}

func NewAudio() Audio {
	ctx := audio.NewContext(AudioSampleRate)
	return Audio{
		audioCtx: ctx,
	}
}

package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"time"
)

type Control int

const (
	ScreenWidth           = 300
	ScreenHeight          = 212
	TimeStep              = 1000 / 60
	TimeStepSec           = float64(TimeStep) / float64(1000)
	LeftButton    Control = 0
	RightButton   Control = 1
	FlapButton    Control = 2
	GodModeButton Control = 3
	PauseButton   Control = 4
	SoundButton   Control = 5
	SkidMillis            = 500
)

var (
	MoveSpeed = []float64{
		0,
		0.5,
		1.0,
		2.0,
		2.5,
	}
	WalkAnimSpeed = []time.Duration{
		time.Millisecond * time.Duration(140),
		time.Millisecond * time.Duration(80),
		time.Millisecond * time.Duration(40),
		time.Millisecond * time.Duration(9),
	}
	SpawnPoints = [][]int{
		{230, 96},  // right
		{132, 168}, // bottom
		{109, 53},  // top
		{16, 104},  // left
	}
	Controls = map[Control]ebiten.Key{
		LeftButton:    ebiten.KeyLeft,
		RightButton:   ebiten.KeyRight,
		FlapButton:    ebiten.KeySpace,
		GodModeButton: ebiten.KeyG,
		PauseButton:   ebiten.KeyP,
		SoundButton:   ebiten.KeyS,
	}

	White = color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	Grey = color.RGBA{
		R: 127,
		G: 127,
		B: 127,
		A: 255,
	}
	Yellow = color.RGBA{
		R: 255,
		G: 255,
		B: 86,
		A: 255,
	}

	SpawnColors = []color.RGBA{
		White,
		Grey,
		Yellow,
	}
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

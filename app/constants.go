package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
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
	CrtButton     Control = 6
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
		{236, 96},  // right
		{132, 168}, // bottom
		{116, 53},  // top
		{16, 104},  // left
	}
	Lanes = []int{
		35,
		89,
		159,
	}

	Controls = map[Control]ebiten.Key{
		LeftButton:    ebiten.KeyLeft,
		RightButton:   ebiten.KeyRight,
		FlapButton:    ebiten.KeySpace,
		GodModeButton: ebiten.KeyG,
		PauseButton:   ebiten.KeyP,
		SoundButton:   ebiten.KeyS,
		CrtButton:     ebiten.KeyC,
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

// WrappedDistance Calculates distance between two points on a playfield that wraps around on the x dimension
// Adapted from:
// https://blog.demofox.org/2017/10/01/calculating-the-distance-between-points-in-wrap-around-toroidal-space/
func WrappedDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(Abs(x2 - x1))
	if dx > float64(ScreenWidth)/2 {
		dx = ScreenWidth - dx
	}

	dy := float64(y2 - y1)
	return math.Sqrt(dx*dx + dy*dy)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

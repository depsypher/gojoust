package entity

import (
	"github.com/depsypher/gojoust/app"
	"github.com/depsypher/gojoust/assets/audio"
	"time"
)

type GameObject interface {
	Update(g *GameState)
}

type GameState struct {
	Buzzards  []*Buzzard
	Cliffs    []*Cliff
	Player    *Player
	Keys      map[app.Control]bool
	GodMode   bool
	SoundOn   bool
	Pause     bool
	Debug     string
	Sounds    audio.GameSounds
	WaveStart time.Time
	NextSpawn time.Time
}

func (gs *GameState) CliffAsSprites() []*Sprite {
	r := make([]*Sprite, len(gs.Cliffs))
	for i := range gs.Cliffs {
		r[i] = gs.Cliffs[i].Sprite
	}
	return r
}

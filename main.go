package main

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/depsypher/gojoust/app"
	"github.com/depsypher/gojoust/assets/audio"
	"github.com/depsypher/gojoust/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"log"
	"math/rand"
	"time"
)

var (
	ss *entity.Sheet

	//go:embed app/crt.go
	crt_go []byte
)

func init() {
	sheet, err := entity.LoadSpriteSheet()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load embedded spritesheet: %s", err))
	}
	ss = sheet
}

type toggleAction func()

type Game struct {
	inited bool
	ss     entity.Sheet
	state  *entity.GameState
	screen *ebiten.Image
	crt    *ebiten.Shader
}

func (g *Game) init() {
	defer func() {
		g.inited = true
		g.state = &entity.GameState{
			Keys:      make(map[app.Control]bool),
			GodMode:   false,
			SoundOn:   true,
			CrtOn:     true,
			WaveStart: time.Now(),
		}

		g.state.Player = entity.MakePlayer(ss)
		sp := app.SpawnPoints[1]
		g.state.Player.SetPos(float64(sp[0]), float64(sp[1]))

		g.state.Cliffs = []*entity.Cliff{
			entity.MakeBottomCliff(ss.C1),
			entity.MakeCliff(ss.C2, 105, 136), // mid-bottom
			entity.MakeCliff(ss.C3, 83, 63),   // mid-top
			entity.MakeCliff(ss.C4, -20, 52),  // top-left
			entity.MakeCliff(ss.C5, 253, 52),  // top-right
			entity.MakeCliff(ss.C6, -17, 114), // bottom-left
			entity.MakeCliff(ss.C7, 257, 114), // bottom-right
			entity.MakeCliff(ss.C8, 202, 106), // mid-right
		}

		var err error
		g.state.Sounds, err = audio.LoadSounds()
		if err != nil {
			errSound := errors.New("error loading sounds")
			log.Fatal(errors.Join(errSound, err))
		}

		s, err := ebiten.NewShader(crt_go)
		if err != nil {
			return
		}
		g.crt = s
	}()
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}
	register(app.FlapButton, g.state.Keys)
	register(app.LeftButton, g.state.Keys)
	register(app.RightButton, g.state.Keys)
	toggle(app.GodModeButton, g.state.Keys, func() {
		g.state.Keys[app.GodModeButton] = true
		g.state.GodMode = !g.state.GodMode
	})
	toggle(app.PauseButton, g.state.Keys, func() {
		g.state.Keys[app.PauseButton] = true
		g.state.Pause = !g.state.Pause
	})
	toggle(app.SoundButton, g.state.Keys, func() {
		g.state.Keys[app.SoundButton] = true
		g.state.SoundOn = !g.state.SoundOn
	})
	toggle(app.CrtButton, g.state.Keys, func() {
		g.state.Keys[app.CrtButton] = true
		g.state.CrtOn = !g.state.CrtOn
	})

	if time.Now().After(g.state.WaveStart.Add(time.Duration(3) * time.Second)) {
		if len(g.state.Buzzards) < 3 && time.Now().After(g.state.NextSpawn) {
			buzz := entity.MakeBuzzard(ss)
			point := app.SpawnPoints[rand.Intn(len(app.SpawnPoints))]
			buzz.SetPos(float64(point[0]), float64(point[1]))
			if rand.Float32() < 0.5 {
				buzz.FacingRight = false
			}
			g.state.Buzzards = append(g.state.Buzzards, buzz)
			g.state.NextSpawn = time.Now().Add(time.Duration(1) * time.Second)
		}
	}

	if !g.state.Pause {
		g.state.Player.Update(g.state)
		for _, b := range g.state.Buzzards {
			b.Update(g.state)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if g.screen == nil {
		g.screen = ebiten.NewImage(w, h)
	} else {
		g.screen.Clear()
	}

	if g.state.GodMode {
		ebitenutil.DebugPrint(g.screen, fmt.Sprintf("FPS: %3.2f\nTPS: %3.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
		ebitenutil.DebugPrintAt(g.screen, g.state.Debug, 70, 0)
		ebitenutil.DebugPrintAt(g.screen, fmt.Sprintf("%f", g.state.Player.Y), 70, 20)

		for _, lane := range app.Lanes {
			y := float32(lane)
			vector.StrokeLine(g.screen, 0, y, app.ScreenWidth, y, 1, app.Yellow, false)
		}
	}
	for _, cliff := range g.state.Cliffs {
		cliff.Sprite.DrawSprite(g.screen)
	}
	for _, b := range g.state.Buzzards {
		b.Draw(g.screen)
	}

	g.state.Player.Draw(g.screen)

	if g.state.CrtOn {
		op := &ebiten.DrawRectShaderOptions{}
		op.Images[0] = g.screen
		screen.DrawRectShader(w, h, g.crt, op)
	} else {
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.screen, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return app.ScreenWidth, app.ScreenHeight
}

func register(control app.Control, keys map[app.Control]bool) {
	key := app.Controls[control]
	if ebiten.IsKeyPressed(key) {
		keys[control] = true
	} else {
		delete(keys, control)
	}
}

func toggle(control app.Control, keys map[app.Control]bool, action toggleAction) {
	key := app.Controls[control]
	if ebiten.IsKeyPressed(key) {
		if !keys[control] {
			keys[control] = true
			action()
		}
	} else {
		delete(keys, control)
	}
}

func main() {
	ebiten.SetWindowSize(app.ScreenWidth*3, app.ScreenHeight*3)
	ebiten.SetWindowTitle("GoJoust")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

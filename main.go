package main

import (
	"errors"
	"fmt"
	"github.com/depsypher/gojoust/app"
	"github.com/depsypher/gojoust/assets/audio"
	"github.com/depsypher/gojoust/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

var (
	ss *entity.Sheet
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
}

func (g *Game) init() {
	defer func() {
		g.inited = true
		g.state = &entity.GameState{
			Keys:    make(map[app.Control]bool),
			GodMode: false,
		}

		for i := 1; i < 2; i++ {
			buzz := entity.MakeBuzzard(ss)
			buzz.SetPos(app.ScreenWidth/2, app.ScreenHeight/(float64(i)+1))
			//			g.state.Buzzards = append(g.state.Buzzards, buzz)
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

	if !g.state.Pause {
		g.state.Player.Update(g.state)
		for _, b := range g.state.Buzzards {
			b.Update(g.state)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state.GodMode {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %3.2f\nTPS: %3.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
		ebitenutil.DebugPrintAt(screen, g.state.Debug, 70, 0)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%f", g.state.Player.Y), 70, 20)
	}
	for _, cliff := range g.state.Cliffs {
		cliff.Sprite.DrawSprite(screen)
	}
	for _, b := range g.state.Buzzards {
		b.Draw(screen)
	}
	g.state.Player.Draw(screen)
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

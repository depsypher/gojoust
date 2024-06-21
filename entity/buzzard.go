package entity

import (
	"github.com/depsypher/gojoust/app"
	"github.com/depsypher/gojoust/assets/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
	"time"
)

type Buzzard struct {
	*MountSprite
	bounder     *ebiten.Image
	lastAnimate time.Time
	state       PlayerState
}

func MakeBuzzard(ss *Sheet) *Buzzard {
	return &Buzzard{
		MountSprite: MakeMountSprite(ss.Buzzard),
		bounder:     ss.Bounder,
		lastAnimate: time.Now(),
	}
}

func (b *Buzzard) spawning(gs *GameState) {
	if time.Now().Before(b.lastAnimate.Add(time.Millisecond * time.Duration(30))) {
		return
	}
	if b.spawn <= 20 {
		// emerging
		b.buildSpawn(b, b.spawn)
		b.spawn += 1
		if b.spawn == 20 {
			if err := gs.Sounds[audio.SpawnSound].Play(gs.SoundOn); err != nil {
				log.Fatal("Error playing sound", err)
			}
		}
	} else {
		b.state = MOUNTED
		b.image = b.buildMount()
		b.spawn = 0
		b.Vy = 1
	}
	b.lastAnimate = time.Now()
}

func (b *Buzzard) buildMount() *ebiten.Image {
	frame := b.Images[b.Frame]
	composite := ebiten.NewImage(frame.Bounds().Dx(), frame.Bounds().Dy())

	body := ebiten.NewImageFromImage(frame)
	if b.state == SPAWNING {
		col := app.SpawnColors[rand.Intn(3)]
		body = b.drawSolid(composite.Bounds(), col, body)
	}
	// draw rider
	op := ebiten.DrawImageOptions{}
	y := 0
	op.GeoM.Translate(float64(4), float64(y))
	rider := ebiten.NewImageFromImage(b.bounder)
	if b.state == SPAWNING {
		col := app.SpawnColors[rand.Intn(3)]
		rider = b.drawSolid(rider.Bounds(), col, rider)
	}
	composite.DrawImage(rider, &op)
	op.GeoM.Reset()
	composite.DrawImage(body, &op)
	if !b.facingRight {
		return b.flipX(composite, op)
	}
	return composite
}

func (b *Buzzard) mounted(gs *GameState) {
	if time.Now().After(b.lastAnimate.Add(time.Millisecond * time.Duration(500))) {
		if b.Frame == 5 {
			b.Frame = 6
		} else {
			b.Frame = 5
		}
		b.lastAnimate = time.Now()
	}
	b.image = b.buildMount()
	b.X += 1
	b.Wrap()

	b.Collisions(gs.CliffAsSprites(), func(c *Sprite) {
		//fmt.Printf("colliding %s\n", c.rect())
	})
}

func (b *Buzzard) unmounted(gs *GameState) {
}

func (b *Buzzard) dead(gs *GameState) {
}

func (b *Buzzard) Draw(screen *ebiten.Image) {
	s := b.Sprite
	s.DrawSprite(screen)
}

func (b *Buzzard) Update(gs *GameState) {
	switch b.state {
	case SPAWNING:
		b.spawning(gs)
	case MOUNTED:
		b.mounted(gs)
	case UNMOUNTED:
		b.unmounted(gs)
	case DEAD:
		b.dead(gs)
	}
}

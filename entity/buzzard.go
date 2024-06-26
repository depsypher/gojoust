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
		lastAnimate: time.Time{},
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
		if b.FacingRight {
			b.xSpeed = 1
		} else {
			b.xSpeed = -1
		}
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
	op := ebiten.DrawImageOptions{}

	// draw rider
	if b.state != UNMOUNTED {
		y := 0
		op.GeoM.Translate(float64(4), float64(y))
		rider := ebiten.NewImageFromImage(b.bounder)
		if b.state == SPAWNING {
			col := app.SpawnColors[rand.Intn(3)]
			rider = b.drawSolid(rider.Bounds(), col, rider)
		}
		composite.DrawImage(rider, &op)
		op.GeoM.Reset()
	}
	composite.DrawImage(body, &op)
	if !b.FacingRight {
		return b.flipX(composite, op)
	}
	return composite
}

func (b *Buzzard) mounted(gs *GameState) {
	b.doFlap()
	b.image = b.buildMount()

	b.velocity()
	b.Wrap()

	b.Collisions(gs.CliffAsSprites(), func(c *Sprite) {
		if b.Y < c.centerY() && xBetween(b.X, c.rect(), 3) {
			// buzzard is above
			b.Vy = 0.5
			b.Y = c.Y - float64(b.Height/2)
			b.walking = true
		} else if b.Y-b.Vy > c.Y && xBetween(b.X, c.rect(), 0) {
			// buzzard is below
			b.Y += 3
			b.Vy = 0.5
		} else if b.centerX() < c.centerX() {
			// buzzard is to left
			b.X -= 6
			b.xSpeed = -b.xSpeed
			b.FacingRight = false
		} else if b.centerX() > c.centerX() {
			// buzzard is to right
			b.X += 6
			b.xSpeed = -b.xSpeed
			b.FacingRight = true
		}
	})

	b.Collisions(gs.BuzzardsAsSprites(), func(c *Sprite) {
		b.bounce(gs, c)
	})
}

func (b *Buzzard) unmounted(gs *GameState) {
	if b.X < app.ScreenWidth/2 {
		b.FacingRight = false
		b.xSpeed = -3
	} else {
		b.FacingRight = true
		b.xSpeed = 3
	}
	b.doFlap()
	b.image = b.buildMount()
	b.velocity()
	if b.X < -float64(b.Width) || b.X > app.ScreenWidth+float64(b.Width/2) {
		for i, buzz := range gs.Buzzards {
			if buzz == b {
				b.state = DEAD
				gs.Buzzards = remove(gs.Buzzards, i)
				break
			}
		}
	}
}

func (b *Buzzard) bounce(gs *GameState, collider *Sprite) bool {
	above := false
	if b.Y < collider.centerY() && xBetween(b.X, collider.rect(), 3) {
		// buzzard is above
		b.Vy = 0.5
		b.Y = collider.Y - float64(b.Height/2)
		b.walking = true
		above = true
	} else if b.Y-b.Vy > collider.Y && xBetween(b.X, collider.rect(), 0) {
		// buzzard is below
		b.Y += 3
		b.Vy = 0.5
	} else if b.centerX() < collider.centerX() {
		// buzzard is to left
		b.X -= 5
		b.xSpeed = -2
	} else if b.centerX() > collider.centerX() {
		// buzzard is to right
		b.X += 5
		b.xSpeed = 2
	}
	return above
}

func remove(s []*Buzzard, i int) []*Buzzard {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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

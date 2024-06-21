package entity

import (
	"fmt"
	"github.com/depsypher/gojoust/app"
	"github.com/depsypher/gojoust/assets/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math/rand"
	"time"
)

type PlayerState int

const (
	SPAWNING  PlayerState = iota
	MOUNTED   PlayerState = iota
	UNMOUNTED PlayerState = iota
	DEAD      PlayerState = iota
)

type Player struct {
	*MountSprite
	rider       *ebiten.Image
	lastAnimate time.Time
	lastAccel   time.Time
	skid        time.Time
	xSpeed      int
	walkStep    bool
	state       PlayerState
}

func MakePlayer(ss *Sheet) *Player {
	return &Player{
		MountSprite: MakeMountSprite(ss.Ostrich),
		rider:       ss.P1Rider,
		lastAnimate: time.Now(),
		lastAccel:   time.Now(),
		skid:        time.Time{},
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	s := p.Sprite
	s.DrawSprite(screen)
}

func (p *Player) Update(gs *GameState) {
	switch p.state {
	case SPAWNING:
		p.spawning(gs)
	case MOUNTED:
		p.mounted(gs)
	case UNMOUNTED:
		p.unmounted(gs)
	case DEAD:
		p.dead(gs)
	}
}

func (p *Player) spawning(gs *GameState) {
	if time.Now().Before(p.lastAnimate.Add(time.Millisecond * time.Duration(30))) {
		return
	}
	if p.spawn <= 20 {
		// emerging
		p.buildSpawn(p, p.spawn)
		p.spawn += 1
		if p.spawn == 20 {
			if err := gs.Sounds[audio.EnergizeSound].Play(gs.SoundOn); err != nil {
				log.Fatal("Error playing sound", err)
			}
		}
	} else if p.spawn < 100 {
		// energizing/waiting
		if gs.Keys[app.FlapButton] || gs.Keys[app.LeftButton] || gs.Keys[app.RightButton] {
			p.state = MOUNTED
			p.image = p.buildMount()
			p.spawn = 0
			p.Vy = 1
			gs.Sounds[audio.EnergizeSound].Stop()
		} else {
			p.image = p.buildMount()
		}
		p.spawn += 1
	} else {
		p.state = MOUNTED
		p.image = p.buildMount()
		p.spawn = 0
		p.Vy = 1
	}
	p.lastAnimate = time.Now()
}

func (p *Player) mounted(gs *GameState) {
	p.walkInput(gs)
	p.flapInput(gs)
	p.velocity(gs)

	aboveCliff := false
	for _, c := range gs.CliffAsSprites() {
		p.Y += 1
		if p.Collides(c) {
			p.Y -= 1
			if p.bounce(gs, c) {
				aboveCliff = true
			}
		} else {
			p.Y -= 1
		}
	}

	p.walkAnimation(gs)
	p.Wrap()
	if !aboveCliff {
		p.walking = false
	}

	p.image = p.buildMount()
}

func (p *Player) bounce(gs *GameState, collider *Sprite) bool {
	above := false
	playBump := false
	if p.Y < collider.centerY() && xBetween(p.X, collider.rect(), 3) {
		// player is above
		p.Vy = 0.5
		p.Y = collider.Y - float64(p.Height/2)
		p.walking = true
		above = true
	} else if p.Y-p.Vy > collider.Y && xBetween(p.X, collider.rect(), 0) {
		// player is below
		p.Y += 3
		p.Vy = 0.5
		playBump = true
	} else if p.centerX() < collider.centerX() {
		// player is to left
		p.X -= 5
		p.xSpeed = -2
		playBump = true
	} else if p.centerX() > collider.centerX() {
		// player is to right
		p.X += 5
		p.xSpeed = 2
		playBump = true
	}
	if playBump {
		if err := gs.Sounds[audio.BumpSound].Play(gs.SoundOn); err != nil {
			log.Fatal("Error playing sound", err)
		}
	}
	return above
}

func xBetween(x float64, rect image.Rectangle, grace int) bool {
	return x <= float64(rect.Max.X-grace) && x >= float64(rect.Min.X+grace)
}

func (p *Player) velocity(gs *GameState) {
	if p.walking {
		if p.xSpeed != 0 {
			p.facingRight = p.xSpeed > 0
		}
	} else {
		p.Fall()
	}

	if p.xSpeed < -4 {
		p.xSpeed = -4
	} else if p.xSpeed > 4 {
		p.xSpeed = 4
	}
	if p.xSpeed < 0 {
		p.X -= app.MoveSpeed[-p.xSpeed]
	} else {
		p.X += app.MoveSpeed[p.xSpeed]
	}
	gs.Debug = fmt.Sprintf("xspeed=%d", p.xSpeed)
	p.Y += p.Vy
}

func (p *Player) walkInput(gs *GameState) {
	now := time.Now()
	canAccel := now.After(p.lastAccel.Add(time.Millisecond * time.Duration(120)))
	if !p.skid.IsZero() {
		if p.skid.After(now) {
			if p.xSpeed != 0 {
				speed := 4
				if p.skid.Sub(now) < (time.Millisecond * time.Duration(app.SkidMillis/2)) {
					speed = 2
				} else if p.skid.Sub(now) < (time.Millisecond * time.Duration(app.SkidMillis/2)) {
					speed = 3
				}
				if p.xSpeed > 0 {
					p.xSpeed = speed
				} else {
					p.xSpeed = -speed
				}
			}
		} else {
			p.xSpeed = 0
			p.lastAccel = time.Now()
			p.skid = time.Time{}
		}
	} else if p.walking && (p.xSpeed > 3 && gs.Keys[app.LeftButton] || (p.xSpeed < -3 && gs.Keys[app.RightButton])) {
		p.skid = now.Add(time.Millisecond * time.Duration(app.SkidMillis))
		if err := gs.Sounds[audio.SkidSound].Play(gs.SoundOn); err != nil {
			log.Fatal("Error playing sound", err)
		}
	} else if gs.Keys[app.LeftButton] {
		if p.walking {
			if canAccel {
				p.Vx = -1
				if p.xSpeed > -4 {
					p.xSpeed -= 1
					p.lastAccel = time.Now()
				}
			}
		} else {
			p.facingRight = false
		}
	} else if gs.Keys[app.RightButton] && canAccel {
		if p.walking {
			if canAccel {
				p.Vx = 1
				if p.xSpeed < 4 {
					p.xSpeed += 1
					p.lastAccel = time.Now()
				}
			}
		} else {
			p.facingRight = true
		}
	}
}

func (p *Player) flapInput(gs *GameState) {
	if gs.Keys[app.FlapButton] {
		p.skid = time.Time{}
		if p.flap == 0 {
			if gs.Keys[app.LeftButton] {
				p.xSpeed -= 1
			}
			if gs.Keys[app.RightButton] {
				p.xSpeed += 1
			}
			p.Vy = -0.4
			p.flap = 2
			gs.Sounds.StopSounds()
			if err := gs.Sounds[audio.FlapDnSound].Play(gs.SoundOn); err != nil {
				log.Fatal("Error playing sound", err)
			}
		} else {
			p.flap = 1
		}
		p.walking = false
	} else {
		if p.flap == 1 {
			gs.Sounds.StopSounds()
			if err := gs.Sounds[audio.FlapUpSound].Play(gs.SoundOn); err != nil {
				log.Fatal("Error playing sound", err)
			}
		}
		p.flap = 0
	}
}

func (p *Player) walkAnimation(gs *GameState) {
	if p.walking {
		if p.xSpeed == 0 {
			p.Frame = 3
			gs.Sounds.StopSounds()
		} else {
			if !p.skid.IsZero() {
				p.Frame = 4
			} else {
				nextFrame := app.WalkAnimSpeed[app.Abs(p.xSpeed)-1]
				if time.Now().After(p.lastAnimate.Add(nextFrame)) {
					p.Frame += 1
					if p.Frame > 3 {
						p.Frame = 0
					}
					if p.Frame == 2 {
						snd := audio.Walk1Sound
						if p.walkStep {
							snd = audio.Walk2Sound
						}
						if err := gs.Sounds[snd].Play(gs.SoundOn); err != nil {
							log.Fatal("Error playing sound", err)
						}
						p.walkStep = !p.walkStep
					}
					p.lastAnimate = time.Now()
				}
			}
		}
	}
}

func (p *Player) unmounted(gs *GameState) {
}

func (p *Player) dead(gs *GameState) {
}

func (p *Player) buildMount() *ebiten.Image {
	if p.flap == 1 {
		p.Frame = 5
	} else if p.flap == 2 || !p.walking {
		p.Frame = 6
	}

	frame := p.Images[p.Frame]
	composite := ebiten.NewImage(frame.Bounds().Dx(), frame.Bounds().Dy())

	body := ebiten.NewImageFromImage(frame)
	if p.state == SPAWNING {
		col := app.SpawnColors[rand.Intn(3)]
		body = p.drawSolid(composite.Bounds(), col, body)
	}

	// draw rider
	op := ebiten.DrawImageOptions{}
	y := 0
	if p.Frame == 4 {
		y = 2
	}
	op.GeoM.Translate(float64(4), float64(y))
	rider := ebiten.NewImageFromImage(p.rider)
	if p.state == SPAWNING {
		col := app.SpawnColors[rand.Intn(3)]
		rider = p.drawSolid(rider.Bounds(), col, rider)
	}
	composite.DrawImage(rider, &op)

	// draw ostrich
	op.GeoM.Reset()
	composite.DrawImage(body, &op)
	if !p.facingRight {
		return p.flipX(composite, op)
	}
	return composite
}

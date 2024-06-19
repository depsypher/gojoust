package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Buzzard struct {
	*Sprite
	bounder     *ebiten.Image
	lastAnimate time.Time
}

func MakeBuzzard(ss *Sheet) *Buzzard {
	return &Buzzard{
		Sprite:      MakeSprite(ss.Buzzard),
		bounder:     ss.Bounder,
		lastAnimate: time.Now(),
	}
}

func (b *Buzzard) buildMount() *ebiten.Image {
	image := ebiten.NewImageFromImage(b.Images[b.Frame])
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(6), float64(0))
	image.DrawImage(b.bounder, &op)
	return image
}

func (b *Buzzard) Draw(screen *ebiten.Image) {
	s := b.Sprite
	s.DrawSprite(screen)
}

func (b *Buzzard) Update(g *GameState) {
	if time.Now().After(b.lastAnimate.Add(time.Millisecond * time.Duration(500))) {
		b.Frame += 1
		b.lastAnimate = time.Now()
	}
	b.Frame = (b.Frame) % len(b.Images)
	b.image = b.buildMount()
	b.X += 1
	if b.X > 300 {
		b.X = float64(-b.Width)
	}

	b.Collisions(g.CliffAsSprites(), func(c *Sprite) {
		//		fmt.Printf("colliding %s\n", c.rect())
	})
}

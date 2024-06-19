package entity

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Cliff struct {
	*Sprite
}

func MakeCliff(image *ebiten.Image, x float64, y float64) *Cliff {
	// not sure why I can't use image directly, but the alpha channel is reversed?!
	img := ebiten.NewImage(image.Bounds().Dx(), image.Bounds().Dy())
	op := ebiten.DrawImageOptions{}
	img.DrawImage(image, &op)

	result := &Cliff{
		Sprite: MakeSprite([]*ebiten.Image{img}, x, y),
	}
	result.image = img
	result.center = false
	return result
}

func MakeBottomCliff(image *ebiten.Image) *Cliff {
	img := ebiten.NewImage(300, 30)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(70), float64(0))
	img.DrawImage(image, &op)

	result := &Cliff{
		Sprite: MakeSprite([]*ebiten.Image{img}, -20, 178),
	}
	result.image = img
	result.center = false
	return result
}

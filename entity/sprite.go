package entity

import (
	"bytes"
	"github.com/depsypher/gojoust/app"
	"github.com/depsypher/gojoust/assets/images"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image"
	"image/color"
	"image/draw"
	_ "image/png"
)

type Recter interface {
	rect() image.Rectangle
}
type Collider interface {
	Recter
	Collides()
}

type Sprite struct {
	Images []*ebiten.Image
	image  *ebiten.Image
	Frame  int
	Width  int
	Height int
	X      float64
	Y      float64
	Vx     float64
	Vy     float64
	Alive  bool
	center bool
}

func MakeSprite(images []*ebiten.Image, pos ...float64) *Sprite {
	position := []float64{0, 0}
	if len(pos) == 2 {
		position = pos
	}
	return &Sprite{
		images,
		nil,
		0,
		images[0].Bounds().Dx(),
		images[0].Bounds().Dy(),
		position[0], position[1],
		0, 0,
		true,
		true,
	}
}

func (s *Sprite) drawSolid(bounds image.Rectangle, color color.Color, mask image.Image) *ebiten.Image {
	img := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	img.Fill(color)
	draw.DrawMask(img, bounds, img, image.Point{}, mask, image.Point{}, draw.Src)
	return img
}

func (s *Sprite) rect() image.Rectangle {
	if s.center {
		w := float64(s.Width) / 2
		h := float64(s.Height) / 2
		return image.Rect(int(s.X-w), int(s.Y-h), int(s.X+w), int(s.Y+h))
	}
	return image.Rect(int(s.X), int(s.Y), int(s.X)+s.Width, int(s.Y)+s.Height)
}

func (s *Sprite) SetPos(x float64, y float64) {
	s.X = x
	s.Y = y
}

func (s *Sprite) centerX() float64 {
	if s.center {
		return s.X
	}
	return s.X + float64(s.Width/2)
}
func (s *Sprite) centerY() float64 {
	if s.center {
		return s.Y
	}
	return s.Y + float64(s.Height/2)
}

func (s *Sprite) Fall() {
	s.Vy += 4 * app.TimeStepSec / 2
	s.Y += s.Vy //* app.TimeStepSec
}

func (s *Sprite) Wrap() {
	if s.X > app.ScreenWidth {
		s.X = float64(-s.Width)
	} else if s.X < float64(-s.Width) {
		s.X = app.ScreenWidth - 1
	}
}

func (s *Sprite) flipX(image *ebiten.Image, op ebiten.DrawImageOptions) *ebiten.Image {
	width := image.Bounds().Dx()
	op.GeoM.Reset()
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(width), 0)
	left := ebiten.NewImage(width, image.Bounds().Dy())
	left.DrawImage(image, &op)
	return left
}

func (s *Sprite) DrawSprite(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if s.center {
		op.GeoM.Translate(-float64(s.Width)/2, -float64(s.Height)/2)
	}
	op.GeoM.Translate(s.X, s.Y)
	if s.image != nil {
		screen.DrawImage(s.image, op)
	}
}

func (s *Sprite) Collides(c *Sprite) bool {
	intersect := s.rect().Intersect(c.rect())
	result := intersect != image.Rectangle{}
	if result {
		// check pixels
		hitBoxMinX := intersect.Min.X
		hitBoxMinY := intersect.Min.Y
		hitBoxMaxX := intersect.Max.X
		hitBoxMaxY := intersect.Max.Y

		for y := hitBoxMinY; y < hitBoxMaxY; y++ {
			for x := hitBoxMinX; x < hitBoxMaxX; x++ {
				alpha1 := s.image.RGBA64At(x-s.rect().Min.X, y-s.rect().Min.Y).A
				if alpha1 != 0 {
					alpha2 := c.image.RGBA64At(x-c.rect().Min.X, y-c.rect().Min.Y).A
					//fmt.Println(x-c.rect().Min.X, y-c.rect().Min.Y, alpha1, alpha2)
					if alpha2 != 0 {
						return true
					}
				}
			}
		}
	}
	return false
}

type doOnCollide func(c *Sprite)

func (s *Sprite) Collisions(group []*Sprite, onCollide doOnCollide) []*Sprite {
	var result []*Sprite
	for _, c := range group {
		if s.Collides(c) {
			onCollide(c)
			result = append(result, c)
		}
	}
	return result
}

type Sheet struct {
	P1Rider *ebiten.Image
	Ostrich []*ebiten.Image
	Buzzard []*ebiten.Image
	Bounder *ebiten.Image
	C1      *ebiten.Image
	C2      *ebiten.Image
	C3      *ebiten.Image
	C4      *ebiten.Image
	C5      *ebiten.Image
	C6      *ebiten.Image
	C7      *ebiten.Image
	C8      *ebiten.Image
}

func LoadSpriteSheet() (*Sheet, error) {
	img, _, err := image.Decode(bytes.NewReader(images.Spritesheet_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	spriteAt := func(x, y, w, h int) *ebiten.Image {
		return sheet.SubImage(image.Rect(x, y, x+w, y+h)).(*ebiten.Image)
	}

	spriteFramesAt := func(x, y, w, h, gap, count int) []*ebiten.Image {
		var result []*ebiten.Image
		for i := 0; i < count; i++ {
			var frameX = i*w + i*gap
			result = append(result, sheet.SubImage(image.Rect(x+frameX, y, x+frameX+w, y+h)).(*ebiten.Image))
		}
		return result
	}

	s := &Sheet{}
	s.P1Rider = spriteAt(58, 79, 12, 7)
	s.Ostrich = spriteFramesAt(348, 19, 16, 20, 5, 8)
	s.Buzzard = spriteFramesAt(191, 44, 20, 20, 3, 7)
	s.Bounder = spriteAt(58, 69, 12, 7)
	s.C1 = spriteAt(0, 19, 190, 30)
	s.C2 = spriteAt(385, 0, 64, 8)  // 315, 420    # mid-bottom
	s.C3 = spriteAt(82, 0, 88, 9)   // 250, 201    # mid-top
	s.C4 = spriteAt(0, 9, 50, 7)    // -60, 168    # top-left
	s.C5 = spriteAt(0, 0, 64, 7)    // 759, 168    # top-right
	s.C6 = spriteAt(173, 0, 80, 8)  // -50, 354    # bottom-left
	s.C7 = spriteAt(319, 0, 63, 7)  // 770, 354    # bottom-right
	s.C8 = spriteAt(254, 0, 58, 11) // 606, 330    # mid-right

	return s, nil
}

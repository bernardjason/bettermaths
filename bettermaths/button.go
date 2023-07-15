package bettermaths

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	rectangle       image.Rectangle
	x, y            float32
	pressed         bool
	enabled         bool
	image           *ebiten.Image
	pressedFunction func(*Game)
	anchorLeft      bool
	anchorRight     bool
}

func loadImage(imageFileName *[]byte) *ebiten.Image {

	img, _, err := image.Decode(bytes.NewReader(*imageFileName))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func NewButton(x float32, y float32, embeddedImage embed.FS, imageFilename string, onPress func(g *Game), anchorLeft bool, anchorRight bool) *Button {
	file, err := embeddedImage.ReadFile(imageFilename)
	if err != nil {
		log.Fatal(err)
	}

	b := &Button{
		x:               x,
		y:               y,
		image:           loadImage(&file),
		pressed:         false,
		enabled:         true,
		pressedFunction: onPress,
		anchorLeft:      anchorLeft,
		anchorRight:     anchorRight,
	}
	b.rectangle = image.Rect(int(x), int(y), int(x)+b.image.Bounds().Dx(), int(y)+b.image.Bounds().Dy())
	return b
}

func (button *Button) drawButton(screen *ebiten.Image, game *Game) {

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	if button.anchorLeft {
		button.x = 0
	}
	if button.anchorRight {
		button.x = float32(game.ScreenWidth) - float32(button.image.Bounds().Dx()) - 1
	}

	op.GeoM.Translate(float64(button.x), float64(button.y))

	screen.DrawImage(button.image, op)
}

func (button *Button) contains(x int, y int) bool {
	return x >= button.rectangle.Min.X && x <= int(float32(button.x)+float32(button.rectangle.Dx())) &&
		y >= button.rectangle.Min.Y && y <= int(float32(button.y)+float32(button.rectangle.Dy()))

}

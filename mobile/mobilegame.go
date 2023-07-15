package mobilegame

import (
	"image/color"

	"github.com/bernardjason/bettermaths/bettermaths"
	"github.com/hajimehoshi/ebiten/v2/mobile"
)

func init() {
	// yourgame.Game must implement ebiten.Game interface.
	// For more details, see
	// * https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Game

	colours := bettermaths.DesktopColour{
		DropBoxBackgroundColourNumber:    color.RGBA{0xff, 0x80, 0x80, 0xff},
		DropBoxBackgroundColourOperation: color.RGBA{0x80, 0x80, 0xff, 0xff},
		DropBoxArithmeticSignColour:      color.RGBA{0xff, 0xff, 0xff, 0xff},
		NumberBackgroundColour:           color.RGBA{0x80, 0x80, 0x80, 0xff},
		NumberColour:                     color.RGBA{0xff, 0xff, 0xff, 0xff},
		BackgroundColourSigns:            color.RGBA{0xff, 0x80, 0x80, 0xff},
		SignColourNumber:                 color.RGBA{0x0, 0xff, 0x80, 0xff},
		TimeTakenColour:                  color.RGBA{0xff, 0xe6, 0x00, 0xff},
	}

	game := bettermaths.NewGame(colours)

	mobile.SetGame(game)

}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}

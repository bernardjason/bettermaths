package main

import (
	"image/color"
	"log"
	"runtime"

	"github.com/bernardjason/bettermaths/bettermaths"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

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

	if runtime.GOOS == "js" {
		colours = bettermaths.DesktopColour{
			DropBoxBackgroundColourNumber:    color.RGBA{0x0, 0x0, 0xff, 0xff},
			DropBoxBackgroundColourOperation: color.RGBA{0x80, 0x80, 0xff, 0xff},
			DropBoxArithmeticSignColour:      color.RGBA{0xff, 0xff, 0xff, 0xff},
			NumberBackgroundColour:           color.RGBA{0x0, 0xff, 0x00, 0xff},
			NumberColour:                     color.RGBA{0x0, 0x0, 0x0, 0xff},
			BackgroundColourSigns:            color.RGBA{0x0, 0xff, 0xff, 0xff},
			SignColourNumber:                 color.RGBA{0x0, 0x0, 0xff, 0xff},
			TimeTakenColour:                  color.RGBA{0xff, 0xe6, 0x00, 0xff},
		}
	}

	thegame := bettermaths.NewGame(colours)
	ebiten.SetWindowResizable(true)
	if err := ebiten.RunGame(thegame); err != nil {
		log.Fatal(err)
	}

}

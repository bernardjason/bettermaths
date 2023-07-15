package bettermaths

import "embed"

//go:embed reset.png
var resetImage embed.FS

//go:embed new.png
var newImage embed.FS

func SetupGuiButtons(screenWidth int, screenHeight int) []*Button {

	button := []*Button{

		NewButton(
			float32(screenWidth)-64,
			0,
			resetImage,
			"reset.png",
			func(g *Game) {
				if !g.GameLogic.finished {
					SetupScreen(g.GameLogic, screenWidth, screenHeight)
					resetDropBoxes(g)
				}
			},
			false, true,
		),
		NewButton(
			0,
			0,
			newImage,
			"new.png",
			func(g *Game) {
				SetupScreen(g.GameLogic, screenWidth, screenHeight)
				clearNewGame(g)
			},
			true, false,
		),
	}

	return button
}

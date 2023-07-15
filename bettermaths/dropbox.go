package bettermaths

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type DropBox struct {
	forNumber        bool
	playingNumber    *PlayingNumber
	playingOperation *PlayingOperation
	available        bool
	hover            bool
	x, y             float32
	boxWidth         float32
	boxHeight        float32
	whereOutOfThree  int
}

type DropBoxes struct {
	//DropBoxes []*DropBox
	Workings []*Workings
	//boxWidth  float32
	//boxHeight float32
}

type Workings struct {
	dropBox  [4]*DropBox
	answered bool
}

func newDropBox(boxWidth float32, boxHeight float32, forNumber bool, whereOutOfThree int) *DropBox {
	return &DropBox{available: true, playingNumber: nil,
		playingOperation: nil, forNumber: forNumber,
		boxWidth: boxWidth, boxHeight: boxHeight,
		whereOutOfThree: whereOutOfThree,
	}
}

/*
func OLDSetupYourWorkingOutDropBoxes(screenWidth int, screenHeight int, startY int) DropBoxes {

	boxNumber := 1

	twoColumns := 1
	if screenWidth >= 640 {
		twoColumns = 2
	}
	boxes := 0
	gap := 4
	xs := [...]int{
		NUMBER_WIDTH - gap*2,
		gap + NUMBER_WIDTH*2,
		gap + NUMBER_WIDTH*3 + gap*3}

	dropBoxes := []*DropBox{}
	workings := []*Workings{}
	for y := startY; boxes <= 5 && y < screenHeight; y += NUMBER_HEIGHT + 10 {

		for xx := 0; boxes <= 5 && xx < screenWidth; xx += screenWidth / twoColumns {
			forNumber := true
			working := Workings{}
			i := 0
			w := 0
			for _, x := range xs {
				width := NUMBER_WIDTH * 1.25
				if !forNumber {
					width = NUMBER_WIDTH
				}
				dropBox := newDropBox(float32(x+xx), float32(y), float32(width), NUMBER_HEIGHT, forNumber, i, boxNumber)
				dropBoxes = append(dropBoxes, dropBox)
				working.dropBox[i] = dropBox
				i++
				w++
				forNumber = !forNumber
				boxNumber++
			}
			boxes = boxes + 1
			workings = append(workings, &working)
		}
	}

	fmt.Printf("drop boxes = %d\n", len(dropBoxes))

	//d := DropBoxes{Workings: workings, boxWidth: float32(NUMBER_WIDTH * 1.5), boxHeight: float32(NUMBER_HEIGHT)}
	d := DropBoxes{Workings: workings}

	return d
}
*/

func SetupYourWorkingOutDropBoxes() DropBoxes {

	dropBoxes := []*DropBox{}
	workings := []*Workings{}

	for boxSets := 1; boxSets <= 5; boxSets++ {

		working := Workings{}
		i := 0
		w := 0
		forNumber := true

		for boxes := 0; boxes < 3; boxes++ {
			width := NUMBER_WIDTH * 1.25
			if !forNumber {
				width = NUMBER_WIDTH
			}
			dropBox := newDropBox(float32(width), NUMBER_HEIGHT, forNumber, i)
			dropBoxes = append(dropBoxes, dropBox)
			working.dropBox[i] = dropBox
			i++
			w++
			forNumber = !forNumber
		}

		workings = append(workings, &working)
	}

	fmt.Printf("drop boxes = %d\n", len(dropBoxes))

	//d := DropBoxes{Workings: workings, boxWidth: float32(NUMBER_WIDTH * 1.5), boxHeight: float32(NUMBER_HEIGHT)}
	d := DropBoxes{Workings: workings}

	return d
}

func resetDropBoxes(game *Game) {
	for _, db := range game.DropBoxes.Workings {
		db.answered = false
		for _, dropdropbox := range db.dropBox {
			if dropdropbox != nil {
				dropdropbox.available = true
				dropdropbox.playingNumber = nil
				dropdropbox.playingOperation = nil
			}
		}
	}
}

func clearNewGame(game *Game) {
	resetDropBoxes(game)
	thechallenge := createTheChallenge()
	game.GameLogic = thechallenge

}

func isIn(dropBox *DropBox, width, height int, x, y int) bool {
	width = (width * 4) / 3
	height = (height * 4) / 3

	isItIn := x > int(dropBox.x)-width/2 && x < int(dropBox.x)+width/2 && y < int(dropBox.y)+height/2 && y > int(dropBox.y)-height/2

	//fmt.Printf("isIn %v   dropbox=%f,%f   w/h=%d,%d   x,y=%d,%d\n", isItIn, dropBox.x, dropBox.y, width, height, x, y)
	return isItIn
}

func numberOrArithmeticDropBoxAt(x int, y int, dropboxes *DropBoxes, isANumberBox bool) *DropBox {
	for _, w := range dropboxes.Workings {
		for _, d := range w.dropBox {
			if d != nil && d.forNumber == isANumberBox && isIn(d, int(d.boxWidth), int(d.boxHeight), x, y) && d.available {
				return d
			}
		}
	}
	return nil
}
func dropBoxAtHighlight(x int, y int, dropboxes *DropBoxes, isItANumber bool) {
	for _, w := range dropboxes.Workings {
		for _, d := range w.dropBox {
			if d != nil {
				if (d.forNumber == isItANumber) && d.available && isIn(d, int(d.boxWidth), int(d.boxHeight), x, y) {
					d.hover = true
				} else {
					d.hover = false
				}
			}
		}
	}

}
func dropBoxResetHighlight(dropboxes *DropBoxes) {
	for _, w := range dropboxes.Workings {
		for _, d := range w.dropBox {
			if d != nil {
				d.hover = false
			}
		}

	}

}
func olddrawDropBoxes(screen *ebiten.Image, dropboxes *DropBoxes, colours *DesktopColour, game *Game) {

	for _, w := range dropboxes.Workings {
		populated := 0
		//x := 0
		//y := 0
		for _, d := range w.dropBox {
			if d == nil {
				continue
			}
			if d.playingNumber != nil || d.playingOperation != nil {
				populated++
			}
			if d.forNumber {
				if d.hover {
					invert := colours.DropBoxBackgroundColourNumber
					invert.B = 128 //255 - invert.B
					invert.G = 128 //255 - invert.G
					invert.R = 128 //255 - invert.R
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, invert, false)
				} else {
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, colours.DropBoxBackgroundColourNumber, false)

				}

			} else {
				if d.hover {
					invert := colours.DropBoxBackgroundColourNumber
					invert.B = 128 //255 - invert.B
					invert.G = 128 //255 - invert.G
					invert.R = 128 //255 - invert.R
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, invert, false)
				} else {
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, colours.DropBoxBackgroundColourOperation, false)
				}
				if d.playingOperation != nil {
					s := OperationToString(d.playingOperation.Operation)
					text.Draw(screen, s, mplusBigFont, int(d.x-d.boxWidth/2+NUMBER_WIDTH/3), int(d.y)+NUMBER_HEIGHT/3, colours.DropBoxArithmeticSignColour)
				}
			}
		}
	}

}

func drawDropBoxes(screen *ebiten.Image, dropboxes *DropBoxes, colours *DesktopColour, game *Game) float32 {

	lineStart := 18
	var x float32 = float32(lineStart)
	var y float32 = NUMBER_HEIGHT * 4.5
	boxesDrawn := 0
	var lastY float32 = 0
	boxesOnARow := 0
	for _, w := range dropboxes.Workings {

		for _, d := range w.dropBox {
			if d == nil {
				continue
			}
			lastY = y + d.boxHeight*1.25
			x = x + d.boxWidth/2
			d.x = x
			d.y = y

			if d.forNumber {
				if d.hover {
					invert := colours.DropBoxBackgroundColourNumber
					invert.B = 128
					invert.G = 128
					invert.R = 128
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, invert, false)
				} else {
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, colours.DropBoxBackgroundColourNumber, false)

				}

			} else {
				if d.hover {
					invert := colours.DropBoxBackgroundColourNumber
					invert.B = 128
					invert.G = 128
					invert.R = 128
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, invert, false)
				} else {
					vector.DrawFilledRect(screen, d.x-d.boxWidth/2, d.y-d.boxHeight/2, d.boxWidth, d.boxHeight, colours.DropBoxBackgroundColourOperation, false)
				}
				if d.playingOperation != nil {
					s := OperationToString(d.playingOperation.Operation)
					text.Draw(screen, s, mplusBigFont, int(d.x-d.boxWidth/2+NUMBER_WIDTH/3), int(d.y)+NUMBER_HEIGHT/3, colours.DropBoxArithmeticSignColour)
				}
			}
			boxesDrawn++

			boxesOnARow++
			x = x + d.boxWidth/2
			if boxesDrawn%3 == 0 {
				x = x + d.boxWidth

				if x+d.boxWidth*3 > float32(game.ScreenWidth) || boxesOnARow >= 9 {
					x = float32(lineStart)
					y = y + d.boxHeight*1.25
					boxesOnARow = 0
				}
			}

		}
	}

	return float32(lastY)

}

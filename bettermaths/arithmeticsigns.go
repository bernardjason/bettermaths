package bettermaths

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var ids = 0

type ArithmeticBox struct {
	id               int
	playingOperation *PlayingOperation
	available        bool

	Gui DragAndDrop
}

type ArithmeticBoxes struct {
	ArithmeticBoxes []*ArithmeticBox
	boxWidth        float32
	boxHeight       float32
}

func (arithmeticBox *ArithmeticBox) ArithmeticDropped(dropBox *DropBox) {
	dropBox.playingOperation = NewPlayingOperation(arithmeticBox.playingOperation.Operation)
	fmt.Println("DROPPED!!!")
}

func (arithmeticBox *ArithmeticBox) In(x, y int, width, height int) bool {
	width = (width * 4) / 3
	height = (height * 4) / 3

	isItIn := x > int(arithmeticBox.Gui.x)-width/2 && x < int(arithmeticBox.Gui.x)+width/2 && y < int(arithmeticBox.Gui.y)+height/2 && y > int(arithmeticBox.Gui.y)-height/2

	//fmt.Printf("arithmetic %s isIn %v   dropbox=%f,%f   w/h=%d,%d   x,y=%d,%d\n", OperationToString(arithmeticBox.playingOperation.Operation), isItIn, arithmeticBox.Gui.x, arithmeticBox.Gui.y, width, height, x, y)
	return isItIn
}

func (s *ArithmeticBox) PickedUp() {
	s.Gui.originalX = s.Gui.x
	s.Gui.originalY = s.Gui.y
}

func newArithmeticBox(x float32, y float32, playingOperation *PlayingOperation) *ArithmeticBox {
	gui := DragAndDrop{x: x, y: y}
	ids++
	return &ArithmeticBox{id: ids, Gui: gui, available: true, playingOperation: playingOperation}
}

func SetupArithmeticBoxes(screenWidth int, screenHeight int, startY int) ArithmeticBoxes {
	width := 52
	height := 52
	gap := 16
	xs := [...]int{gap / 2, gap + width + gap*1, gap/2 + width*2 + gap*3, gap/2 + width*3 + gap*4}
	arithmeticBoxes := []*ArithmeticBox{}

	signs := [...]PlayingOperation{*NewPlayingOperation(0), *NewPlayingOperation(1), *NewPlayingOperation(2), *NewPlayingOperation(3)}
	i := 0

	for _, x := range xs {

		arithmeticBox := newArithmeticBox(float32(width)+float32(x), float32(startY), &signs[i])
		arithmeticBoxes = append(arithmeticBoxes, arithmeticBox)
		i++

	}

	fmt.Printf("arithmetic boxes = %d\n", len(arithmeticBoxes))

	d := ArithmeticBoxes{ArithmeticBoxes: arithmeticBoxes, boxWidth: float32(width), boxHeight: float32(height)}

	return d
}

func drawArithmeticBoxes(screen *ebiten.Image, arithmeticboxes *ArithmeticBoxes, colours *DesktopColour) {

	for _, d := range arithmeticboxes.ArithmeticBoxes {

		s := OperationToString(d.playingOperation.Operation)
		vector.DrawFilledRect(screen,
			d.Gui.x-arithmeticboxes.boxWidth/2, d.Gui.y-arithmeticboxes.boxHeight/2,
			arithmeticboxes.boxWidth, arithmeticboxes.boxHeight,
			colours.BackgroundColourSigns, false)
		text.Draw(screen, s, mplusBigFont, int(d.Gui.x-arithmeticboxes.boxWidth/3), int(d.Gui.y+arithmeticboxes.boxHeight/3), colours.SignColourNumber)

	}

}

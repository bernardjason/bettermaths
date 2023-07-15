package bettermaths

import (
	"fmt"

	"image/color"
	_ "image/png"
	"log"

	"time"

	//"math/rand"
	//"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	mplusWellDoneFont font.Face
	mplusBigFont      font.Face
	endMessage        = "Well done"
)

const NUMBER_WIDTH = 80
const NUMBER_HEIGHT = 64

type DragAndDrop struct {
	x, y         float32
	originalX    float32
	originalY    float32
	value        string
	canBeDragged bool
	dragging     bool
}

func SetupFonts() {

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72

	mplusWellDoneFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    72,
		DPI:     dpi,
		Hinting: font.Hinting(font.WeightExtraBold),
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

}

type DesktopColour struct {
	DropBoxBackgroundColourNumber    color.RGBA
	DropBoxBackgroundColourOperation color.RGBA
	DropBoxArithmeticSignColour      color.RGBA
	NumberBackgroundColour           color.RGBA
	NumberColour                     color.RGBA
	BackgroundColourSigns            color.RGBA
	SignColourNumber                 color.RGBA
	TimeTakenColour                  color.RGBA
}

type Game struct {
	ScreenWidth           int
	ScreenHeight          int
	GameLogic             *GameLogic
	Strokes               map[*Stroke]struct{}
	touchIDs              []ebiten.TouchID
	DropBoxes             DropBoxes
	ArithmeticBoxes       ArithmeticBoxes
	Buttons               []*Button
	tick                  int
	colours               DesktopColour
	positionArtefactsDone bool
	fred                  int
}

func (g *Game) Update() error {
	// Initialize the glyphs for special (colorful) rendering.

	if !g.positionArtefactsDone {
		var space float32 = float32(g.ScreenWidth) / 4
		x := space / 2

		for _, d := range g.ArithmeticBoxes.ArithmeticBoxes {
			d.Gui.x = x
			x = x + space
		}

		g.positionArtefactsDone = true
	}

	g.tick++

	for check := range g.GameLogic.Current_game {
		populated := g.GameLogic.Current_game[check]
		if populated == nil {
			continue
		}
		if populated.Number == g.GameLogic.Current_answer.Number {
			g.GameLogic.finished = true
		}

	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.fred = g.fred + 1
		fmt.Println(g.fred)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s := NewStroke(&MouseStrokeSource{})

		mapScreenTouchPressToObject(g, s, g.GameLogic.finished)

	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		s := NewStroke(&TouchStrokeSource{id})
		mapScreenTouchPressToObject(g, s, g.GameLogic.finished)

	}

	for s := range g.Strokes {

		g.dropOrHighlightWhereToDrop(s)
		if s.IsReleased() {
			delete(g.Strokes, s)
		}
	}

	for s := range g.Strokes {

		switch v := s.DraggingObject().(type) {
		case *PlayingNumber:
			if v != nil { // android ??
				v.Gui.x = float32(s.currentX) - NUMBER_WIDTH/4
				v.Gui.y = float32(s.currentY)
			}

		case *ArithmeticBox:
			if v != nil { // android ??
				v.Gui.x = float32(s.currentX) - NUMBER_WIDTH/2
				v.Gui.y = float32(s.currentY)
			}
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}

	}

	for _, w := range g.DropBoxes.Workings {
		if w.answered {
			continue
		}
		populated := 0
		for _, d := range w.dropBox {
			if d == nil {
				continue
			}

			if d.playingNumber != nil || d.playingOperation != nil {
				populated++
			}
		}
		if populated == 3 {
			w.answered = true
			w.dropBox[0].playingNumber.Gui.canBeDragged = false
			w.dropBox[1].available = false
			w.dropBox[2].playingNumber.Gui.canBeDragged = false
			_, answer := perform(float32(w.dropBox[0].playingNumber.Number), float32(w.dropBox[2].playingNumber.Number), w.dropBox[1].playingOperation.Operation)
			number := NewPlayingNumber(answer)
			number.Gui.x = w.dropBox[2].playingNumber.Gui.x + w.dropBox[2].boxWidth
			number.Gui.y = w.dropBox[2].playingNumber.Gui.y

			g.GameLogic.Current_game = append(g.GameLogic.Current_game, number)

		}

	}

	if !g.GameLogic.finished {
		timeTaken := time.Since(g.GameLogic.started) //.ParseDuration("1m30s")
		g.GameLogic.taken = fmt.Sprintf("Time %0.f", timeTaken.Seconds())
	}

	return nil
}

func mapScreenTouchPressToObject(g *Game, s *Stroke, finished bool) {
	if !finished {
		if number := g.numberAt(s.Position()); number != nil && number.Gui.canBeDragged {
			number.PickedUp()
			s.SetDraggingObject(number)
			g.Strokes[s] = struct{}{}
			return
		} else if sign := g.signAt(s.Position()); sign != nil {
			fmt.Println("picked up")
			sign.PickedUp()
			s.SetDraggingObject(sign)
			g.Strokes[s] = struct{}{}
			return
		}

	}
	for _, b := range g.Buttons {
		if b.contains(s.Position()) {
			b.pressedFunction(g)
		}
	}

}

func (g *Game) numberAt(x, y int) *PlayingNumber {

	for i := len(g.GameLogic.Current_game) - 1; i >= 0; i-- {
		s := g.GameLogic.Current_game[i]
		if s.In(x, y) {
			return s
		}
	}
	return nil
}
func (g *Game) signAt(x, y int) *ArithmeticBox {

	for i := len(g.ArithmeticBoxes.ArithmeticBoxes) - 1; i >= 0; i-- {
		s := g.ArithmeticBoxes.ArithmeticBoxes[i]

		if s.In(x, y, int(g.ArithmeticBoxes.boxWidth), int(g.ArithmeticBoxes.boxHeight)) {
			return s
		}
	}
	return nil
}

func SetupScreen(g *GameLogic, screenWidth int, screenHeight int) {
	g.Current_answer.Gui.x = 20
	g.Current_answer.Gui.y = NUMBER_HEIGHT

	const y = NUMBER_HEIGHT * 3.5
	left := screenWidth / 6 / 4
	between := (screenWidth - left) / 6
	for i, x := 0, left; x < screenWidth && i < MAX_NUMBERS_IN_RANGE; i, x = i+1, x+between {
		g.Current_game[i].Gui.x = float32(x)
		g.Current_game[i].Gui.y = float32(y)
		g.Current_game[i].Gui.canBeDragged = true

	}
	g.Current_game = g.Current_game[0:MAX_NUMBERS_IN_RANGE]

}

func (g *Game) dropOrHighlightWhereToDrop(stroke *Stroke) {
	stroke.Update()
	dropBoxResetHighlight(&g.DropBoxes)

	switch s := stroke.DraggingObject().(type) {
	case *PlayingNumber:
		if s != nil { // android
			dropBoxAtHighlight(stroke.currentX, stroke.currentY, &g.DropBoxes, true)

		}
	case *ArithmeticBox:
		if s != nil {
			if s.available == true {
				dropBoxAtHighlight(stroke.currentX, stroke.currentY, &g.DropBoxes, false)
			}
		}

	default:
	}

	if !stroke.IsReleased() {
		return
	}

	switch s := stroke.DraggingObject().(type) {
	case *PlayingNumber:
		if s != nil { // android
			dropboxAt := numberOrArithmeticDropBoxAt(int(s.Gui.x), int(s.Gui.y), &g.DropBoxes, true)
			if dropboxAt == nil {
				s.NumberDropped()
			} else {
				if dropboxAt.available {
					dropboxAt.available = false
					dropboxAt.playingNumber = s
					s.Gui.x = dropboxAt.x - dropboxAt.boxWidth/2 //- NUMBER_WIDTH //.DropBoxes.boxWidth/3
					s.Gui.y = dropboxAt.y + NUMBER_HEIGHT/4
					if dropboxAt.whereOutOfThree == 0 {
						if s.Number <= 9 { // justify right
							s.Gui.x = s.Gui.x + dropboxAt.boxWidth*0.6
						} else if s.Number <= 99 {
							s.Gui.x = s.Gui.x + dropboxAt.boxWidth/3
						}
					}
				}
			}
		}
	case *ArithmeticBox:
		if s != nil {
			dropboxAt := numberOrArithmeticDropBoxAt(int(s.Gui.x), int(s.Gui.y), &g.DropBoxes, false)
			if dropboxAt != nil {
				s.ArithmeticDropped(dropboxAt)
			}
			s.Gui.x = s.Gui.originalX
			s.Gui.y = s.Gui.originalY
		}
	default:
		fmt.Printf("I don't know about type %T!\n", s)
	}
	dropBoxResetHighlight(&g.DropBoxes)

	stroke.SetDraggingObject(nil)
}

func (g *Game) Draw(screen *ebiten.Image) {

	lastY := drawDropBoxes(screen, &g.DropBoxes, &g.colours, g)

	g.drawNumber(g.GameLogic.Current_answer.Gui.value,
		screen,
		int(g.GameLogic.Current_answer.Gui.x), int(g.GameLogic.Current_answer.Gui.y),
		g.colours.NumberBackgroundColour, g.colours.NumberColour, true, false)

	for i := 0; i < len(g.GameLogic.Current_game); i++ {

		g.drawNumber(g.GameLogic.Current_game[i].Gui.value,
			screen,
			int(g.GameLogic.Current_game[i].Gui.x), int(g.GameLogic.Current_game[i].Gui.y),
			g.colours.NumberBackgroundColour, g.colours.NumberColour, false, false)

	}

	for _, b := range g.Buttons {
		b.drawButton(screen, g)
	}

	if g.GameLogic.finished {
		tickWellDone(float32(g.ScreenWidth)/4, 300, float32(g.ScreenWidth)/2, float32(NUMBER_HEIGHT*2), screen)

	} else {
		drawArithmeticBoxes(screen, &g.ArithmeticBoxes, &g.colours)
	}

	g.drawTimeTaken(screen, lastY)

}

func tickWellDone(x float32, y float32, wellDoneXMiddle float32, wellDoneY float32, screen *ebiten.Image) {

	//vector.StrokeLine(screen, x, y, x+52, y+56, 20, color.RGBA{0, 255, 0, 255}, true)
	//vector.StrokeLine(screen, x+39, y+57, x+250, y-160, 20, color.RGBA{0, 255, 0, 255}, true)

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(wellDoneXMiddle-float32(len(endMessage))*18), float64(wellDoneY))
	op.ColorScale.Scale(1, 1, 0, 1)
	op.Filter = ebiten.FilterLinear

	text.DrawWithOptions(screen, endMessage, mplusWellDoneFont, op)
}

func (g *Game) drawTimeTaken(screen *ebiten.Image, lastY float32) {
	x := 14

	text.Draw(screen, g.GameLogic.taken, mplusBigFont, x, int(lastY), g.colours.TimeTakenColour)
}

func (g *Game) drawNumber(number string, screen *ebiten.Image, x int, y int, backgroundColour color.RGBA, textColour color.RGBA, centreScreenX bool, centreScreenY bool) {
	letters := []string{
		"",
		"X",
		"XX",
		"XXX",
		"XXXX",
		"XXXXX",
		"XXXXXX",
	}
	const pad = 5

	b := text.BoundString(mplusBigFont, letters[len(number)])
	width := float32(b.Dx())
	height := float32(b.Dy()) + pad
	if len(number) == 1 {
		width = width + pad
	}
	standout := 1
	borderColour := backgroundColour
	borderColour.A = borderColour.A / 2
	borderColour.R = 255

	if centreScreenX {
		x = g.ScreenWidth/2 - int(width/2)
	}
	if centreScreenY {
		x = g.ScreenHeight/2 - int(height/2)
	}

	vector.DrawFilledRect(screen, float32(b.Min.X+x-standout), float32(b.Min.Y+y-standout), width+float32(standout), height+float32(standout*2), borderColour, false)
	vector.DrawFilledRect(screen, float32(b.Min.X+x), float32(b.Min.Y+y), width-float32(standout), height-float32(standout), backgroundColour, false)
	text.Draw(screen, number, mplusBigFont, x+pad/2, y+int(pad/2), textColour)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth > 1280 {
		outsideWidth = 1280
	}
	g.ScreenWidth = outsideWidth
	g.ScreenHeight = outsideHeight
	return g.ScreenWidth, g.ScreenHeight
}

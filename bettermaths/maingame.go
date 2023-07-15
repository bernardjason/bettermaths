package bettermaths

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const MAX_NUMBERS_IN_RANGE = 6
const MAX_NUMBER = 10
const MIN_NUMBER = 2

type PlayingNumber struct {
	Number int
	Gui    DragAndDrop
}

type PlayingOperation struct {
	Operation int
	inUse     bool
	Gui       DragAndDrop
}

func NewPlayingOperation(operation int) *PlayingOperation {
	gui := DragAndDrop{value: fmt.Sprintf("%d", operation), dragging: false}
	return &PlayingOperation{Operation: operation, Gui: gui}
}
func NewPlayingNumber(value int) *PlayingNumber {
	gui := DragAndDrop{value: fmt.Sprintf("%d", value), dragging: false, canBeDragged: true}
	return &PlayingNumber{Number: value, Gui: gui}
}

type GameLogic struct {
	Current_game   []*PlayingNumber
	Operation      []*PlayingOperation
	Current_answer *PlayingNumber
	finished       bool
	started        time.Time
	taken          string
}

func (s *PlayingNumber) In(x, y int) bool {
	width := len(s.Gui.value) * NUMBER_WIDTH
	height := NUMBER_HEIGHT
	return x > int(s.Gui.x) && x < int(s.Gui.x)+width && y < int(s.Gui.y)+height/2 && y > int(s.Gui.y)-height/2
}

func (s *PlayingNumber) NumberDropped() {
	s.Gui.x = s.Gui.originalX
	s.Gui.y = s.Gui.originalY
}
func (s *PlayingNumber) PickedUp() {
	s.Gui.originalX = s.Gui.x
	s.Gui.originalY = s.Gui.y
}

func NewChallenge() *GameLogic {

	g := GameLogic{
		Current_game: make([]*PlayingNumber, MAX_NUMBERS_IN_RANGE),
		Operation:    make([]*PlayingOperation, MAX_NUMBERS_IN_RANGE-1),
		finished:     false,
		started:      time.Now(),
		//guesses:      make([]*PlayingNumber, MAX_NUMBERS_IN_RANGE*2),
	}
	return &g
}

var (
	screenWidth  = 393 //800 //411 //1024
	screenHeight = 851 //640 // 817 // 640
)

func NewGame(colours DesktopColour) *Game {

	fmt.Printf("Hello from %s\n", runtime.GOOS)
	if runtime.GOOS == "js" {
		screenWidth, screenHeight = ebiten.ScreenSizeInFullscreen()
		//screenHeight = int(float32(screenHeight) * 0.8)

	} else if runtime.GOOS == "android" {
		screenWidth, screenHeight = 360, 700
	}

	fmt.Printf("Screen size %v %v", screenWidth, screenHeight)

	SetupFonts()
	thechallenge := createTheChallenge()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bettermaths")

	thegame := &Game{
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		GameLogic:       thechallenge,
		Strokes:         map[*Stroke]struct{}{},
		DropBoxes:       SetupYourWorkingOutDropBoxes(),
		ArithmeticBoxes: SetupArithmeticBoxes(screenWidth, screenHeight, NUMBER_HEIGHT*2),
		Buttons:         SetupGuiButtons(screenWidth, screenHeight),
		colours:         colours,
	}
	return thegame
}

func createTheChallenge() *GameLogic {
	game := NewChallenge()

	game.Setup(screenWidth, screenHeight)

	for game.Current_answer.Number <= 20 || game.Current_answer.Number > 999 {
		fmt.Println("too hard...")
		game.Setup(screenWidth, screenHeight)
	}

	fmt.Printf("Answer %v\n", game.Current_answer)

	fmt.Print("Numbers are ")
	for i := 0; i < MAX_NUMBERS_IN_RANGE; i++ {
		fmt.Print(game.Current_game[i].Number, "   ")
	}
	fmt.Println()
	fmt.Println()

	fmt.Println("My answer is ", game.Answer())
	return game
}

func (game *GameLogic) Setup(screenWidth int, screenHeight int) {

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < MAX_NUMBERS_IN_RANGE; i++ {
		number := NewPlayingNumber(rand.Intn(MAX_NUMBER-MIN_NUMBER) + MIN_NUMBER)
		number.Gui.canBeDragged = true
		game.Current_game[i] = number
	}
	for i := 0; i < MAX_NUMBERS_IN_RANGE-1; i++ {
		operation := NewPlayingOperation(rand.Intn(4))
		game.Operation[i] = operation
	}

	answer := NewPlayingNumber(game.Current_game[0].Number)
	for i := 1; i < MAX_NUMBERS_IN_RANGE; i++ {

		b := game.Current_game[i]

		ok := false
		add_to := 0
		for !ok {
			operation := game.Operation[i-1]
			ok, add_to = perform(float32(answer.Number), float32(b.Number), operation.Operation)
			if !ok {
				operation := rand.Intn(4)
				game.Operation[i-1] = NewPlayingOperation(operation)
			}
		}
		answer = NewPlayingNumber(add_to)
	}
	game.Current_answer = answer

	/*
		for i := 0; i < MAX_NUMBERS_IN_RANGE; i++ {
			number := NewPlayingNumber(game.Current_game[i].Number)

			game.guesses = append(game.guesses, number)

		}
	*/
	SetupScreen(game, screenWidth, screenHeight)
}

/*
func TestCase2(game *GameLogic) {
	game.Current_answer = 52
	game.Current_game = []int{9, 3, 2, 5, 4, 2}
	game.Operation = []int{3, 2, 2, 1, 2}
}

func TestCase1(game *GameLogic) {
	game.Current_answer = 23
	game.Current_game = []int{5, 2, 3, 4, 5, 7}
	game.Operation = []int{0, 0, 1, 2, 1}
}
*/

func (game *GameLogic) Answer() string {

	finalAnswer := NewPlayingNumber(game.Current_game[0].Number)
	working := fmt.Sprintf("%d", finalAnswer.Number)
	for i := 1; i < MAX_NUMBERS_IN_RANGE; i++ {

		b := game.Current_game[i]

		operation := game.Operation[i-1]

		var ok bool
		ok, finalAnswer.Number = perform(float32(finalAnswer.Number), float32(b.Number), operation.Operation)

		if !ok {
			panic("HOW CAN IT NOT BE ROUND NUMBER")
		}

		next := ""
		if i < MAX_NUMBERS_IN_RANGE-1 {
			next = OperationToString(game.Operation[i].Operation)
		}

		switch next {
		case "*", "/":
			working = fmt.Sprintf("%s %s %d = %d", working, OperationToString(operation.Operation), b.Number, finalAnswer.Number)
		default:
			working = fmt.Sprintf("%s %s %d", working, OperationToString(operation.Operation), b.Number)
		}

	}
	working = fmt.Sprintf("%s = %d  ", working, finalAnswer.Number)

	return working
}

func OperationToString(operation int) string {
	switch operation {
	case 0:
		return "+"
	case 1:
		return "\u002d"
	case 2:
		return "\u00D7"
	case 3:
		return "\u00F7"
	}
	panic("Unrecognised operation to string")
}

func perform(a float32, b float32, operation int) (bool, int) {
	var result float32 = 0.0
	switch operation {
	case 0:
		result = a + b
	case 1:
		result = a - b
	case 2:
		result = a * b
	case 3:
		result = a / b
	}

	//fmt.Println(a, b, game.OperationToString(operation))
	if math.Round(float64(result)) != float64(result) {
		return false, 0
	}

	return true, int(result)
}

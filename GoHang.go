package main

// Packages are "namespaces and compilation units"
// Every file in Go belongs to a package

import (
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// --- Globals ---

// This is a pragma to tell compiler to read these contents into the binary
// At compile time
// Also words.txt came from https://github.com/dwyl/english-words/blob/master/words_alpha.txt
// May need a better dictionary, but one thing at a time
//
//go:embed words.txt
var wordData string

//go:embed help.txt
var helpText string

// "Game implements ebiten.Game interface. ebiten.Game has necessary functions for an Ebitengine game: Update, Draw and Layout"
type Game struct {
	word    string
	width   int
	height  int
	state   int
	lastEsc time.Time
	// So, rune is kind of an interesting data type. Go doesn't have a char type, only rune, which is actually stored as an int32
	// ChatGPT: "var c rune = 'a' is actually var c int32 = 97 because 'a' has Unicode codepoint 97."
	// In short there are no chars, only bytes (uint8) and runes (int32)
	// Strings are byte sequences, but some unicode characters are multiple bytes, necessitating runes
	// "s := "hello" Internally this is: 68 65 6c 6c 6f But Unicode characters can be multiple bytes. Example: é"
	// So, essentially, the following line just creates a map that accepts any unicode characters, and maps them to bools as the value to the key:value pair
	guessed map[rune]bool
	// Btw UTF-8 stands for Unicode Transformation Format - 8-bit, so it's shortened to just "unicode"
	incorrectGuesses int
}

const (
	menu = iota
	help
	game
	victory
	loss
)

// "iota" make this equivalent to:
// const Menu State = 0
// const Help State = 1
// const Gameplay State = 2

var startRect = image.Rect(300, 200, 500, 240)
var helpRectR = image.Rect(300, 260, 500, 300)
var quitRect = image.Rect(300, 320, 500, 360)

// --- Main ---

// Function order does not matter, can put other funcs later
func main() {
	// Using Ebitengine to draw window & run game
	sw, sh := ebiten.Monitor().Size()

	w := int(float64(sw) * 0.8)
	h := int(float64(sh) * 0.8)

	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("GoHang")

	if err := ebiten.RunGame(NewGame(w, h)); err != nil {
		log.Fatal(err)
	}

}

// --- Create New Game ---

// WARNING: Vibe Coded
func NewGame(w, h int) *Game {
	// Returning a pointer to a new Game struct

	// words := strings.Split(strings.TrimSpace(wordData), "\n")
	// TrimSpace removes leading/trailing whitespace for the file, but isn't recognizing DOS /r/n endlines, only /n
	// Fields splits on any whitespace
	words := strings.Fields(wordData)
	return &Game{
		word:    words[rand.Intn(len(words))],
		width:   w,
		height:  h,
		state:   menu,
		lastEsc: time.Now(),
		guessed: make(map[rune]bool),
	}
}

// --- For Ebitengine ---

// "Defines (*Game).Update function, that is called every tick. Tick is a time unit for logical updating.
// The default value is 1/60 [s], then Update is called 60 times per second by default"
//
// This is where most of the work happens
func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()

	if g.state == menu {

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

			if inside(x, y, startRect) {
				println("Game")
				g.state = game
			}

			if inside(x, y, helpRectR) {
				println("Help")
				g.state = help
			}

			if inside(x, y, quitRect) {
				println("Quit")
				return ebiten.Termination
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			now := time.Now()

			if now.Sub(g.lastEsc) > 2*time.Second {
				println("ESC")
				return ebiten.Termination
			}
		}
	}

	if g.state == help {
		g.back2Menu()
	}

	if g.state == game {
		// Check for ESC key press to quit to menu
		g.back2Menu()

		g.gameKeyPressChecker()

		// Victory condition checker
		if g.allLettersGuessed() {
			g.state = victory
			print("Victory!\n")
		}
		// Defeat condition checker
		if g.incorrectGuesses >= 6 {
			g.state = loss
			print("Defeat!\n")
		}
	}

	if g.state == victory {
		g.back2Menu()
	}

	if g.state == loss {
		// Seems 2 work
		g.back2Menu()
	}

	return nil
}

// Helper function to check if a point is inside a rectangle
func inside(x, y int, r image.Rectangle) bool {
	return x >= r.Min.X && x <= r.Max.X && y >= r.Min.Y && y <= r.Max.Y
}

// Check for victory condition: all letters guessed correctly
func (g *Game) allLettersGuessed() bool {
	for _, c := range g.word {
		if !g.guessed[c] {
			return false
		}
	}
	return true
}

// "Defines (*Game).Draw function, that is called every frame" "Takes an *ebiten.Image as an argument"
//
//	func (g *Game) Draw(screen *ebiten.Image) {
//		ebitenutil.DebugPrint(screen, "Random word: "+g.word)
//	}
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {

	case menu:
		ebitenutil.DebugPrintAt(screen, "[ Start Game ]", 320, 210)
		ebitenutil.DebugPrintAt(screen, "[ Help ]", 350, 270)
		ebitenutil.DebugPrintAt(screen, "[ Quit ]", 350, 330)

	case help:
		ebitenutil.DebugPrint(screen, helpText)

	case game:

		g.drawGallows(screen)

		g.drawWordAndInfo(screen)

		// Call to drawHangman helper function to draw the hangman figure
		g.drawHangman(screen)

	case victory:
		// ebitenutil.DebugPrint(screen, "Random word: "+g.word)

		g.drawGallows(screen)

		g.drawWordAndInfo(screen)

		// Call to drawHangman helper function to draw the hangman figure
		g.drawHangman(screen)

		victoryMsg := "You won!"
		g.drawGameOverScreen(screen, victoryMsg)

	case loss:
		// ebitenutil.DebugPrint(screen, "Random word: "+g.word)

		g.drawGallows(screen)

		g.drawWordAndInfo(screen)

		// Call to drawHangman helper function to draw the hangman figure
		g.drawHangman(screen)
		lossMsg := "You lost! :("
		g.drawGameOverScreen(screen, lossMsg)

	}
}

// "Defines (*Game).Layout function. Layout accepts an outside size, which is a window size on desktop, and returns the game's logical screen size."
func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return g.width, g.height // Vibe coded
}

// By the way, didn't notice this before but the (g *Game) thing indicates that drawHangman is a method of the Game struct
func (g *Game) drawHangman(screen *ebiten.Image) {
	// Check number of incorrect guesses and draw the corresponding hangman figure
	if g.incorrectGuesses >= 1 {
		// Head
		vector.StrokeCircle(screen, 302, 370, 30, 2, color.White, false)
		// print("Drawing Head")
	}
	if g.incorrectGuesses >= 2 {
		// Body
		vector.StrokeLine(screen, 302, 400, 302, 500, 2, color.White, false)
		// print("Drawing Body")
	}
	if g.incorrectGuesses >= 3 {
		// Left arm
		vector.StrokeLine(screen, 302, 420, 260, 460, 2, color.White, false)
		// print("Drawing Left Arm")
	}
	if g.incorrectGuesses >= 4 {
		// Right arm
		vector.StrokeLine(screen, 302, 420, 344, 460, 2, color.White, false)
		// print("Drawing Right Arm")
	}
	if g.incorrectGuesses >= 5 {
		// Left leg
		vector.StrokeLine(screen, 302, 500, 260, 560, 2, color.White, false)
		// print("Drawing Left Leg")
	}
	if g.incorrectGuesses >= 6 {
		// Right leg
		vector.StrokeLine(screen, 302, 500, 344, 560, 2, color.White, false)
		// print("Drawing Right Leg")
	}
	// Thank you claude, this would have taken me forever to trial and error
}

// Optimization / cleanup functions

// Only used for game, victory, loss, or help states, NOT for menu state
func (g *Game) back2Menu() {

	// Check if the player just pressed ESC
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		now := time.Now()

		// Compare current time with last ESC press time to prevent auto quitting the game
		if now.Sub(g.lastEsc) > 2*time.Second {
			println("ESC")
			// Set the game state and reset game data
			g.state = menu
			g.lastEsc = now
			words := strings.Fields(wordData)
			g.word = words[rand.Intn(len(words))]
			g.guessed = make(map[rune]bool)
			g.incorrectGuesses = 0
		}
	}
}

func (g *Game) gameKeyPressChecker() {
	for k := ebiten.KeyA; k <= ebiten.KeyZ; k++ {
		if ebiten.IsKeyPressed(k) {

			// This is weird because the way Go handles characters is via runes (UTF-8 code points)
			// so this acts like a int + offset deal which converts any key to lowercase
			// From claude:
			// "ebiten.KeyA is just an integer constant representing the A key — let's say it's 65 (the actual value doesn't matter, just that the keys A through Z are consecutive integers)
			// k - ebiten.KeyA gives you the offset of whichever key you're on — so KeyA gives 0, KeyB gives 1, KeyC gives 2, and so on up to KeyZ giving 25
			// 'a' + offset then adds that offset to the rune value of lowercase 'a' (which is 97 in Unicode), giving you 97 for 'a', 98 for 'b', 99 for 'c', etc.
			// The outer rune(...) is just an explicit type cast to make Go happy, since the arithmetic produces a plain integer"
			letter := rune('a' + (k - ebiten.KeyA))

			// If it's not in the guessed map, add it
			if !g.guessed[letter] {
				g.guessed[letter] = true
				println("Guessed:", string(letter))
				// Incorrect guess handling
				if !g.correctGuess(letter) {
					g.incorrectGuesses++
				}
			}
		}
	}
}

// Checking if a letter is in the word
func (g *Game) correctGuess(letter rune) bool {
	for _, c := range g.word {
		if c == letter {
			return true
		}
	}
	return false
}

func (g *Game) drawGallows(screen *ebiten.Image) {
	// Drawing gallows
	vector.DrawFilledRect(screen, 100, 500, 200, 10, color.White, false)
	vector.DrawFilledRect(screen, 180, 300, 10, 200, color.White, false)
	vector.DrawFilledRect(screen, 180, 300, 120, 10, color.White, false)
	vector.DrawFilledRect(screen, 300, 300, 2, 40, color.White, false)
}

func (g *Game) drawWordAndInfo(screen *ebiten.Image) {
	// Info for user
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Word length: %d", len(g.word)), 50, 50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Incorrect Guesses: %d", g.incorrectGuesses), 50, 70)
	// This line creates the word spaces on the screen, but generates one extra space after the word
	ebitenutil.DebugPrintAt(screen, g.displayWord(), 300, 100)
}

// This is called when a letter is guessed
func (g *Game) displayWord() string {

	result := ""

	// "range on a string iterates over it character by character, yielding two values on each iteration
	// — the byte index of the character, and the character itself as a rune. The _ is Go's blank identifier,
	// meaning "I don't care about this value, throw it away." So _, c means "give me the character but discard the index."
	// c then holds each rune in the word one at a time."
	for _, c := range g.word {

		if g.guessed[c] {
			result += string(c) + " "
		} else {
			result += "_ "
		}
	}

	// Fix the trailing space before returning
	strings.TrimRight(result, " ")
	return result
}

func (g *Game) drawGameOverScreen(screen *ebiten.Image, finalMsg string) {

	wordMsg := fmt.Sprintf("The word was: %s", g.word)
	guessMsg := fmt.Sprintf("Incorrect guesses: %d/6", g.incorrectGuesses)
	escMsg := "Press ESC to return to menu"

	ebitenutil.DebugPrintAt(screen, finalMsg, (g.width-len(finalMsg)*6)/2, g.height-80)
	ebitenutil.DebugPrintAt(screen, wordMsg, (g.width-len(wordMsg)*6)/2, g.height-60)
	ebitenutil.DebugPrintAt(screen, guessMsg, (g.width-len(guessMsg)*6)/2, g.height-40)
	ebitenutil.DebugPrintAt(screen, escMsg, (g.width-len(escMsg)*6)/2, g.height-20)
}

// didn't actually save that much space by refactoring, lol

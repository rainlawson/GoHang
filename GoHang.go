package main

// Packages are "namespaces and compilation units"
// Every file in Go belongs to a package

import (
	_ "embed"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// This is a pragma to tell compiler to read these contents into the binary
// At compile time
// Also words.txt came from https://github.com/dwyl/english-words/blob/master/words_alpha.txt
// May need a better dictionary, but one thing at a time
//
//go:embed words.txt
var wordData string

// Function order does not matter, can put other funcs later
// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	words := strings.Split(strings.TrimSpace(wordData), "\n")

// 	random := words[rand.Intn(len(words))]
// 	fmt.Println(random)
// }

type Game struct {
	word   string
	width  int
	height int
}

func NewGame(w, h int) *Game {
	rand.Seed(time.Now().UnixNano())

	words := strings.Split(strings.TrimSpace(wordData), "\n")
	return &Game{
		word:   words[rand.Intn(len(words))],
		width:  w,
		height: h,
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Random word: "+g.word)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return g.width, g.height // critical fix
}

func main() {
	sw, sh := ebiten.Monitor().Size()

	w := int(float64(sw) * 0.8)
	h := int(float64(sh) * 0.8)

	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("GoHang")

	if err := ebiten.RunGame(NewGame(w, h)); err != nil {
		log.Fatal(err)
	}
}

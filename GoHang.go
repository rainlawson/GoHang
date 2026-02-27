package main

// Packages are "namespaces and compilation units"
// Every file in Go belongs to a package

import (
	_ "embed"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// This is a pragma to tell compiler to read these contents into the binary
// At compile time
// Also words.txt came from https://github.com/dwyl/english-words/blob/master/words_alpha.txt
// May need a better dictionary, but one thing at a time
//
//go:embed words.txt
var wordData string

// Function order does not matter, can put other funcs later
func main() {
	rand.Seed(time.Now().UnixNano())

	words := strings.Split(strings.TrimSpace(wordData), "\n")

	random := words[rand.Intn(len(words))]
	fmt.Println(random)
}

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

type Global struct {
	alphabet string
	taille   int
	N        uint64
}

func newGlobal(alphabet string, taille int) Global {
	n := uint64(1)
	for i := 0; i < taille; i++ {
		n *= uint64(len(alphabet))
	}
	return Global{
		alphabet: alphabet,
		taille:   taille,
		N:        n,
	}
}

func sha1Hex(texte string) []byte {
	sum := sha1.Sum([]byte(texte))
	return sum[:]
}

func emprinteToString(empreinte []byte) string {
	return strings.ToUpper(hex.EncodeToString(empreinte))
}

func main() {
	var global = newGlobal("abcdefghijklmnopqrstuvwxyz", 4)

	// QUESTION 1
	fmt.Println("Question1|SHA1|: Salut ->", emprinteToString(sha1Hex("Salut"))) // 9F57098C5534762DD32802302DB78ADA1BA864F5

	// QUESTION 2
	fmt.Println("Question2|N|", global.N)
}
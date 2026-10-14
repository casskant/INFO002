package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Global struct {
	hashName string 
	alphabet string
	taille   int  
	N        uint64
}

func newGlobal(hashName string, alphabet string, taille int) Global {
	// N = (nombre de lettres de l'alphabet) puissance taille
	nbLettres := uint64(len(alphabet))
	N := uint64(1)
	for i := 0; i < taille; i++ {
		N = N * nbLettres
	}

	return Global{
		hashName: hashName,
		alphabet: alphabet,
		taille:   taille,
		N:        N,
	}
}

// ---------- Question 1 : hachage ----------

func sha1Hex(texte string) []byte {
	octetsDuTexte := []byte(texte)
	empreinte := sha1.Sum(octetsDuTexte) 
	return empreinte[:]                  
}

func emprinteToString(empreinte []byte) string {
	enHexa := hex.EncodeToString(empreinte) 
	return strings.ToUpper(enHexa)         
}

// ---------- Question 2 : configuration ----------

func afficherConfig(g Global) {
	fmt.Println("Fonction de hachage :", g.hashName)
	fmt.Println("Alphabet            :", g.alphabet)
	fmt.Println("Taille des clairs   :", g.taille)
	fmt.Println("Nombre de clairs N  :", g.N)
}

// ---------- Question 3 : indice -> texte clair ----------

// i2c transforme un indice n en texte clair :
// c'est l'écriture de n en base (nombre de lettres de l'alphabet).
func i2c(g Global, n uint64, texte []byte) {
	nbLettres := uint64(len(g.alphabet))

	for position := g.taille - 1; position >= 0; position-- {
		numeroLettre := n % nbLettres
		texte[position] = g.alphabet[numeroLettre]
		n = n / nbLettres
	}
}

// ---------- Question 4 : empreinte -> indice ----------

// h2i transforme une empreinte y en indice : (y[0..7] + t) % N.
func h2i(g Global, y []byte, t uint64) uint64 {
	// On lit les 8 premiers octets en little-endian :
	// y[0] est l'octet de poids faible (il vaut y[0] * 1),
	// y[1] vaut y[1] * 256, y[2] vaut y[2] * 256*256, etc.
	valeur := uint64(0)
	poids := uint64(1)
	for i := 0; i < 8; i++ {
		valeur = valeur + uint64(y[i])*poids
		poids = poids * 256
	}

	return (valeur + t) % g.N
}

// ---------- Lecture des fichiers test ----------
func lireLignes(nomFichier string) []string {
	fichier, err := os.Open(nomFichier)
	if err != nil {
		fmt.Println("Impossible d'ouvrir le fichier :", nomFichier)
		os.Exit(1)
	}
	defer fichier.Close()

	var lignes []string
	scanner := bufio.NewScanner(fichier)
	for scanner.Scan() {
		lignes = append(lignes, scanner.Text())
	}
	return lignes
}

// executerCommande lance la commande demandée (ligne 4) avec ses arguments (lignes 5 et suivantes).
func executerCommande(g Global, commande string, arguments []string) {
	switch commande {

	case "hash":
		texte := arguments[0]
		empreinte := sha1Hex(texte)
		fmt.Println("hash("+texte+") =", emprinteToString(empreinte))

	case "config":
		// rien à faire : la configuration est déjà affichée pour chaque fichier test

	case "i2c":
		n, _ := strconv.ParseUint(arguments[0], 10, 64)
		texte := make([]byte, g.taille)
		i2c(g, n, texte)
		fmt.Println("i2c("+arguments[0]+") =", string(texte))

	case "h2i":
		texte := arguments[0]
		t, _ := strconv.ParseUint(arguments[1], 10, 64)
		empreinte := sha1Hex(texte)
		indice := h2i(g, empreinte, t)
		fmt.Println("hash("+texte+") =", emprinteToString(empreinte))
		fmt.Println("h2i(hash("+texte+"), "+arguments[1]+") =", indice)

	default:
		fmt.Println("Commande inconnue :", commande)
	}
}

func main() {
	// usage : go run main.go fichier_test [fichier_test ...]
	nomsFichiers := os.Args[1:]

	for _, nomFichier := range nomsFichiers {
		fmt.Println("=====", nomFichier, "=====")

		lignes := lireLignes(nomFichier)

		hashName := lignes[0]
		alphabet := lignes[1]
		taille, _ := strconv.Atoi(lignes[2])
		commande := lignes[3]
		arguments := lignes[4:]

		g := newGlobal(hashName, alphabet, taille)
		afficherConfig(g)
		fmt.Println("Commande            :", commande)
		fmt.Println("-----")

		executerCommande(g, commande, arguments)
		fmt.Println()
	}
}

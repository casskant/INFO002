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

// Global regroupe la configuration lue dans les 3 premières lignes d'un fichier test.
type Global struct {
	hashName string // fonction de hachage (ligne 1)
	alphabet string // alphabet des textes clairs (ligne 2)
	taille   int    // taille des textes clairs (ligne 3)
	N        uint64 // nombre total de textes clairs : len(alphabet)^taille
}

func newGlobal(hashName string, alphabet string, taille int) Global {
	n := uint64(1)
	for i := 0; i < taille; i++ {
		n *= uint64(len(alphabet))
	}
	return Global{
		hashName: hashName,
		alphabet: alphabet,
		taille:   taille,
		N:        n,
	}
}

// ---------- Question 1 : hachage ----------

func sha1Hex(texte string) []byte {
	sum := sha1.Sum([]byte(texte))
	return sum[:]
}

func emprinteToString(empreinte []byte) string {
	return strings.ToUpper(hex.EncodeToString(empreinte))
}

// ---------- Question 2 : configuration ----------

func afficherConfig(g Global) {
	fmt.Println("Fonction de hachage :", g.hashName)
	fmt.Printf("Alphabet            : %s (%d caractères)\n", g.alphabet, len(g.alphabet))
	fmt.Println("Taille des clairs   :", g.taille)
	fmt.Println("Nombre de clairs N  :", g.N)
}

// ---------- Question 3 : numérotation des textes clairs ----------

// i2c transforme un entier n (0 <= n < N) en texte clair :
// c'est l'écriture de n en base len(alphabet), sur exactement "taille" lettres.
// La dernière lettre correspond à n modulo len(alphabet).
func i2c(g Global, n uint64) string {
	base := uint64(len(g.alphabet))
	texte := make([]byte, g.taille)
	for i := g.taille - 1; i >= 0; i-- { // on remplit de droite à gauche
		texte[i] = g.alphabet[n%base]
		n /= base
	}
	return string(texte)
}

// ---------- Lecture des fichiers test ----------

// lireFichierTest lit un fichier test et renvoie la configuration,
// la commande (ligne 4) et ses arguments (lignes 5 et suivantes).
func lireFichierTest(nomFichier string) (Global, string, []string, error) {
	f, err := os.Open(nomFichier)
	if err != nil {
		return Global{}, "", nil, err
	}
	defer f.Close()

	// bufio.Scanner retire le "\n" et l'éventuel "\r" (fichiers Windows)
	var lignes []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lignes = append(lignes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return Global{}, "", nil, err
	}

	if len(lignes) < 4 {
		return Global{}, "", nil, fmt.Errorf("le fichier doit contenir au moins 4 lignes (%d trouvées)", len(lignes))
	}

	hashName := strings.TrimSpace(lignes[0])
	if hashName != "SHA1" {
		return Global{}, "", nil, fmt.Errorf("fonction de hachage non supportée : %q", hashName)
	}

	alphabet := strings.TrimSpace(lignes[1])

	taille, err := strconv.Atoi(strings.TrimSpace(lignes[2]))
	if err != nil {
		return Global{}, "", nil, fmt.Errorf("taille des textes clairs invalide : %q", lignes[2])
	}

	commande := strings.TrimSpace(lignes[3])
	arguments := lignes[4:]

	return newGlobal(hashName, alphabet, taille), commande, arguments, nil
}

// executerCommande lance la commande demandée avec ses arguments.
func executerCommande(g Global, commande string, arguments []string) error {
	switch commande {
	case "hash":
		if len(arguments) < 1 {
			return fmt.Errorf("la commande hash attend 1 argument (ligne 5)")
		}
		texte := arguments[0]
		fmt.Printf("%s(%q) = %s\n", g.hashName, texte, emprinteToString(sha1Hex(texte)))

	case "config":
		// la configuration est déjà affichée pour chaque fichier test
		fmt.Println("(configuration affichée ci-dessus)")

	case "i2c":
		if len(arguments) < 1 {
			return fmt.Errorf("la commande i2c attend 1 argument (ligne 5)")
		}
		n, err := strconv.ParseUint(strings.TrimSpace(arguments[0]), 10, 64)
		if err != nil {
			return fmt.Errorf("entier invalide : %q", arguments[0])
		}
		if n >= g.N {
			return fmt.Errorf("%d est hors limites (il faut 0 <= n < N = %d)", n, g.N)
		}
		fmt.Printf("i2c(%d) = %q\n", n, i2c(g, n))

	default:
		return fmt.Errorf("commande inconnue : %q", commande)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage : go run main.go fichier_test [fichier_test ...]")
		os.Exit(1)
	}

	erreur := false
	for _, nomFichier := range os.Args[1:] {
		fmt.Println("=====", nomFichier, "=====")

		g, commande, arguments, err := lireFichierTest(nomFichier)
		if err != nil {
			fmt.Fprintln(os.Stderr, "erreur :", err)
			erreur = true
			continue
		}

		afficherConfig(g)
		fmt.Println("Commande            :", commande)
		fmt.Println("-----")

		if err := executerCommande(g, commande, arguments); err != nil {
			fmt.Fprintln(os.Stderr, "erreur :", err)
			erreur = true
		}
		fmt.Println()
	}

	if erreur {
		os.Exit(1)
	}
}

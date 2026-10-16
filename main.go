package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
	"strconv"
	"strings"
	"encoding/binary"
)

// Global regroupe la configuration lue dans les 3 premières lignes d'un fichier test.
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

// sha1Hex calcule l'empreinte SHA-1 d'un texte (20 octets).
func sha1Hex(texte string) []byte {
	octetsDuTexte := []byte(texte)
	empreinte := sha1.Sum(octetsDuTexte) 
	return empreinte[:]          
}

// emprinteToString affiche une empreinte en hexadécimal, en majuscules.
func emprinteToString(empreinte []byte) string {
	enHexa := hex.EncodeToString(empreinte) // "9f57..."
	return strings.ToUpper(enHexa)          // "9F57..."
}

// ---------- Question 2 : configuration ----------

func afficherConfig(g Global) {
	fmt.Println("Fonction de hachage :", g.hashName)
	fmt.Println("Alphabet            :", g.alphabet)
	fmt.Println("Taille des clairs   :", g.taille)
	fmt.Println("Nombre de clairs N  :", g.N)
}

// ---------- Question 3 : indice -> texte clair ----------

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
	valeur := binary.LittleEndian.Uint64(y[:8])
	return (valeur + t) % g.N
}

// ---------- Question 5 : chaînes d'indices ----------
var TEXTE []byte
var EMPREINTE [sha1.Size]byte

func i2i(g Global, indice uint64, t uint64) uint64 {
	i2c(g, indice, TEXTE)
	EMPREINTE = sha1.Sum(TEXTE)
	return h2i(g, EMPREINTE[:], t)
}

type Chaine struct {
	debut uint64
	fin   uint64
}

// nouvelleChaine calcule une chaîne de "largeur" indices en partant de idx1 :
// idx1 --i2i(_, 1)--> i2 --i2i(_, 2)--> i3 ... --i2i(_, largeur-1)--> fin
func nouvelleChaine(g Global, idx1 uint64, largeur int) Chaine {
	indice := idx1
	for t := 1; t < largeur; t++ {
		indice = i2i(g, indice, uint64(t))
	}

	return Chaine{debut: idx1, fin: indice}
}

// ---------- Question 6 : création de la table ----------

// indexAleatoire renvoie un indice au hasard entre 0 et N-1.
func indexAleatoire(g Global) uint64 {
	return rand.Uint64N(g.N)
}

type Table struct {
	hauteur int // nombre de chaînes
	largeur int // longueur de chaque chaîne
	chaines []Chaine
}

// creerTable calcule "hauteur" chaînes de "largeur" indices, triées par fin croissante.
func creerTable(g Global, largeur int, hauteur int) Table {
	chaines := make([]Chaine, hauteur)

	for i := 0; i < hauteur; i++ {
		chaines[i] = nouvelleChaine(g, indexAleatoire(g), largeur)
	}

	sort.Slice(chaines, func(a int, b int) bool {
		return chaines[a].fin < chaines[b].fin
	})

	return Table{hauteur: hauteur, largeur: largeur, chaines: chaines}
}

// ---------- Question 7 : sauvegarde / ouverture / affichage ----------

func sauveTable(g Global, table Table, nomFichier string) {
	fichier, err := os.Create(nomFichier)
	if err != nil {
		fmt.Println("Impossible de créer le fichier :", nomFichier)
		os.Exit(1)
	}
	defer fichier.Close()

	ecrivain := bufio.NewWriter(fichier)

	fmt.Fprintln(ecrivain, g.hashName)
	fmt.Fprintln(ecrivain, g.alphabet)
	fmt.Fprintln(ecrivain, g.taille)
	fmt.Fprintln(ecrivain, table.hauteur)
	fmt.Fprintln(ecrivain, table.largeur)

	for _, chaine := range table.chaines {
		fmt.Fprintln(ecrivain, chaine.debut, chaine.fin)
	}

	ecrivain.Flush()
}

func ouvreTable(nomFichier string) (Global, Table) {
	lignes := lireLignes(nomFichier)

	hashName := lignes[0]
	alphabet := lignes[1]
	taille, _ := strconv.Atoi(lignes[2])
	hauteur, _ := strconv.Atoi(lignes[3])
	largeur, _ := strconv.Atoi(lignes[4])

	chaines := make([]Chaine, hauteur)
	for i := 0; i < hauteur; i++ {
		morceaux := strings.Fields(lignes[5+i])
		debut, _ := strconv.ParseUint(morceaux[0], 10, 64)
		fin, _ := strconv.ParseUint(morceaux[1], 10, 64)
		chaines[i] = Chaine{debut: debut, fin: fin}
	}

	g := newGlobal(hashName, alphabet, taille)
	table := Table{hauteur: hauteur, largeur: largeur, chaines: chaines}
	return g, table
}

func afficheTable(g Global, table Table, nbAffichees int) {
	afficherConfig(g)
	fmt.Println("Hauteur             :", table.hauteur)
	fmt.Println("Largeur             :", table.largeur)
	fmt.Println("Contenu :")

	for i := 0; i < table.hauteur; i++ {
		estAuDebut := i < nbAffichees
		estALaFin := i >= table.hauteur-nbAffichees

		if estAuDebut || estALaFin {
			chaine := table.chaines[i]
			fmt.Printf("  %09d: %d --> %d\n", i, chaine.debut, chaine.fin)
		} else if i == nbAffichees {
			fmt.Println("  ...")
		}
	}
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
		i2c(g, n, TEXTE)
		fmt.Println("i2c("+arguments[0]+") =", string(TEXTE))

	case "h2i":
		texte := arguments[0]
		t, _ := strconv.ParseUint(arguments[1], 10, 64)
		empreinte := sha1Hex(texte)
		indice := h2i(g, empreinte, t)
		fmt.Println("hash("+texte+") =", emprinteToString(empreinte))
		fmt.Println("h2i(hash("+texte+"), "+arguments[1]+") =", indice)

	case "i2i":
		indice, _ := strconv.ParseUint(arguments[0], 10, 64)
		t, _ := strconv.ParseUint(arguments[1], 10, 64)
		resultat := i2i(g, indice, t)
		// après l'appel, TEXTE et EMPREINTE contiennent les étapes intermédiaires
		fmt.Printf("%d --i2c--> %s --h--> %s --h2i(%d)--> %d\n",
			indice, string(TEXTE), emprinteToString(EMPREINTE[:]), t, resultat)
		fmt.Printf("%d  --i2i(%d)-->  %d\n", indice, t, resultat)

	case "chain":
		largeur, _ := strconv.Atoi(arguments[0])
		idx1, _ := strconv.ParseUint(arguments[1], 10, 64)
		chaine := nouvelleChaine(g, idx1, largeur)
		fmt.Printf("chain of length %d: %d ... %d\n", largeur, chaine.debut, chaine.fin)

	case "create":
		hauteur, _ := strconv.Atoi(arguments[0])
		largeur, _ := strconv.Atoi(arguments[1])
		nomFichier := arguments[2]
		table := creerTable(g, largeur, hauteur)
		sauveTable(g, table, nomFichier)
		fmt.Println("Table sauvegardée dans", nomFichier)

	case "info":
		nomFichier := arguments[0]
		gTable, table := ouvreTable(nomFichier)
		fmt.Println("Table lue dans", nomFichier)
		afficheTable(gTable, table, 10)

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
		TEXTE = make([]byte, g.taille)
		afficherConfig(g)
		fmt.Println("Commande            :", commande)
		fmt.Println("-----")

		executerCommande(g, commande, arguments)
		fmt.Println()
	}
}

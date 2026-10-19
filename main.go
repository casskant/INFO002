package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"os"
	"sort"
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

// ---------- hachage ----------

func sha1Hex(texte string) []byte {
	empreinte := sha1.Sum([]byte(texte))
	return empreinte[:]
}

func emprinteToString(empreinte []byte) string {
	return strings.ToUpper(hex.EncodeToString(empreinte))
}

// ---------- configuration ----------

func afficherConfig(g Global) {
	fmt.Println("Fonction de hachage :", g.hashName)
	fmt.Println("Alphabet            :", g.alphabet)
	fmt.Println("Taille des clairs   :", g.taille)
	fmt.Println("Nombre de clairs N  :", g.N)
}

// ---------- indice -> texte clair ----------

func i2c(g Global, n uint64, texte []byte) {
	nbLettres := uint64(len(g.alphabet))

	for position := g.taille - 1; position >= 0; position-- {
		texte[position] = g.alphabet[n%nbLettres]
		n = n / nbLettres
	}
}

// ---------- empreinte -> indice ----------

func h2i(g Global, y []byte, t uint64) uint64 {
	valeur := binary.LittleEndian.Uint64(y[:8])
	return (valeur + t) % g.N
}

// ---------- chaînes d'indices ----------

// variables globales réutilisées par i2i, pour éviter d'allouer à chaque appel
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

func nouvelleChaine(g Global, idx1 uint64, largeur int) Chaine {
	indice := idx1
	for t := 1; t < largeur; t++ {
		indice = i2i(g, indice, uint64(t))
	}

	return Chaine{debut: idx1, fin: indice}
}

// ---------- création de la table ----------

func indexAleatoire(g Global) uint64 {
	return rand.Uint64N(g.N)
}

type Table struct {
	hauteur int
	largeur int
	chaines []Chaine
}

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

// ---------- sauvegarde / ouverture / affichage ----------

// Format binaire : hashName, alphabet (longueur puis octets), taille, hauteur,
// largeur, puis (debut, fin) par chaîne. Chaque entier sur 8 octets little-endian.
func sauveTable(g Global, table Table, nomFichier string) {
	fichier, err := os.Create(nomFichier)
	if err != nil {
		fmt.Println("Impossible de créer le fichier :", nomFichier)
		os.Exit(1)
	}
	defer fichier.Close()

	ecrivain := bufio.NewWriter(fichier)

	ecrireTexte(ecrivain, g.hashName)
	ecrireTexte(ecrivain, g.alphabet)
	ecrireEntier(ecrivain, uint64(g.taille))
	ecrireEntier(ecrivain, uint64(table.hauteur))
	ecrireEntier(ecrivain, uint64(table.largeur))

	for _, chaine := range table.chaines {
		ecrireEntier(ecrivain, chaine.debut)
		ecrireEntier(ecrivain, chaine.fin)
	}

	ecrivain.Flush()
}

func ouvreTable(nomFichier string) (Global, Table) {
	fichier, err := os.Open(nomFichier)
	if err != nil {
		fmt.Println("Impossible d'ouvrir le fichier :", nomFichier)
		os.Exit(1)
	}
	defer fichier.Close()

	lecteur := bufio.NewReader(fichier)

	hashName := lireTexte(lecteur)
	alphabet := lireTexte(lecteur)
	taille := int(lireEntier(lecteur))
	hauteur := int(lireEntier(lecteur))
	largeur := int(lireEntier(lecteur))

	chaines := make([]Chaine, hauteur)
	for i := 0; i < hauteur; i++ {
		debut := lireEntier(lecteur)
		fin := lireEntier(lecteur)
		chaines[i] = Chaine{debut: debut, fin: fin}
	}

	g := newGlobal(hashName, alphabet, taille)
	table := Table{hauteur: hauteur, largeur: largeur, chaines: chaines}
	return g, table
}

func ecrireEntier(ecrivain *bufio.Writer, x uint64) {
	binary.Write(ecrivain, binary.LittleEndian, x)
}

func ecrireTexte(ecrivain *bufio.Writer, texte string) {
	ecrireEntier(ecrivain, uint64(len(texte)))
	ecrivain.WriteString(texte)
}

func lireEntier(lecteur *bufio.Reader) uint64 {
	var x uint64
	binary.Read(lecteur, binary.LittleEndian, &x)
	return x
}

func lireTexte(lecteur *bufio.Reader) string {
	longueur := lireEntier(lecteur)
	octets := make([]byte, longueur)
	io.ReadFull(lecteur, octets)
	return string(octets)
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

// ---------- recherche et inversion ----------

// recherche renvoie [a, b], les lignes dont la fin vaut idx (table triée).
// Renvoie (0, -1) si idx est absent.
func recherche(table Table, idx uint64) (int, int) {
	chaines := table.chaines

	a := sort.Search(len(chaines), func(i int) bool {
		return chaines[i].fin >= idx
	})

	if a == len(chaines) || chaines[a].fin != idx {
		return 0, -1
	}

	b := a
	for b+1 < len(chaines) && chaines[b+1].fin == idx {
		b++
	}
	return a, b
}

// verifieCandidat rejoue une chaîne de son début jusqu'à la colonne t-1 et
// vérifie que H du candidat vaut h. Le texte clair est écrit dans "clair".
func verifieCandidat(g Global, h []byte, t uint64, debut uint64, clair []byte) bool {
	indice := debut
	for col := uint64(1); col < t; col++ {
		indice = i2i(g, indice, col)
	}

	i2c(g, indice, clair)
	empreinte := sha1.Sum(clair)
	return bytes.Equal(empreinte[:], h)
}

func inverse(g Global, table Table, h []byte) (string, bool) {
	largeur := uint64(table.largeur)
	clair := make([]byte, g.taille)

	for t := largeur - 1; t >= 1; t-- {
		idx := h2i(g, h, t)
		for i := t + 1; i < largeur; i++ {
			idx = i2i(g, idx, i)
		}

		a, b := recherche(table, idx)
		for ligne := a; ligne <= b; ligne++ {
			if verifieCandidat(g, h, t, table.chaines[ligne].debut, clair) {
				return string(clair), true
			}
		}
	}

	return "", false
}

// ---------- recherche exhaustive ----------

func bruteforce(g Global, h []byte) (string, bool) {
	clair := make([]byte, g.taille)
	var cible [sha1.Size]byte
	copy(cible[:], h)

	for idx := uint64(0); idx < g.N; idx++ {
		i2c(g, idx, clair)
		if sha1.Sum(clair) == cible {
			return string(clair), true
		}
	}
	return "", false
}

// ---------- estimations (stats) ----------

// débit de i2i mesuré sur une machine de référence ; à réajuster selon la machine
const hashesParSeconde = 2211000.0

func estimationCouverture(N uint64, hauteur int, largeur int) float64 {
	Nf := float64(N)
	m := float64(hauteur)
	v := 1.0
	for i := 0; i < largeur; i++ {
		v = v * (1 - m/Nf)
		m = Nf * (1 - math.Exp(-m/Nf))
	}
	return 100 * (1 - v)
}

func tailleLisible(octets uint64) string {
	o := float64(octets)
	switch {
	case o >= 1024*1024:
		return fmt.Sprintf("%.1f M", o/1024/1024)
	case o >= 1024:
		return fmt.Sprintf("%.1f K", o/1024)
	default:
		return fmt.Sprintf("%d B", octets)
	}
}

func tempsLisible(secondes float64) string {
	s := int(secondes + 0.5)
	switch {
	case s >= 3600:
		return fmt.Sprintf("%d h %d min %d s", s/3600, (s%3600)/60, s%60)
	case s >= 60:
		return fmt.Sprintf("%d min %d s", s/60, s%60)
	default:
		return fmt.Sprintf("%d s", s)
	}
}

func afficheStats(g Global, hauteur int, largeur int) {
	octets := uint64(hauteur) * 16 // 2 entiers de 8 octets par chaîne

	// création : ~ hauteur*largeur appels à i2i ; recherche : ~ largeur² (pire cas)
	tCreation := float64(hauteur) * float64(largeur) / hashesParSeconde
	tRecherche := float64(largeur) * float64(largeur) / hashesParSeconde

	fmt.Println("Hauteur               :", hauteur)
	fmt.Println("Largeur               :", largeur)
	fmt.Println("Taille estimée        :", tailleLisible(octets))
	fmt.Println("Temps de création     :", tempsLisible(tCreation))
	fmt.Printf("Probabilité de succès : %.2f%%\n", estimationCouverture(g.N, hauteur, largeur))
	fmt.Println("Temps de recherche    :", tempsLisible(tRecherche))
}

// ---------- fichiers test ----------

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

func executerCommande(g Global, commande string, arguments []string) {
	switch commande {

	case "hash":
		texte := arguments[0]
		empreinte := sha1Hex(texte)
		fmt.Println("hash("+texte+") =", emprinteToString(empreinte))

	case "config":

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

	case "crack":
		empreinteHex := arguments[0]
		nomFichier := arguments[1]
		h, _ := hex.DecodeString(empreinteHex)

		gTable, table := ouvreTable(nomFichier) // config prise dans la table
		TEXTE = make([]byte, gTable.taille)

		clair, trouve := inverse(gTable, table, h)
		if trouve {
			fmt.Printf("crack(%s) = %q\n", emprinteToString(h), clair)
		} else {
			fmt.Printf("crack(%s) : aucun antécédent trouvé dans la table\n", emprinteToString(h))
		}

	case "bruteforce":
		empreinteHex := arguments[0]
		h, _ := hex.DecodeString(empreinteHex)
		clair, trouve := bruteforce(g, h)
		if trouve {
			fmt.Printf("bruteforce(%s) = %q\n", emprinteToString(h), clair)
		} else {
			fmt.Printf("bruteforce(%s) : aucun antécédent trouvé\n", emprinteToString(h))
		}

	case "stats":
		hauteur, _ := strconv.Atoi(arguments[0])
		largeur, _ := strconv.Atoi(arguments[1])
		afficheStats(g, hauteur, largeur)

	case "create_crack":
		hauteur, _ := strconv.Atoi(arguments[0])
		largeur, _ := strconv.Atoi(arguments[1])
		nomFichier := arguments[2]
		empreinteHex := arguments[3]

		table := creerTable(g, largeur, hauteur)
		sauveTable(g, table, nomFichier)
		fmt.Println("Table sauvegardée dans", nomFichier)

		h, _ := hex.DecodeString(empreinteHex)
		clair, trouve := inverse(g, table, h)
		if trouve {
			fmt.Printf("crack(%s) = %q\n", emprinteToString(h), clair)
		} else {
			fmt.Printf("crack(%s) : aucun antécédent trouvé dans la table\n", emprinteToString(h))
		}

	default:
		fmt.Println("Commande inconnue :", commande)
	}
}

func main() {
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

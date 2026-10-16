# INFO002 - Crytologie
## TP1 : compromis temps / mémoire, les tables arc-en-ciel
#### Par : Lilian Cusumano - Cassandre Kant 

___
## Utilisation de l'IA
Dans le cadre de ce TP, nous utilisons claude pour nous donner le plan d'action des certaines questions compliquées du TP. 
L'objectif est de lui donner le contexte de la fonction, la demande de résolution en le limitant sur les points d'attention.
Une fois le code ou une partie du code généré, nous relisons le code pour épurer ce qui doit l'être et corrigé en fonction de ce qu'on veut 
Le code est donc relu et tester pour chaque utilisation. 

Pour la partie compte rendu, on l'utilise pour reformuler une réponse lorsqu'on trouve qu'elle n'est pas assez compréhensible. Pour ce TP ce n'est pas le cas mais li nous arrive de générer les grands titres du compte rendu en fonction du sujet et de modifier en conséquence si ce n'est pas totalement ce qui est voulu. 


## 0. Préliminaire 

**Question 1** : 
Depuis la librairie ```crypto/sha1```, on a importé la fonction et nous avons pu tester la chaine de caractère salut, qui retourne le même résultat présenté dans le TP. 
On a donc une fonction qui ressmeble à : 
```
func sha1Hex(texte string) []byte {
	sum := sha1.Sum([]byte(texte))
	return sum[:]
}
```
## 0.1 Configuration globale

**Question 2** :
En go, nous avons crée la structure Global qui nous permet d'initialiser l'alphabet, la taille et le nombre de message en clair 
```
type Global struct {
	alphabet string
	taille   int
	N        uint64
}
```
## 0.2 Interface 

**Question 3** : 
- Das un premier temps, on ajoute la fonction de hachage dans la structure Global. 
- Pour créer un nouveau type, on a donc notre constructeur ```newGlobal```
- Pour la lecture des fichiers de test, nous avons dans la fonction ```lireFichierTest``` qui correspond à la fonction hash demandé dans le tp
- Pour l'affichage des données, nous avons cette fonction : 

```
func afficherConfig(g Global) {
	fmt.Println("Fonction de hachage :", g.hashName)
	fmt.Printf("Alphabet            : %s (%d caractères)\n", g.alphabet, len(g.alphabet))
	fmt.Println("Taille des clairs   :", g.taille)
	fmt.Println("Nombre de clairs N  :", g.N)
}
```

## 0.3 Indice et texte claire

**Question 4** : Nous avons cette fonction qui nous permet de transformer n en texte claire
 ```
func i2c(g Global, n uint64, texte []byte) {
nbLettres := uint64(len(g.alphabet))

	for position := g.taille - 1; position >= 0; position-- {
		numeroLettre := n % nbLettres
		texte[position] = g.alphabet[numeroLettre]
		n = n / nbLettres
	}
}
```
Et pour l'intégrer à la gestion des commandes, on a ajouté ce casx² dans la fonction ```executerCommande```
```
case "i2c":
		n, _ := strconv.ParseUint(arguments[0], 10, 64)
		texte := make([]byte, g.taille)
		i2c(g, n, texte)
		fmt.Println("i2c("+arguments[0]+") =", string(texte))
```
Nous avons donc un case par commande à executer (hash,config,i2c)

## 1. Idée du compromis temps-mémoire
**Question 5** (bonus) :
## 1.1 Principe et complexité : 

Le principe du compromis temps-mémoire repose sur le stockage de la première et de la dernière valeur de chaque chaîne, évitant ainsi d'enregistrer l'ensemble des maillons intermédiaires.
- Pour une taille totale d'espace de recherche $N$ et des chaînes de longueur $t$ :
  - Complexité mémoire : $\mathcal{O}(N/t)$ (seules les extrémités des chaînes sont conservées)
  - Complexité temporelle : $\mathcal{O}(t^2)$ (temps nécessaire pour rebalayer les chaînes lors de la recherche)
## 1.2. Comparaison avec les approches extrêmes
- Recherche exhaustive :
  - Temps : $\mathcal{O}(N)$
  - Mémoire : $\mathcal{O}(1)$ (ou $\mathcal{O}(\log N)$ si l'espace de recherche est trié)
- Précalcul complet (table de hachage) :
  - Temps : $\mathcal{O}(1)$
  - Mémoire : $\mathcal{O}(N)$

Cette méthode intermédiaire offre ainsi un équilibre optimal en réduisant l'empreinte mémoire par rapport au précalcul complet, tout en restant nettement plus rapide qu'une recherche exhaustive.
## 2. Mise en œuvre concrète : précalcul de la table
## 2.1. Indices et empreintes
**Question 6 :**
Voici notre implémentation de h2i : 
```
func h2i(g Global, y []byte, t uint64) uint64 {
	valeur := uint64(0)
	poids := uint64(1)
	for i := 0; i < 8; i++ {
		valeur = valeur + uint64(y[i])*poids
		poids = poids * 256
	}
	return (valeur + t) % g.N
}
```

De même que pour la question 4, nous avons ajouté un cas supplémentaire : 
```
case "h2i":
		texte := arguments[0]
		t, _ := strconv.ParseUint(arguments[1], 10, 64)
		empreinte := sha1Hex(texte)
		indice := h2i(g, empreinte, t)
		fmt.Println("hash("+texte+") =", emprinteToString(empreinte))
		fmt.Println("h2i(hash("+texte+"), "+arguments[1]+") =", indice)
```

## 2.2. Table arc-en-ciel
Question 7 : 
Voici l'implémentation de i2i : 
```
func i2i(g Global, indice uint64, t uint64) uint64 {
	i2c(g, indice, TEXTE)
	EMPREINTE = sha1.Sum(TEXTE)
	return h2i(g, EMPREINTE[:], t)
}
```
Pour récupérer la chaine nous avons : 
```
func nouvelleChaine(g Global, idx1 uint64, largeur int) Chaine {
	indice := idx1
	for t := 1; t < largeur; t++ {
		indice = i2i(g, indice, uint64(t))
	}

	return Chaine{debut: idx1, fin: indice}
}
```
Et comme les autres questions, on ajoute deux cas supplémentaires dans executer_fonction :
``` 
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
```
Question 8 : L'ajout du paramètre t dans la fonction h2i permet d'augmenter la couverture car elle évite les colisions entre les chaines. La couverture reste quand même inferieur à 100% mais elle déjà plus grande que si la variable t n'était pas dans la fonction.

Question 9 : Pour l'index aléatoire, on utilise la fonction random de la librairie math et pour créer table, on crée un tableau avec les index aléatoires à l'interieur. Ce qui ressemble à : 
```
// indexAleatoire renvoie un indice au hasard entre 0 et N-1.
func indexAleatoire(g Global) uint64 {
	return rand.Uint64N(g.N)
}

// creerTable calcule "hauteur" chaînes de "largeur" indices, triées par fin croissante.
func creerTable(g Global, largeur int, hauteur int) []Chaine {
	table := make([]Chaine, hauteur)

	for i := 0; i < hauteur; i++ {
		table[i] = nouvelleChaine(g, indexAleatoire(g), largeur)
	}

	sort.Slice(table, func(a int, b int) bool {
		return table[a].fin < table[b].fin
	})

	return table
}
```
Question 10 : Pour cette question, nous avons stocké nos tables dans des fichiers en bit, au début nous l'avions fait en stockant l'ascii mais on a trouvé plus pertinent d'enregistrer les bits. 
Donc nous avons sauver table : 

```
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
```

# 3. Mise en œuvre concrète : recherche dans une table arc-en-ciel
**Question 11** : Pour récupérer les premières lignes, nous avons implémenter recherche : 
```
func recherche(table Table, idx uint64) (int, int) {
	chaines := table.chaines

	// a = première ligne dont la fin est >= idx
	a := sort.Search(len(chaines), func(i int) bool {
		return chaines[i].fin >= idx
	})

	if a == len(chaines) || chaines[a].fin != idx {
		return 0, -1 // idx absent de la dernière colonne
	}

	// b = dernière ligne dont la fin vaut idx
	b := a
	for b+1 < len(chaines) && chaines[b+1].fin == idx {
		b++
	}
	return a, b
}

```
Et l'ajout du cas de la commande crack ressemble à : 
```
case "crack":
		empreinteHex := arguments[0]
		nomFichier := arguments[1]
		h, _ := hex.DecodeString(empreinteHex)

		// on utilise la configuration enregistrée dans la table elle-même
		gTable, table := ouvreTable(nomFichier)
		TEXTE = make([]byte, gTable.taille)

		clair, trouve := inverse(gTable, table, h)
		if trouve {
			fmt.Printf("crack(%s) = %q\n", emprinteToString(h), clair)
		} else {
			fmt.Printf("crack(%s) : aucun antécédent trouvé dans la table\n", emprinteToString(h))
		}

```

**Question 12** (bonus): Pour m chaînes, longueur t, réductions R₁…R_t (une par colonne), hash h.

Pour un haché h à retrouver, on suppose successivement qu’il se trouve à la colonne i = t−1, t−2, …, 1 :

On calcule la fin de chaîne à partir de cette hypothèse : R_i, h, R_{i+1}, h, …, R_{t−1}. Cela coûte (t−i) évaluations de h et de R.
On cherche le point obtenu parmi les fins de chaînes. Avec une table triée, c’est O(log m).
En cas de correspondance, on reconstruit la chaîne depuis son début (O(t)) et on vérifie que h y apparaît.

- Coût du calcul des fins de chaîne : O(t²) évaluations de h et R
- Coût des recherches : t × O(log m) 

Donc la complexité et égale à  O(t² + t·log m) = O(t²) (car log m ≪ t en pratique)

**Question 13** : Pour les estimations, nous aurons décidé de faire une fonction par estimation ce qui nous donne :
```
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
```
Et pour la gestion la commande stats : 
```
case "stats":
		hauteur, _ := strconv.Atoi(arguments[0])
		largeur, _ := strconv.Atoi(arguments[1])
		afficheStats(g, hauteur, largeur)
```

**Question 14** : 

**Question 15** : 

**Question 16** : 
# 4. Améliorations
Question 17 : 
# 4.1. Parallélisation du code
# 4.2. Plusieurs tables
# 4.3. Tailles min et max
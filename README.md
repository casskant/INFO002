# INFO002 - Crytologie
## TP1 : compromis temps / mémoire, les tables arc-en-ciel
#### Par : Lilian Cusumano - Cassandre Kant 

___
## Utilisation de l'IA
- Utilisation de Claude code, bon résualtat 
- Codé en GO

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
Question 10 : 

# 3. Mise en œuvre concrète : recherche dans une table arc-en-ciel
Question 11 : 

Question 12 (bonus):

Question 13 : 

Question 14 : 

Question 15 : 

Question 16 : 
# 4. Améliorations
Question 17 : 
# 4.1. Parallélisation du code
# 4.2. Plusieurs tables
# 4.3. Tailles min et max
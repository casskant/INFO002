# INFO002 - Crytologie
## TP1 : compromis temps / mémoire, les tables arc-en-ciel
#### Par : Lilian Cusumano - Cassandre Kant 

___
## Utilisation de l'IA
- Utilisation de Claude code, bon résualtat 
- Codé en GO

## Préliminaire 

**Question 1** : 
Depuis la librairie ```crypto/sha1```, on a importé la fonction et nous avons pu tester la chaine de caractère salut, qui retourne le même résultat présenté dans le TP. 
On a donc une fonction qui ressmeble à : 
```
func sha1Hex(texte string) []byte {
	sum := sha1.Sum([]byte(texte))
	return sum[:]
}
```
1. Configuration globale

**Question 2** :
En go, nous avons crée la structure Global qui nous permet d'initialiser l'alphabet, la taille et le nombre de message en clair 
```
type Global struct {
	alphabet string
	taille   int
	N        uint64
}
```
2. Interface 

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

3. Indice et texte claire
Question 4 :


# 1. Idée du compromis temps-mémoire
Question 5 (bonus) :
1. Principe et complexité : 

Le principe du compromis temps-mémoire repose sur le stockage de la première et de la dernière valeur de chaque chaîne, évitant ainsi d'enregistrer l'ensemble des maillons intermédiaires.
- Pour une taille totale d'espace de recherche $N$ et des chaînes de longueur $t$ :
  - Complexité mémoire : $\mathcal{O}(N/t)$ (seules les extrémités des chaînes sont conservées)
  - Complexité temporelle : $\mathcal{O}(t^2)$ (temps nécessaire pour rebalayer les chaînes lors de la recherche)
2. Comparaison avec les approches extrêmes
- Recherche exhaustive :
  - Temps : $\mathcal{O}(N)$
  - Mémoire : $\mathcal{O}(1)$ (ou $\mathcal{O}(\log N)$ si l'espace de recherche est trié)
- Précalcul complet (table de hachage) :
  - Temps : $\mathcal{O}(1)$
  - Mémoire : $\mathcal{O}(N)$

Cette méthode intermédiaire offre ainsi un équilibre optimal en réduisant l'empreinte mémoire par rapport au précalcul complet, tout en restant nettement plus rapide qu'une recherche exhaustive.
# 2. Mise en œuvre concrète : précalcul de la table
## 2.1. Indices et empreintes
Question 6 : 

## 2.2. Table arc-en-ciel
Question 7 : 
Question 8 : 
Question 9 : 
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
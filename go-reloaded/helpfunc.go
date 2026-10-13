package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// reCommande reconnaît (hex), (bin), (up), (low), (cap), avec ou sans nombre :
// "(up)", "(up, 2)", "( cap ,3 )"...
var reCommande = regexp.MustCompile(`\(\s*(hex|bin|up|low|cap)\s*(?:,\s*(\d+)\s*)?\)`)

func hexVersDecimal(mot string) string {
	valeur, err := strconv.ParseInt(mot, 16, 64) // lit le mot en base 16
	if err != nil {
		return mot // pas un hexa valide : on ne touche à rien
	}
	return strconv.FormatInt(valeur, 10) // on renvoie le nombre en base 10
}

func binVersDecimal(mot string) string {
	valeur, err := strconv.ParseInt(mot, 2, 64) // lit le mot en base 2
	if err != nil {
		return mot // pas un binaire valide : on ne touche à rien
	}
	return strconv.FormatInt(valeur, 10) // on renvoie le nombre en base 10
}

// changerCasse modifie les n derniers mots de la liste selon le mode
// ("up", "low" ou "cap").
func changerCasse(mots []string, n int, mode string) {
	debut := len(mots) - n
	if debut < 0 {
		debut = 0 // on ne remonte pas avant le premier mot
	}
	for i := debut; i < len(mots); i++ {
		switch mode {
		case "up":
			mots[i] = strings.ToUpper(mots[i])
		case "low":
			mots[i] = strings.ToLower(mots[i])
		case "cap":
			lettres := []rune(strings.ToLower(mots[i]))
			if len(lettres) > 0 {
				lettres[0] = unicode.ToUpper(lettres[0])
			}
			mots[i] = string(lettres)
		}
	}
}

// separerCommandes met une espace autour de chaque commande et la normalise
// ("(up, 2)" devient " (up,2) "), pour que le découpage en mots la garde
// d'un seul bloc, même collée à une ponctuation : "word(up),".
func separerCommandes(texte string) string {
	return reCommande.ReplaceAllStringFunc(texte, func(cmd string) string {
		m := reCommande.FindStringSubmatch(cmd)
		if m[2] == "" {
			return " (" + m[1] + ") "
		}
		return " (" + m[1] + "," + m[2] + ") "
	})
}

func traiterCommandes(mots []string) []string {
	var res []string
	for _, mot := range mots {
		// pas une commande : on le range et on passe au suivant
		if !strings.HasPrefix(mot, "(") || !strings.HasSuffix(mot, ")") {
			res = append(res, mot)
			continue
		}

		contenu := strings.Trim(mot, "()") // "up,2" ou "hex"
		nom := contenu
		n := 1
		if strings.Contains(contenu, ",") {
			parties := strings.Split(contenu, ",")
			nom = parties[0]
			if nombre, err := strconv.Atoi(parties[1]); err == nil {
				n = nombre
			}
		}

		switch nom {
		case "hex":
			if len(res) > 0 {
				res[len(res)-1] = hexVersDecimal(res[len(res)-1])
			}
		case "bin":
			if len(res) > 0 {
				res[len(res)-1] = binVersDecimal(res[len(res)-1])
			}
		case "up", "low", "cap":
			changerCasse(res, n, nom)
		default:
			res = append(res, mot) // parenthèses qui ne sont pas une commande
		}
	}
	return res
}

func formaterPonctuation(texte string) string {
	ponctuation := ".,!?:;"
	var res []rune
	for _, c := range texte {
		if strings.ContainsRune(ponctuation, c) {
			// on retire les espaces avant la ponctuation
			for len(res) > 0 && res[len(res)-1] == ' ' {
				res = res[:len(res)-1]
			}
			res = append(res, c)
		} else if c == ' ' {
			if len(res) > 0 && res[len(res)-1] != ' ' {
				res = append(res, c)
			}
		} else {
			// juste après une ponctuation : on remet une espace
			if len(res) > 0 && strings.ContainsRune(ponctuation, res[len(res)-1]) {
				res = append(res, ' ')
			}
			res = append(res, c)
		}
	}
	return string(res)
}

// estApostropheDeMot indique si l'apostrophe en position i est au milieu d'un
// mot (don't, it's...) : ce n'est alors pas un guillemet simple.
func estApostropheDeMot(lettres []rune, i int) bool {
	return i > 0 && i+1 < len(lettres) &&
		unicode.IsLetter(lettres[i-1]) && unicode.IsLetter(lettres[i+1])
}

func formaterApostrophes(texte string) string {
	lettres := []rune(texte)
	var res []rune
	ouvert := false
	for i := 0; i < len(lettres); i++ {
		c := lettres[i]
		if c != '\'' || estApostropheDeMot(lettres, i) {
			res = append(res, c)
			continue
		}
		if !ouvert {
			res = append(res, c)
			// on saute les espaces après l'apostrophe ouvrante
			for i+1 < len(lettres) && lettres[i+1] == ' ' {
				i++
			}
			ouvert = true
		} else {
			// on retire les espaces avant l'apostrophe fermante
			for len(res) > 0 && res[len(res)-1] == ' ' {
				res = res[:len(res)-1]
			}
			res = append(res, c)
			ouvert = false
		}
	}
	return string(res)
}

func corrigerArticles(mots []string) []string {
	for i := 0; i < len(mots)-1; i++ {
		if mots[i] != "a" && mots[i] != "A" {
			continue
		}
		premiere := unicode.ToLower([]rune(mots[i+1])[0]) // 1re lettre du mot suivant
		if strings.ContainsRune("aeiouh", premiere) {
			mots[i] += "n" // "a" -> "an", "A" -> "An"
		}
	}
	return mots
}

// traiterLigne applique toutes les transformations à une ligne de texte.
func traiterLigne(ligne string) string {
	mots := strings.Fields(separerCommandes(ligne))
	mots = traiterCommandes(mots)
	mots = corrigerArticles(mots)

	texte := strings.Join(mots, " ")
	texte = formaterPonctuation(texte)
	texte = formaterApostrophes(texte)
	return strings.TrimSpace(texte)
}

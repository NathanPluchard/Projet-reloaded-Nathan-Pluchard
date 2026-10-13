package main

import "testing"

func TestTraiterLigne(t *testing.T) {
	cas := []struct {
		nom, entree, attendu string
	}{
		{"hex", "1E (hex) files were added", "30 files were added"},
		{"bin", "It has been 10 (bin) years", "It has been 2 years"},
		{"up", "Ready, set, go (up) !", "Ready, set, GO!"},
		{"low", "I should stop SHOUTING (low)", "I should stop shouting"},
		{"cap", "Welcome to the Brooklyn bridge (cap)", "Welcome to the Brooklyn Bridge"},
		{"up avec nombre", "This is so exciting (up, 2)", "This is SO EXCITING"},
		{"nombre plus grand que le texte", "hello world (up, 10)", "HELLO WORLD"},
		{"commande collée à la ponctuation", "wow(up), nice", "WOW, nice"},
		{"ponctuation", "I was sitting over there ,and then BAMM !!", "I was sitting over there, and then BAMM!!"},
		{"groupe de ponctuation", "I was thinking ... You were right", "I was thinking... You were right"},
		{"point d'interrogation", "what do you think ?", "what do you think?"},
		{"apostrophes", "I am exactly how they describe me: ' awesome '", "I am exactly how they describe me: 'awesome'"},
		{"apostrophes plusieurs mots", "As Elton John said: ' I am the most well-known homosexual in the world '", "As Elton John said: 'I am the most well-known homosexual in the world'"},
		{"contraction", "I don't know", "I don't know"},
		{"article a -> an", "There it was. A amazing rock!", "There it was. An amazing rock!"},
		{"article devant h", "a hour ago", "an hour ago"},
		{"article inchangé", "a rock", "a rock"},
		{"hex invalide", "zz (hex)", "zz"},
		{"commande en début de ligne", "(up) hello", "hello"},
		{"ligne vide", "", ""},
		{"exemple du sujet 1", "it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6)", "It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness"},
		{"exemple du sujet 2", "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.", "Simply add 66 and 2 and you will see the result is 68."},
		{"exemple du sujet 3", "There is no greater agony than bearing a untold story inside you.", "There is no greater agony than bearing an untold story inside you."},
		{"exemple du sujet 4", "Punctuation tests are ... kinda boring ,what do you think ?", "Punctuation tests are... kinda boring, what do you think?"},
	}
	for _, c := range cas {
		t.Run(c.nom, func(t *testing.T) {
			if got := traiterLigne(c.entree); got != c.attendu {
				t.Errorf("traiterLigne(%q)\n  obtenu  : %q\n  attendu : %q", c.entree, got, c.attendu)
			}
		})
	}
}

func TestHexBin(t *testing.T) {
	if got := hexVersDecimal("1E"); got != "30" {
		t.Errorf("hexVersDecimal(1E) = %q, attendu 30", got)
	}
	if got := binVersDecimal("1010"); got != "10" {
		t.Errorf("binVersDecimal(1010) = %q, attendu 10", got)
	}
	if got := hexVersDecimal("xyz"); got != "xyz" {
		t.Errorf("un mot invalide doit rester inchangé, obtenu %q", got)
	}
}

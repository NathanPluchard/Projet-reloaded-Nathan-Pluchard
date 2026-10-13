package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage : go run . entree.txt sortie.txt")
		os.Exit(1)
	}

	contenu, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture :", err)
		os.Exit(1)
	}

	lignes := strings.Split(string(contenu), "\n")
	for i, ligne := range lignes {
		lignes[i] = traiterLigne(ligne)
	}

	err = os.WriteFile(os.Args[2], []byte(strings.Join(lignes, "\n")), 0o644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur d'écriture :", err)
		os.Exit(1)
	}
}

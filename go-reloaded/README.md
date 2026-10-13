# Go-reloaded

Petit outil en Go qui lit un fichier texte, applique des corrections et des
éditions automatiques, puis écrit le résultat dans un second fichier.

## Utilisation

```sh
go run . sample.txt result.txt
cat result.txt
```

## Transformations

| Règle | Exemple |
|-------|---------|
| `(hex)` : mot précédent (hexadécimal) → décimal | `1E (hex) files` → `30 files` |
| `(bin)` : mot précédent (binaire) → décimal | `10 (bin) years` → `2 years` |
| `(up)` : majuscules | `go (up)` → `GO` |
| `(low)` : minuscules | `SHOUTING (low)` → `shouting` |
| `(cap)` : initiale en majuscule | `bridge (cap)` → `Bridge` |
| `(up, n)`, `(low, n)`, `(cap, n)` : agit sur les `n` mots précédents | `so exciting (up, 2)` → `SO EXCITING` |
| Ponctuation `. , ! ? : ;` collée au mot précédent, espace après | `there ,and` → `there, and` |
| Groupes `...`, `!?` gardés ensemble | `thinking ... You` → `thinking... You` |
| Apostrophes `' '` collées au texte qu'elles entourent | `' awesome '` → `'awesome'` |
| `a` devient `an` devant une voyelle ou un `h` | `a untold story` → `an untold story` |

## Organisation

- `main.go` : lecture du fichier d'entrée, écriture du fichier de sortie.
- `helpfunc.go` : les fonctions de transformation (`traiterLigne` les enchaîne).
- `helpfunc_test.go` : tests unitaires (exemples du sujet + cas limites).

## Tests

```sh
go test ./...
```

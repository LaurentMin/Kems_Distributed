# Langage

Go

## Concept

Variante du jeu de carte Kem's où la partie se joue de manière individuelle, les joueurs doivent rassembler 4 cartes identiques puis annoncer Kem's avant qu'un autre joueur les contre. Pour rassembler ces cartes, les joueurs ont chacun 4 cartes en main et partagent une pioche de `nombre de joueurs` cartes. Si aucun joueur ne souhaite échanger (ou au bout d'un certain temps, les cartes sont renouvelées).
La partie s'arrête si un joueur gagne. Chaque joueur peut contrer en se trompant 1 fois.
Le jeu est paramétrable (nombre de joueurs, de cartes en mains, de cartes dans la pioche, de "faux" contre Kem's).

## Donnée partagée entre les sites

- La pioche
- Chaque site a un réplicat de celle-ci (copie locale)

## Cohérence des réplicats

- La pioche ne peut pas être modifiée que par un site à la fois (2 joueurs ne peuvent pas prendre la même carte)
- La modification est ensuite propagée aux autres sites
- Exclusion mutuelle [algorithme de la file d'attente répartie](https://moodle.utc.fr/pluginfile.php/172574/mod_resource/content/1/5-POLY-file-attente-2018.pdf) avec une *section critique* pour modifier la donnée
- Besoin d'implémenter les *estampilles* pour cet algo
- *Si donnée est diffusée après modification, pas besoin d'exclusion mutuelle pour la lecture (uniquement pour l'écriture)* ?

## Fonctionnalité de sauvegarde répartie datée

- Algorithme de calcul d'instantanés
- Horloges vectorielles pour dater les sauvegardes

## Application répartie est clairement structurée

- Architecture qui distingue les fonctionnalités applicatives des fonctionnalités de contrôle
- Avoir un réseau de test convaincant (voir **Etape 3**)
# Consignes du projet

## Objectifs pédagogiques

[moodle](https://moodle.utc.fr/mod/page/view.php?id=155894)

- Appliquer les concepts et les algorithmes de la première partie du cours.
- Programmer une application répartie et comprendre les difficultés à surmonter.
- Travailler en groupe, s'organiser, savoir gérer d'éventuels problèmes.
- Progresser en autonomie, rendre compte de son travail lors des TD dédiés et en fin de période par une soutenance orale en groupe.

## Cahier des charges

[moodle](https://moodle.utc.fr/mod/page/view.php?id=155900)

- L'application répartie utilise une donnée partagée entre les sites
  - Définir un scénario qui nécessite le partage d'au moins une donnée entre plusieurs "sites" : les instances de l'application réparties s'exécutant sur chaque site travaillent sur des _réplicats_ qui sont des copies locales de la donnée partagée.
- Les réplicats restent cohérents
  - N'autoriser qu'une seule modification de réplicat à la fois et propager les modifications aux autres réplicats.
  - Implémenter pour cela l'algorithme de la file d'attente répartie qui organise une exclusion mutuelle. La section critique correspond à l'accès exclusif à la donnée. À vous de voir s'il faut une exclusion mutuelle pour l'écriture et la lecture de la donnée partagée. À vous de voir comment adapter l'algorithme pour diffuser la mise à jour de la donnée partagée.
  - Cet algorithme utilise lui-même les estampilles, qu'il est donc nécessaire d'implémenter.
- L'application répartie inclut une fonctionnalité de sauvegarde répartie datée
  - Implémenter pour cela un algorithme de calcul d'instantanés du cours.
  - Pour dater la sauvegarde, utiliser des horloges vectorielles.
- L'application répartie est clairement structurée
  - Utiliser une architecture qui distingue les fonctionnalités applicatives des fonctionnalités de contrôle.
  - Définir au moins un réseau convaincant pour les tests.

## Etapes

### Etape 1 : algorithme de contrôle, messages structurés

[moodle](https://moodle.utc.fr/mod/page/view.php?id=167707)
Fait en TD (principe de base)

### Etape 2 : affichages dans le terminal

[moodle](https://moodle.utc.fr/mod/page/view.php?id=167772)
Fait en TD (fonctions utilitaires)

### Etape 3 : construire un réseau avec le shell

[moodle](https://moodle.utc.fr/mod/page/view.php?id=155924)
A faire (projet de base sans interface graphique en utilisant les pipes pour construire un réseau)

### Etape 4. Ajouter une interface graphique avec un client web

[moodle](https://moodle.utc.fr/mod/page/view.php?id=156211)
A faire (ajout d'une interface graphique séparée de l'application grâce à un serveur web, des clients web et des WebSocket en go)

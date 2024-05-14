# SR05 PROJET KEMS

Geffrelot Elouan, Huichalaf Kilapang, Minatchy Laurent

## Démarrage rapide (interface terminal)

Plus d'infos sur l'[interface terminal](./docs/terminal_interface.md)

### Mise en place

1. Se rendre dans la racine du projet `cd projet_sr05/`
2. Ajouter les droits d'exécution à tous les scripts `chmod +x *.sh`
3. Build les exécutables `./buildScript.sh`

### Débuter une exécution

4. Exécuter le clear script au moins une fois entre chaque exécution `./clearScript.sh`
5. Lancer le script d'initialisation `./reseauComplet.sh`
6. Dans un nouveau terminal lancer l'affichage `./displayScript.sh` (permet de suivre l'état du jeu)
7. Dans un nouveau terminal lancer le joueur 1 `./player1.sh` (permet de contrôler le joueur 1)
8. Dans un nouveau terminal lancer le joueur 2 `./player2.sh` (permet de contrôler le joueur 2)
9. Dans un nouveau terminal lancer le joueur 3 `./player3.sh` (permet de contrôler le joueur 3)
10. Il est maintenant possible de jouer :
    - En suivant les instructions qui s'affichent sur le terminal d'affichage
    - En entrant les commandes dans les terminaux de contrôle respectifs des joueurs
11. Remarques :
    - Il faut connecter 3 joueurs avant de débuter une partie (car lorsqu'un joueur se connecte la partie se réinitialise)

### Mettre fin à une exécution

12. Exécuter le clear script `./clearScript.sh` (cela met fin à tous les processus)

## Démarrage rapide (interface graphique)

Plus d'infos sur l'[interface graphique](./docs/graphical_interface.md)

## Fonctionnement

### L'application

L'application est en charge de recevoir les entrées utilisateur puis de calculer l'état suivant de jeu. Elle demande au contrôleur l'accès exclusif à la donnée avant de la modifier. Ensuite elle diffuse sa donnée via son contrôleur. Enfin elle envoi sa donnée à l'affichage afin qu'il se mette aussi à jour.

Plus d'infos sur l'[application](./docs/app.md)

### Le contrôleur

Le contrôleur est en charge de demander les accès exclusifs aux autres contrôleurs puis de l'attribuer à son application. Elle transmet aussi les nouveaux états de son application aux autres contrôleurs afin qu'ils les envoient à leurs propres applications afin qu'elles se mettent à jour.

Plus d'infos sur l'[application](./docs/app.md)

### Le réseau

Le réseau est simulé avec des lectures/écritures dans des pipes nommées. Chaque instance lit dans un premier fichier (son entrée) et écrit dans un second (sa sortie). Pour faire un lien bi-directionnel entre deux sites, il suffit de dubliquer la sortie de l'un dans l'entrée de l'autre et inversement.

<br>

Notre réseau comporte 3 applications (une par membre du groupe). Chaque application est reliée à son contrôleur et les contrôleurs sont liés entre-eux.
Tous ces liens sont bi-directionnels.

<br>

L'interface terminal est composée de 4 instances. Une pour afficher l'état du jeu (qui prend en entrée le sortie de toutes les applications). Les 3 autres pour envoyer les actions des utilisateurs aux application (la sortie de chacune de ces instances est liée à l'entrée d'une des applications).
Les communications de ces 2 derniers programmes et les applications sont uni-directionnels.

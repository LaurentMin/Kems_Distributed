# L'affichage terminal

## Deux programmes nécessaire

Pour l'affichage terminal, nous avons fait face à la difficulté de gérer à la fois les transmission de l'état du jeu par l'app, et les inputs du joueur. Comme il n'est pas possible pour un programme d'avoir deux stdin nous avons opté pour la solution suivante.

# Terminal-display

Le programme terminal-display est connecté au réseau et reçoit les mise à jour de l'état du jeu par les apps. A la réception il affiche l'état du jeu avec les cartes communes et les mains des joueurs. Le programme recherche aussi à chaque réception si il y a un potentiel gagnant (un joueur possédant un kems) pour en informer tous les joueurs, le potentiel gagnant doit entrer kems avant d'être contré. De plus il garde en mémoire le dernier état ce qui permet à la réception d'un nouvel état de le comparer à l'ancien pour savoir si il y a eu un gagnant ou un perdant et d'afficher ainsi les scores.

# Terminal-input

Le programme terminal-input est lui connecté au flux standard d'input d'une app (correspondant à un player). Il permet au joueur d'entrer des commandes comme échanger des cartes, passer à un nouveau tour, etc. Chaque commande est très simpliste pour les écrire rapidement, le programme se charge de transformer ces commandes en messages encodés et de les transmettre a l'app qui se chargera d'effectuer les commandes et les transmettra aux autres sites.

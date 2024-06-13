# L'affichage graphique

## Le web proxy

Le web proxy est un programme en go à la fois connecté à la sortie et l'entrée d'une app (correspondant à un joueur). Son rôle est de transmettre les mises à jour de l'app à l'interface graphique et de transmettre les actions du joueur dans l'interface graphique à l'app. Pour communiquer avec l'app il utilise comme les autres programmes stdin et sdtout et pour communiquer avec le site web il utilise une websocket. La partie écoute de l'app se fait dans une goroutine pour la séparer du service web qui permet de créer la websocket.  Le web proxy a aussi pour rôle de comparer les états de jeu comme terminal-display et d'informer l'interface graphique de la présence d'un potentiel gagnant. Enfin les communications entre l'interface graphique et web proxy se font via l'envoi de message avec le format JSON ce qui permet de travailler plus facilement avec le site web qui utilise Javascript. 

## L'interface graphique

L'interface graphique consiste en une page web par trois fichiers, un html, un css et un javascript. On affiche la main du joueur et les cartes communes qui sont des boutons avec un comportement "radio", après avoir selectionné deux cartes un bouton swap cards apparait qui permet d'échanger les cartes. La pioche fait office de bouton pour changer de tour et redistribuer des cartes. Il y a aussi une modal avec les règles, le bouton 'reset game' et une autre modal qui apprait à la fin d'un round pour afficher un résumé des scores.

Le fonctionnement dynamique du site est assuré par le code en Javascript qui permet de modifier le contenu de l'html, les classes des éléments html. Avant la connection une interface permet de choisir parmi les trois joueurs, le choix effectué determine quel port va être utilisé pour ouvrir la websocket correspondant à celle créé par les scripts `web-proxy.sh` correspondant elle même à un joueur. 

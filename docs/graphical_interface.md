# L'affichage graphique

## Le web proxy

Le web proxy est un programme en go à la fois connecté à la sortie et l'entrée d'une app (correspondant à un joueur). Son rôle est de transmettre les mises à jour de l'app à l'interface graphique et de transmettre les actions du joueur dans l'interface graphique à l'app. Pour communiquer avec l'app il utilise comme les autres programmes stdin et sdtout et pour communiquer avec le site web il utilise une websocket. La partie écoute de l'app se fait dans une goroutine pour la séparer du service web qui permet de créer la websocket.  Le web proxy a aussi pour rôle de comparer les états de jeu comme terminal-display et d'informer l'interface graphique de la présence d'un potentiel gagnant. Enfin les communications entre l'interface graphique et web proxy se font via l'envoi de message avec le format JSON ce qui permet de travailler plus facilement avec le site web qui utilise Javascript. 

## L'interface graphique

L'interface graphique consiste en une page web par trois fichiers, un html, un css et un javascript. Nous avons décidé d'intégrer 

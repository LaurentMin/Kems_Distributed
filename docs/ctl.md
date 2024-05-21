# Le contrôleur

## Intérêt

Chaque application est liée à un contrôleur. Il permet de gérer l'interaction de l'application avec les autres instances. Cela permet d'avoir une disctinction des fonctionnalités applicatives des contrôle.

## Exclusion mutuelle

Les contrôleurs sont responsables d'autoriser ou non leur application à modifier la donnée partagée. Pour ce faire, l'algorithme de la file d'attente répartie a été implémenté.

<br>

Cet algorithme a été implémenté fidèlement à celui vu en cours avec une légère particularité. La manière de simuler un réseau fait qu'on ne peut pas envoyer de message à un contrôleur spécifique. Dans le cas d'un accusé de réception il est donc nécessaire d'indiquer le contrôleur pour lequel on répond.

<br>

Cette notion de filtrage est utilisée après toutes les lectures de messages (dans l'application, le contrôleur, ...) pour s'assurer que la simulation de réseau ne pénalise pas l'application. Ainsi, un message reçu par une instance qui ne lui est pas destiné sera ignoré.

## Sauvegarde

On utilise un algorithme de marqueur pour faire la sauvegarde sur le système.
Le contrôleur est chargé de la sauvegarde de l'état du jeu. Il peut recevoir une demande de la par de son application ou d'un autre controleur.
Dans le premier cas, cela veut dire que son application est iitiateur de la sauvegarde. Il fait alors la sauvegarde avec l'état du jeu envoyé par son application et envoie ensuite une demande aux autres controleur.
Dans le second cas, il demande l'état du jeu à son application pour faire la sauvegarde.

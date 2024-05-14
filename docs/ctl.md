# Le contrôleur

## Sauvegarde

On utilise un algorithme de marqueur pour faire la sauvegarde sur le système.
Le contrôleur est chargé de la sauvegarde de l'état du jeu. Il peut recevoir une demande de la par de son application ou d'un autre controleur.
Dans le premier cas, cela veut dire que son application est iitiateur de la sauvegarde. Il fait alors la sauvegarde avec l'état du jeu envoyé par son application et envoie ensuite une demande aux autres controleur.
Dans le second cas, il demande l'état du jeu à son application pour faire la sauvegarde.

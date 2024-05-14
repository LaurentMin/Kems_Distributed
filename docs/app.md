# L'application

## Concept

Nous nous sommes inspirés du célèbre jeu de cartes **Kem's** pour créer notre application. Nous proposons une variante du jeu de carte Kem's où la partie se joue de manière individuelle. Les joueurs doivent rassembler `N` cartes identiques puis annoncer Kem's avant qu'un autre joueur les contre. Lorsqu'ils y parviennent ils gagnent un point si ils se font contrer, ils en perdent un. Le gagnant est celui qui a le plus de points.

<br>

Pour rassembler ces cartes, les joueurs ont chacun `N` cartes en main et partagent une pioche de `M` cartes. Lorsqu'aucun joueur ne souhaite échanger de cartes, cela marque la fin du tour et les joueurs peuvent passer au tour suivant. Les cartes de la pioche sont alors renouvelées puis l'échange est à nouveau possible...

## Donnée partagée

La donnée que chaque site partage est l'état actuel du jeu. C'est à dire, les joueurs avec leurs scores et leurs mains; la pile dans laquelle les joueurs peuvent piocher; le paquet de carte qui sert à renouveler la pile ou redistribuer des cartes aux joueurs; la défausse.

<br>

Ainsi, peu importe la modification apportée à l'état du jeu, il suffit d'envoyer cette donnée (qui est transformée en chaine de caractères pour l'envoi). Les applications qui la reçoivent peuvent la transformer dans un format qui leur permet de se mettre à jour, tout comme les instances qui font de l'affichage.

## Cohérence des réplicats

Chaque application a donc un réplicat de la donnée partagée. La cohérence de cette donnée est assurée par le contrôleur. Avant chaque modification de la donnée, l'application attend d'avoir l'autorisation du contrôleur pour la modifier.

<br>

Dans notre système, nous avons décidé que dans le cas où l'application reçoit une action `a1`, puis fait une demande de modification; mais reçoit une action `a2` avant d'avoir pu appliquer `a1`, elle oublie `a1` et applique `a2` lorsqu'elle obtient le droit de modification. Comme en principe, l'application est plus rapide que l'utilisateur cette situation n'arrive jamais. Et dans le cas où le système est surchargée, il n'est pas une bonne idée de lui faire faire plus d'actions.

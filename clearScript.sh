#!/bin/bash

# Fonction de nettoyage
nettoyer () {
  # Suppression des processus de l'application affichage terminal
  killall terminal-display 2> /dev/null

  # Suppression des processus de l'application input terminal
  killall terminal-input 2> /dev/null

  # Suppression des processus de l'application app
  killall app 2> /dev/null

  # Suppression des processus de l'application ctl
  killall ctl 2> /dev/null
 
  # Suppression des processus NET
  killall net 2> /dev/null
  killall netadmin 2> /dev/null

  # Suppression des processus tee et cat
  killall tee 2> /dev/null
  killall cat 2> /dev/null
 
  # Suppression des tubes nommés
  \rm -f /tmp/in* /tmp/out*
  echo "Nettoyage terminé."
  exit 0
}
 
# Appel de la fonction nettoyer
echo "Nettoyage en cours..."
nettoyer
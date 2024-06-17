#!/bin/bash

if [ "$1" = "-p" ]; then
    regexNumber='^[0-9]+$'
    if ! [[ $2 =~ $regexNumber ]] ; then
        echo "error: Not a number" >&2; exit 1
    fi
    echo "Launching input for player $2 !!"
    ./web-proxy -p 444$2 -a localhost -n $2 > /tmp/in_A$2  < /tmp/in_D$2
else    
    echo "Invalid argument"
fi
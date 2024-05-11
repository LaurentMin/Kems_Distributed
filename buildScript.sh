#!/bin/bash

echo "Building..."
go build -o app app.go app-utils.go utils.go logs.go
echo "Built application"
go build -o ctl ctl.go utils.go logs.go
echo "Built controller"
go build -o player player.go utils.go logs.go
echo "Built player"
go build -o display display-game.go app-utils.go utils.go logs.go
echo "Built display"
echo "Done."
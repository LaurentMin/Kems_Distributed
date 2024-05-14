#!/bin/bash

echo "Building..."
go build -o app app.go app-utils.go utils.go logs.go io-utils.go
echo "Built application"
go build -o ctl ctl.go utils.go logs.go io-utils.go
echo "Built controller"
go build -o terminal-input terminal-input.go utils.go logs.go app-utils.go io-utils.go
echo "Built terminal-input"
go build -o terminal-display terminal-display.go app-utils.go utils.go logs.go io-utils.go
echo "Built display"
echo "Done."
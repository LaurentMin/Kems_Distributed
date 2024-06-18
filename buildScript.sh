#!/bin/bash

echo "Building..."
go build -o app app.go app-utils.go utils.go logs.go io-utils.go
echo "Built application"
go build -o ctl ctl.go utils.go logs.go io-utils.go
echo "Built controller"
go build -o terminal-input terminal-input.go utils.go logs.go app-utils.go client-utils.go io-utils.go
echo "Built terminal-input"
go build -o terminal-display terminal-display.go app-utils.go utils.go logs.go client-utils.go io-utils.go
echo "Built display"
go build -o web-proxy web-proxy.go app-utils.go utils.go logs.go client-utils.go io-utils.go
echo "Built web-proxy"

echo "Building NET..."
cd ./NET/
go build -o net net.go net-utils.go pastouche3.go pastouche2.go pastouche1.go
echo "Done."
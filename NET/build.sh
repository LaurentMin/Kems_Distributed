#!/bin/bash

echo "Building NET..."
go build -o net net.go net-utils.go pastouche3.go pastouche2.go pastouche1.go
echo "Done."

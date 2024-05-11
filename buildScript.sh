echo "Building..."
go build -o app app.go app-utils.go utils.go logs.go
echo "Built app"
go build -o ctl ctl.go utils.go logs.go
echo "Built ctl"
go build -o player player.go utils.go logs.go
echo "Built player"
echo "Done."
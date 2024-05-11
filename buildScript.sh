echo "Building..."
go build -o app app.go app-utils.go utils.go logs.go
echo "Built app"
go build -o ctl ctl.go utils.go logs.go
echo "Built ctl"
echo "Done."
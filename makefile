source := main.go

release:
	GOOS=linux GOARCH=amd64 go build -o miniflux-exporter-linux64 ${source}

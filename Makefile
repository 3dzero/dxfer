clean:
	rm bin/*

bin/dxfer.exe:
	GOOS=windows GOARCH=amd64 go build -o bin/dxfer.exe cmd/server/server.go


bin/dxfer-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/dxfer-linux cmd/server/server.go

bin/dxfer:
	GOARCH=amd64 go build -o bin/dxfer cmd/server/server.go
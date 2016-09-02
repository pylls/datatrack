install:
	go build -o bin/datatrack main.go 

test:
	go test ./... -cover

doc:
	godoc -http :6060

# clean:
# 	rm -rf pkg/

distclean: clean
	rm -rf ./bin/
	rm -rf ./datatrack.db

cross:
	env GOOS=windows GOARCH=amd64 go build -o bin/datatrack.win64.exe main.go
	env GOOS=darwin GOARCH=amd64 go build -o bin/datatrack.darwin64 main.go
	env GOOS=linux GOARCH=amd64 go build -o bin/datatrack.linux64 main.go

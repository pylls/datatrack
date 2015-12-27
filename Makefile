GOPATH=$(PWD)
export GOPATH
install:
	go build -o bin/datatrack datatrack/server

test:
	go test datatrack... -cover

doc:
	godoc -http :6060

clean:
	rm -rf pkg/

distclean: clean
	rm -rf bin/
	rm -rf datatrack.db

cross:
	env GOOS=windows GOARCH=amd64 go build -o bin/datatrack.win64.exe datatrack/server
	env GOOS=darwin GOARCH=amd64 go build -o bin/datatrack.darwin64 datatrack/server
	env GOOS=linux GOARCH=amd64 go build -o bin/datatrack.linux64 datatrack/server

env:
	@echo export GOPATH=$(GOPATH)

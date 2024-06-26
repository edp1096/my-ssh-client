build:
	go build -trimpath -o bin/ .

build-vendor:
	go work vendor
	go build -mod vendor -trimpath -o bin/ .

debug:
	go build -o bin/ .

dist:
	go get -d github.com/mitchellh/gox
	go build -mod=readonly -o ./bin/ github.com/mitchellh/gox
	go mod tidy
	go env -w GOFLAGS=-trimpath
	bin/gox -mod="readonly" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="windows/amd64 linux/amd64 linux/arm linux/arm64 darwin/amd64 darwin/arm64" .
	rm ./bin/gox*

test:
	go test ./... -race -cover -count=1

clean:
	rm -rf ./bin/*
	rm -f ./coverage.html
	rm -f ./coverage.out

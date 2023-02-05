.PHONY: build

build:
	rm -fr ./bin
	go build -o ./bin/toml2cm ./main.go

install: build
	cp ./bin/toml2cm ~/.local/bin

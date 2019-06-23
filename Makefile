PROJECT?=auto-docs
V:=$(shell git log | head -n 1 | cut -d\  -f2)
LDFLAGS:=-ldflags "-X main.GitCommit=${V}"

bin-prep:
	GO111MODULE=off go get -u github.com/kevinburke/go-bindata/...

bin-dist:
	go-bindata -o $(PROJECT)/server/assets.go -prefix dist/ dist/...
	sed -i "s/package main/package server/" $(PROJECT)/server/assets.go

binaries:
	$(MAKE) binary GOARCH=amd64 GOOS=linux
	$(MAKE) binary GOARCH=amd64 GOOS=windows
	$(MAKE) binary GOARCH=amd64 GOOS=darwin
	$(MAKE) binary GOARCH=386 GOOS=linux
	$(MAKE) binary GOARCH=386 GOOS=windows

binary: GOARCH?=amd64
binary: GOOS?=linux
binary:
	go build $(LDFLAGS) -o build/$(PROJECT).$(GOARCH)-$(GOOS) ./auto-docs
	if [ "$(GOOS)" = "windows" ]; then \
		mv build/$(PROJECT).$(GOARCH)-$(GOOS) build/$(PROJECT).$(GOARCH)-$(GOOS).exe; \
	fi

build-fe:
	yarn
	yarn build

clean:
	rm -f $(PROJECT)/server/assets.go
	rm -r dist

compile:
	go install $(LDFLAGS) ./$(PROJECT)/

coverage:
	go test -race -covermode=atomic -coverprofile=c.out ./...
	sed -i '/^github.com\/cloudcloud\/auto-docs\/auto-docs\/server\/assets.go.*/d' c.out
	go tool cover -html=c.out -o cover.html

# at this time, there's no watch enabled for the go binary
dev-be: bin-prep bin-dist install
	$(PROJECT) server

# serve is a watch task with built-in node server
dev-fe:
	yarn serve

image: VERSION?=1.0.0
image:
	docker build -t cloudcloud/auto-docs:v$(VERSION) .

install: build-fe compile

test: bin-dist install
	go test -v -race ./...

complete-binaries: clean bin-prep build-fe bin-dist binaries

#
# Docker Compose management commands.
#
DN?=auto-docs
DC?=docker-compose -p $(DN)

up:
	$(DC) up -d

shell:
	docker exec -it $(DN)_autodocs_1 bash

down:
	$(DC) down


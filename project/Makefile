VERSION=$(shell git describe --tags --always)

.PHONY: build
build:
	mkdir -p build/ && \
	go mod tidy && \
	go build -ldflags "-X main.Version=$(VERSION)" -o ./build/app ./cmd

.PHONY: gorm
gorm:
	go mod tidy
	cd ./tool ; go run gormgen.go

.PHONY: api
api:
	go mod tidy
	echo projectName: $(name)
	cd ./tool ; go run apigen.go --name $(name)

.PHONY: wire
wire:
	go install github.com/google/wire/cmd/wire@latest
	go mod tidy
	cd ./app ; wire

.PHONY: swag
swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	go mod tidy
	swag fmt; swag init -d ./handler,./po,./vo -g ../api/api.go --generatedTime   

.PHONY: trans 
trans:
	go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
	go mod tidy
	goi8n extract; goi18n merge active.en.toml translate.zh.toml




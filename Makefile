# Unset GOPATH so Go uses module mode when project is inside GOPATH (e.g. ~/Go on macOS).
GO := GOPATH= go

.PHONY: deps build run test tidy
deps tidy:
	$(GO) mod tidy
build:
	$(GO) build -o app .
run: build
	./app
test:
	$(GO) test ./...

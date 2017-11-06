.PHONY: all clean deps fmt vet test docker

EXECUTABLE ?= cryptotrader
IMAGE ?= euskadi31/$(EXECUTABLE)
VERSION ?= $(shell git describe --match 'v[0-9]*' --dirty='-dev' --always)
COMMIT ?= $(shell git rev-parse --short HEAD)

LDFLAGS = -X "main.Revision=$(COMMIT)" -X "main.Version=$(VERSION)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

release:
	@echo "Release v$(version)"
	@git pull
	@git checkout master
	@git pull
	@git checkout develop
	@git flow release start $(version)
	@echo "$(version)" > .version
	@sed -e "s/version: .*/version: \"v$(version)\"/g" docs/swagger.yaml > docs/swagger.yaml.new && rm -rf docs/swagger.yaml && mv docs/swagger.yaml.new docs/swagger.yaml
	@git add .version docs/swagger.yaml
	@git commit -m "feat(project): update version file" .version docs/swagger.yaml
	@git flow release finish $(version)
	@git push
	@git push --tags
	@git checkout master
	@git push
	@git checkout develop
	@echo "Release v$(version) finished."

all: deps build test

clean:
	@go clean -i ./...

deps:
	@dep ensure

fmt:
	@go fmt $(PACKAGES)

vet:
	@go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -ldflags '-s -w $(LDFLAGS)' -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

travis:
	@for PKG in $(PACKAGES); do go test -ldflags '-s -w $(LDFLAGS)' -cover -covermode=count -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

cover: test
	@echo ""
	@for PKG in $(PACKAGES); do go tool cover -func $$GOPATH/src/$$PKG/coverage.out; echo ""; done;

docker:
	@sudo docker build --no-cache=true --rm -t $(IMAGE) .

publish: docker
	@sudo docker tag $(IMAGE) $(IMAGE):latest
	@sudo docker push $(IMAGE)

bindata.go: docs/index.html docs/swagger.yaml
	@echo "Bin data..."
	@go-bindata docs/

$(EXECUTABLE): $(shell find . -type f -print | grep -v vendor | grep "\.go")
	@echo "Building $(EXECUTABLE)..."
	@CGO_ENABLED=0 go build -ldflags '-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

run: docker
	@sudo docker run -e "HYPERPIC_AUTH_SECRET=c8da8ded-f9a2-429c-8811-9b2a07de8ede" -p 8574:8080 -v $(shell pwd)/var/lib/hyperpic:/var/lib/hyperpic --rm $(IMAGE)

dev: $(EXECUTABLE)
	@./$(EXECUTABLE)

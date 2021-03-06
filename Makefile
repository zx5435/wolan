.PHONY: build

default:
	cat Makefile

build:
	./gradlew bootJar
	docker build -f __cicd__/Dockerfile -t zx5435/wolan .

up:
	docker stop wolan
	docker rm wolan
	docker run -it -d --name wolan -p8080:8080 \
	    -v "$$PWD/__task__":/www/__task__ \
	    zx5435/wolan

up222222222222222:
	cd __work__ && docker run -it -d --name wolan -p 4321:23456 \
	    -v "$$PWD":/app/__work__ \
	    -v "/var/run/docker.sock:/var/run/docker.sock" \
	    zx5435/wolan

build-go:
	docker run -it --rm \
	    -v "$$GOPATH/src":/go/src \
	    -w /go/src/github.com/zx5435/wolan/cmd/wolan-server \
	    golang:1.10.2 \
        sh -c "CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo -o wolan-server"

ingress-build:
	docker run -it --rm \
	    -v "$$GOPATH/src":/go/src \
	    -w /go/src/github.com/zx5435/wolan/cmd/wolan-ingress \
	    golang:1.10.2 \
        sh -c "CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo -o wolan-ingress"

ingress-test:
	docker run -it -d --name wolan-ingress -p80:80 -p443:443 zx5435/wolan:ingress

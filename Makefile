.PHONY: all
all:
	go generate ./pkg/version
	CGO_ENABLED=0 go build -o bin/ktranslate ./cmd/ktranslate

.PHONY: windows
windows:
	go generate ./cmd/version
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/ktranslate.exe ./cmd/ktranslate

.PHONY: arm
arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/ktranslate ./cmd/ktranslate

.PHONY: test
test:
	go test ./cmd/... ./pkg/...

.PHONY: bench
bench:
	go test -bench=. ./cmd/... ./pkg/...

.PHONY: ktranslate
ktranslate:
	go install ./cmd/ktranslate

.PHONY: clean
clean:
	rm -f bin/ktranslate

.PHONY: generate
generate:
	go generate ./...

.PHONY: install
install:
	mkdir -p $(DESTDIR)/usr/local/bin
	install -m 0755 bin/ktranslate $(DESTDIR)/usr/local/bin

.PHONY: docker
docker: all
	docker pull ubuntu:20.04
	docker build -t ktranslate:v2 -f Dockerfile .

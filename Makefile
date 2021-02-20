.PHONY: all
all:
	GOPRIVATE=github.com/kentik go build -tags dynamic -o bin/ktranslate ./cmd/ktranslate

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
	GOPRIVATE=github.com/kentik go generate ./...

.PHONY: install
install:
	mkdir -p $(DESTDIR)/usr/local/bin
	install -m 0755 bin/ktranslate $(DESTDIR)/usr/local/bin

.PHONY: docker
docker: all
	docker build -t ktranslate:v1 -f Dockerfile .

.PHONY: publish
publish: docker
	docker tag ktranslate:v1 gcr.io/kentik-continuous-delivery/ktranslate:v1
	docker push gcr.io/kentik-continuous-delivery/ktranslate:v1
	docker tag ktranslate:v1 kentik/ktranslate:v1
	docker push kentik/ktranslate:v1

.PHONY: pub_latest
pub_latest: publish
	docker tag ktranslate:v1 kentik/ktranslate:latest
	docker push kentik/ktranslate:latest

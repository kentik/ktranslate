.PHONY: all
all:
	go build -tags dynamic -o bin/ktranslate ./cmd/ktranslate

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

.PHONY: publish
publish: docker
	docker tag ktranslate:v2 gcr.io/kentik-continuous-delivery/ktranslate:v2
	docker push gcr.io/kentik-continuous-delivery/ktranslate:v2
	docker tag ktranslate:v2 kentik/ktranslate:v2
	docker tag ktranslate:v2 kentik/ktranslate:v1
	docker push kentik/ktranslate:v2
	docker push kentik/ktranslate:v1

.PHONY: pub_latest
pub_latest: publish
	docker tag ktranslate:v2 kentik/ktranslate:latest
	docker push kentik/ktranslate:latest

.PHONY: pub_aws
pub_aws: all
	docker pull public.ecr.aws/lambda/provided:al2
	docker build -t ktranslate:aws -f DockerfileAws .
	docker tag ktranslate:aws kentik/ktranslate:aws
	docker push kentik/ktranslate:aws

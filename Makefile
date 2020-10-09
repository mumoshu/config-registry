# find or download addlicense
addlicense:
ifeq (, $(shell which addlicense))
	@{ \
	set -e ;\
	INSTALL_TMP_DIR=$$(mktemp -d) ;\
	cd $$INSTALL_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/google/addlicense ;\
	rm -rf $$INSTALL_TMP_DIR ;\
	}
ADDLICENSE=$(GOBIN)/addlicense
else
ADDLICENSE=$(shell which addlicense)
endif

# find or download goreleaser
goreleaser:
ifeq (, $(shell which goreleaser))
	@{ \
	set -e ;\
	INSTALL_TMP_DIR=$$(mktemp -d) ;\
	cd $$INSTALL_TMP_DIR ;\
	go mod init tmp ;\
	go get github.com/goreleaser/goreleaser ;\
	rm -rf $$INSTALL_TMP_DIR ;\
	}
GORELEASER=$(GOBIN)/goreleaser
else
GORELEASER=$(shell which goreleaser)
endif

.PHONY: format
format:
	gofmt -w .

.PHONY: test/format
test/format:
	gofmt -s -d . > gofmt.out
	test -z "$$(cat gofmt.out)" || (cat gofmt.out && rm gofmt.out && false)

.PHONY: test/release
test/release: goreleaser
	$(GORELEASER) release --snapshot --skip-publish --rm-dist

.PHONY: test/krew
test/krew:
	docker run -v $(PWD)/.krew-release-bot.yaml:/tmp/template-file.yaml rajatjindal/krew-release-bot:v0.0.38 \
	krew-release-bot template --tag v0.1.0 --template-file /tmp/template-file.yaml

.PHONY: test/go
test/go:
	go test ./...

.PHONY: check
check: test/format test/release test/krew test/go

.PHONY: build
build:
	go build ./cmd/kubeconf

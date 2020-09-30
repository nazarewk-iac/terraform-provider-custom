GPG_FINGERPRINT=nazarewk@gmail.com
export GPG_FINGERPRINT

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: test-release
test-release:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release:
	GITHUB_TOKEN ?= $(shell bash -c 'read -s -p "GITHUB_TOKEN: " GITHUB_TOKEN ; echo "$GITHUB_TOKEN"')
	export GITHUB_TOKEN
	goreleaser
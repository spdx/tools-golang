.PHONY: test
test: unit fuzz

.PHONY: bootstrap
bootstrap:
	go install github.com/rinchsan/gosimports/cmd/gosimports@v0.3.8

.PHONY: format
format:
	gofmt -w .
	go mod tidy
	gosimports -w -local github.com/spdx .

.PHONY: unit
unit:
	go test -v -covermode=count -coverprofile=profile.cov.tmp ./...
	cat profile.cov.tmp | grep -v /model.go > profile.cov # ignore generated model file

.PHONY: fuzz
fuzz:
	go test -v -run=Fuzz -fuzz=FuzzShouldIgnore ./utils -fuzztime=10s
	go test -v -run=Fuzz -fuzz=FuzzPackageCanGetVerificationCode ./utils -fuzztime=10s

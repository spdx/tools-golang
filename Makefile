.PHONY: test
test: unit fuzz

.PHONY: format
format:
	gofmt -w .
	gosimports -w -local github.com/spdx .

.PHONY: unit
unit:
	go test -v -covermode=count -coverprofile=profile.cov ./...

.PHONY: fuzz
fuzz:
	go test -v -run=Fuzz -fuzz=FuzzShouldIgnore ./utils -fuzztime=10s
	go test -v -run=Fuzz -fuzz=FuzzPackageCanGetVerificationCode ./utils -fuzztime=10s

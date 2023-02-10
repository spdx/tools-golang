test :
	go test -v -cover ./...
fuzz :
	go test -v -cover -fuzz=FuzzShouldIgnore ./utils -fuzztime=10s
	go test -v -cover -fuzz=FuzzPackageCanGetVerificationCode ./utils -fuzztime=10s

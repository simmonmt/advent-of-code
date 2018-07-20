DEPS=$(shell find src -name \*.go \! -name \*_test.go)
PKGS=$(shell find src -type d -depth 1 |xargs -n1 basename)

otp: $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o $@ src/main.go

test:
	GOPATH=$$(PWD):$$GOPATH go test $(PKGS)

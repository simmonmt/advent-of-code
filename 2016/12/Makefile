DEPS=$(shell find src -name \*.go \! -name \*_test.go)

asm: $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o $@ src/main.go

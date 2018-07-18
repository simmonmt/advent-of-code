DEPS=$(shell find src -name \*.go \! -name \*_test.go)

decoder: $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o $@ src/main.go

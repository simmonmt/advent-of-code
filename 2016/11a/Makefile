DEPS=$(shell find src -name \*.go \! -name \*_test.go)

game: $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o $@ src/main.go

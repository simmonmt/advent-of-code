DEPS=$(shell find src -name \*.go \! -name \*_test.go)

sleigh: src/main.go $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o sleigh src/main.go

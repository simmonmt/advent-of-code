DEPS=$(shell find src -name \*.go \! -name \*_test.go)
PKGS=$(shell find src -type d -depth 1 |xargs -n1 basename)

solution: $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o $@ src/*.go

test:
	GOPATH=$$(PWD):$$GOPATH go test $(PKGS)

clean:
	$(RM) solution

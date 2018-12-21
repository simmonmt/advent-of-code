DEPS=$(shell find src -name \*.go \! -name \*_test.go)
PKGS=$(shell find src -type d -depth 1 |xargs -n1 basename)

solution: $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o $@ src/*.go

test:
	@for test in testdata/test*.txt ; do \
		expected=$$(echo $$test |sed -e 's:/test:/expected:') ; \
		if ./solution <$$test |diff $$expected - >/dev/null ; then \
			echo $$test PASS ; \
		else \
			echo $$test FAIL ; \
		fi ; \
	done

#	GOPATH=$$(PWD):$$GOPATH go test $(PKGS)
#	GOPATH=$$(PWD):$$GOPATH go test src/*.go

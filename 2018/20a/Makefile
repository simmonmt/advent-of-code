DEPS=$(shell find src -name \*.go \! -name \*_test.go)
PKGS=$(shell find src -type d -depth 1 |xargs -n1 basename)

solution: $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build -o $@ src/*.go

test:	test_input test_go

test_input:
	@for test in testdata/test*.txt ; do \
		echo $$test \\c ; \
		expected=$$(echo $$test |sed -e 's:/test:/expected:') ; \
		if ./solution <$$test |diff $$expected - >/dev/null ; then \
			echo PASS ; \
		else \
			echo FAIL ; \
		fi ; \
	done

test_go:
	GOPATH=$$(PWD):$$GOPATH go test $(PKGS)
	GOPATH=$$(PWD):$$GOPATH go test src/*.go

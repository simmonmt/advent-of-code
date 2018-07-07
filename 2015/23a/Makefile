DEPS= \
	src/reg/file.go \
	src/reg/reg.go

lock: src/lock.go $(DEPS)
	GOPATH=$$(PWD):$$GOPATH go build src/lock.go

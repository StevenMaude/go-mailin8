build:
	docker build . -t go-mailin8 
	docker run --rm --entrypoint cat go-mailin8 /go/bin/go-mailin8 > go-mailin8
	chmod u+x go-mailin8

.PHONY: build

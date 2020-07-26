build-linux:
	docker build . -t go-mailin8 
	docker run --rm --entrypoint cat go-mailin8 /go/bin/go-mailin8 > go-mailin8_linux_amd64
	chmod u+x go-mailin8_linux_amd64

build-windows:
	docker build --build-arg GOOS=windows . -t go-mailin8
	docker run --rm --entrypoint cat go-mailin8 /go/bin/windows_amd64/go-mailin8.exe > go-mailin8.exe

.PHONY: build-linux build-windows

# go-mailin8

Display the latest mail from a temporary [nada](https://getnada.com)
email address, at the command line.

Useful for getting confirmation emails for throwaway purposes, or for
testing.

## History (or "why is this called go-mailin8?")

Originally, this used Mailinator's service to retrieve email; however,
they rebuilt their site and this became a WebSocket powered application,
making it more difficult to work with.

This was a version of a [bash script I
made](https://gist.github.com/StevenMaude/914e9187c09027866fe88958798acb7e)
that uses [jq](https://stedolan.github.io/jq/). However, this Go version
requires no other dependencies, and now uses a different disposable
email service.

## Build

### Download

If you don't want to build from source, [Linux and Windows
binaries](https://github.com/StevenMaude/go-mailin8/releases) are built
from tagged versions via GitHub Actions.

### With Go installed

`go build` or `go install`

### With Docker installed

* `make build-linux` (for a Linux build)
* `make build-windows` (for a Windows build)

## Usage

1. Send, or get an email sent to a [nada](https://getnada.com) email
   address of your choosing.
2. Run the binary supplying `<email address>` as an argument, e.g.
   `go-mailin8_linux_amd64 <email_address>`

# go-mailin8

Display the latest mail from a temporary Mailinator email address, at
the command line.

Useful for getting confirmation emails for throwaway purposes, or for
testing.

This is a version of a [bash script I
made](https://gist.github.com/StevenMaude/914e9187c09027866fe88958798acb7e)
that uses [jq](https://stedolan.github.io/jq/). However, this Go version
requires no other dependencies.

## Build

`go build`

## Usage

1. Send, or get an email sent to a Mailinator email address of your
   choosing.
2. Note the local-part (the part before @, e.g.
   `somefakeemail1234@mailinator.com`).
3. Run `go-mailin8 <local-part of email address>`

## Notes

NB: you may get unpredictable results from the server if you run this
too often in a short time. See TODO.

## TODO:

* Consider introducing slight pauses between requests to try and
  improve reliability.
* Possibly consider selecting other than latest message (though outside
  of my original use case).

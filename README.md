# Aries API implementation for IIIF

This is a web service to implement the Aries API for IIIF.
It supports the following endpoints:

* / : returns version information
* /resources/iiif/pid/[PID] : returns a JSON object with some information about the image referenced by PID

### System Requirements

* GO version 1.9.2 or greater
* DEP (https://golang.github.io/dep/) version 0.4.1 or greater

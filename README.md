Reservoir-API
=============

Reservoir API is a Golang-backed JSON REST API that represents a fully functioning build system.
It will contain implementation backing for dependencies, remote workers, and storage mirroring.

Preferably, this API is the one to serve all requests incoming. If you wish to use virtual domains,
proxying with nginx is supported. We do NOT recommend Apache for an API application such as this.

Where can I find the API specifications?
========================================

You can find them in apidocs. They're written manually so that we can keep track easily of changes.
(And I don't feel like making some sort of automatic documentation tool.)

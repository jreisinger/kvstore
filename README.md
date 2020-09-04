Example of a cloud native application taken from the "Cloud Native Go" book.
kvstore is a key/value store. It stores the (resource) state by writing
transaction log into a file and then reading from it upon start.

Endpoints:

* PUT /v1/{key}
* GET /v1/{key}
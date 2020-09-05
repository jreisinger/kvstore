kvstore is a RESTful service that implements a key/value store. It persists
the (resource) state by writing transaction log into a file and then reading
from it upon start.

Endpoints:

* `PUT /v1/<key>`
* `GET /v1/<key>`
* `DELETE /v1/<key>`

Build and run the service:

```
go build
./kvstore
```

Build and run the service in Docker:

```
docker build -t kvstore .
docker run -d -p 8000:8000 kvstore
```

Usage:

```
curl localhost:8000/v1/my-key-1 -XPUT -d 'my-value1'
curl localhost:8000/v1/my-key-2 -XPUT -d 'my-value2'
curl localhost:8000/v1/my-key-2

cat transaction.log
1	2	my-key-1	my-value1
2	2	my-key-2	my-value2
```

kvstore is an example of a cloud native application and it's mostly taken
from the [Cloud Native
Go](https://learning.oreilly.com/library/view/cloud-native-go/9781492076322/)
book.

kvstore is a minimum viable product with these possible improvements:

* There's no `Close` method to gracefully close the transactions file.
* The application can close with events still in the write buffer: events can
  get lost.
* Keys and values aren't encoded in the transaction log: multiple lines or
  whitespace will fail to parse correctly.
* The sizes of keys and values are unbound: huge keys/values can be added,
  filling the disk.
* The transaction log is written in plain text: it will take up more space than
  it probably needs to.
* The log retains records of deleted values forever: it will grow indefinitely.
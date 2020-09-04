Example of a cloud native application taken from the [Cloud Native
Go](https://learning.oreilly.com/library/view/cloud-native-go/9781492076322/)
book. kvstore is a key/value store. It stores the (resource) state by writing
transaction log into a file and then reading from it upon start.

Endpoints:

* PUT /v1/{key}
* GET /v1/{key}

Usage:

```
curl localhost:8000/v1/my-key-1 -XPUT -d 'my-value1'
curl localhost:8000/v1/my-key-2 -XPUT -d 'my-value2'
curl localhost:8000/v1/my-key-2

cat transaction.log
1	2	my-key-1	my-value1
2	2	my-key-2	my-value2
```

See also Kelsey Hightower's [memkv](https://github.com/kelseyhightower/memkv).
ApiProxy
========

Created for fun, inspired by [https://www.youtube.com/watch?v=sVpMc0hwqps](https://www.youtube.com/watch?v=sVpMc0hwqps)

<img src="http://i.imgur.com/atfrFO2.png" />

api-proxy has one primary api endpoint: `/call`. Example request body:
```json
[{"Url": "http://localhost:8002/users/slawosz/repos", "Method": "Get"},
{"Url": "http://localhost:8002/users/baisa/repos", "Method": "Get"},
{"Url": "http://localhost:8002/users/lukesarnacki/repos", "Method": "Get"}]
```
(right now api-proxy starts fake server on 8002)

Response is array with response body for given call. If request timed out,
url link will be provided to get delayed body (comming soon).


Installation
============


```
export GOPATH=`pwd`
$ make
```

Running
=======

```
bin/api-proxy-server
```

Testing
=======


Httpie:
```
http POST localhost:8001/call < bodies/simple.json
```

Curl (to measure responses):
```
curl -w "@curl-format.txt" -s http://localhost:8001/call  -H "Content-Type: application/json" --data @bodies/simple.json
```

Why?
====
* lot of fun with go
* focus on better api design - more smaller calls make api better designed
* save resources (and battery) on mobile client
* faster calls
* less time to get all api calls

TODO:
====
* Limit number of requests
* Add timeout for every request

// maybe
* support for hypermedia api ie. Github

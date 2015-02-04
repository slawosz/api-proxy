ApiProxy
========

Created for fun, inspired by (paste youtube link here)

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
* focus on better api design - more smaller calls make api better designed
* save resources (and battery) on mobile client
* faster calls
* less time to get all api calls

// maybe
* support for hypermedia api ie. Github

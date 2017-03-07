# SSL Muxy Tests

Tests the following features:

* Run Proxy with HTTPS enabled
* Run Proxy with HTTPS enabled + custom certificate
* Proxy HTTPS target
* Proxy HTTPS target with invalid (untrusted) certificate
* Proxy HTTPS target requiring client certificates


### Start MASSL server

```
cd examples/ssl/massl-server
go run main.go
```

From this directory, you should be able to `curl` the server to ensure it's up:

```
curl --cacert ca.pem -E ./client.p12:password https://localhost:8080/hello
# responds with "hello, world!"
```

### Start Muxy

```
cd examples/ssl
muxy proxy --config certificate.yml
```

### cURL muxy

```
curl -k -v https://localhost:8000/hello
```

You should see "Server certificate: localhost" if the correct certificates are being used.

### Add some chaos

Now that you have things working, time to add some chaos - uncomment the `http_tamperer`
in `certificate.yml`:

```
## HTTP Tamperer - Messes with Layer 7.
##
## Useful for messing with the HTTP protocol
##
- name: http_tamperer
  config:
    request:
      path:   "/nothello"
      body:   "wow, new body!" # Override request body
    response:
      status: 201              # Override HTTP Status code
      body:   "my new body"    # Override response body
```

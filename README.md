These are some tools I've been using to test performance issues related to long running HTTP requests or TCP/SSL handshake

Usage:
```sh
concurent-http-requests <requests> <concurrency> <url>
```
```sh
concurrent-ssl-only <requests> <concurrency> <host:port>
```

Depending where the tools will run, you should usually make sure that open file limits is high enough

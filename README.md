These are some tools I've been using to test performance issues related to long-running HTTP requests or TCP/SSL handshake

Usage:
```sh
concurrent-http-requests <requests> <concurrency> <url>
```
```sh
concurrent-ssl-only <requests> <concurrency> <host:port>
```

Depending on where the tools will run, you should usually make sure that open file limits is high enough

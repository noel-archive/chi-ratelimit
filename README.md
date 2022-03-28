# 🔬 chi-ratelimit • [![Lint](https://github.com/Noelware/chi-ratelimit/actions/workflows/lint.yml/badge.svg)](https://github.com/Noelware/chi-ratelimit/actions/workflows/lint.yml) [![GoDoc](https://godoc.org/github.com/noelware/chi-ratelimit?status.svg)](http://godoc.org/github.com/noelware/chi-ratelimit) [![Go Report Card](https://goreportcard.com/badge/github.com/noelware/chi-ratelimit)](https://goreportcard.com/report/github.com/noelware/chi-ratelimit)

> *Simple production-ready ratelimiter for Chi applications*

## What is this?
**chi-ratelimit** is middleware to implement production-safe ratelimiting into your [go-chi](https://github.com/go-chi/chi) applications.

## Providers
- In-memory
- [Redis](https://redis.io) under [`chi-ratelimit/redis`](https://github.com/Noelware/chi-ratelimit-redis)
- [etcd](https://etcd.io) under [`chi-ratelimit/etcd`](https://github.com/Noelware/chi-ratelimit-etcd)

## How to use?
```shell
$ go get github.com/noelware/chi-ratelimit
```

## License
**chi-ratelimit** is released under the **MIT License** by Noelware.


![Go](https://github.com/dalmarcogd/gwp/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/dalmarcogd/gwp/branch/master/graph/badge.svg)](https://codecov.io/gh/dalmarcogd/go-worker-pool)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/beee52f22195471abea544a19ee6304a)](https://www.codacy.com/manual/dalmarco.gd/go-worker-pool?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=dalmarcogd/go-worker-pool&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/dalmarcogd/go-worker-pool)](https://goreportcard.com/report/github.com/dalmarcogd/go-worker-pool)

# gwp

This package wants to offer the community and implement workers with the pure Go code for Golangers, without any other dependency just Uuid. It allows you to expose an http server to answer the response of health checks, stats, debug pprof and the main "workers". Workers for consumer queues, channel processes and other things that you think worker needs.

## Prerequisites
Golang version >= [1.14](https://golang.org/doc/devel/release.html#go1.14)

## Features
- Setup http server to monitoring yours;
  - /stats with workers, showing statuses her, number of goroutines, number of cpus and more;
  - /health-check that look for status of workers;
  - /debug/pprof expose all endpoints of investivate golang runtime [http](https://golang.org/pkg/net/http/pprof/);
- Allow multiple concurrencies of work, handle errors and restart always worker;

## Documentation
TODO

## Examples

#### [Simple Worker](https://github.com/dalmarcogd/test-go-worker-pool/blob/master/simpleWorker.go) ###

```go
```

#### [Simple Worker Consume SQS](https://github.com/dalmarcogd/test-go-worker-pool/blob/master/simpleWorkerConsumeSQS.go) ###
```go
```

#### [Simple Worker Consume Rabbit](https://github.com/dalmarcogd/test-go-worker-pool/blob/master/simpleWorkerConsumeRabbit.go) ###
```go
```


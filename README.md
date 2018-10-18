# runningvariance (Go)

[![GoDoc](https://godoc.org/github.com/andreyvit/runningvariance?status.svg)](https://godoc.org/github.com/andreyvit/runningvariance)

Computes running mean, variance, standard deviation, skewness and kurtosis
using O(1) memory.

Install:

```sh
go get github.com/andreyvit/runningvariance
```

Example:

```go
var s runningvariance.Stat
s.Push(0)
s.Push(2)
s.Push(4)
fmt.Println("Mean:", s.Mean())
fmt.Println("StdDev:", s.StdDev())
```

See [the docs](https://godoc.org/github.com/andreyvit/runningvariance) for more.
### What version of Go are you using (`go version`)?

<pre>
$ go version
go version go1.18.4 darwin/amd64
</pre>

Also present on older versions *(tested `go1.16` and `go1.17.11` as well)*.

### Does this issue reproduce with the latest release?

Yes.

### What operating system and processor architecture are you using (`go env`)?

Easily reproducible *(after a few tries)* on macOS Monterey *(both on ARM and Intel-based CPUs, tried on both macOS 12.2 and 12.4)*.

<details><summary><code>go env</code> Output</summary><br><pre>
$ go env
GO111MODULE="auto"
GOARCH="amd64"
GOBIN=""
GOCACHE="/Users/vagrant/Library/Caches/go-build"
GOENV="/Users/vagrant/Library/Application Support/go/env"
GOEXE=""
GOEXPERIMENT=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="darwin"
GOINSECURE=""
GOMODCACHE="/go/pkg/mod"
GONOPROXY=""
GONOSUMDB=""
GOOS="darwin"
GOPATH="/go"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOROOT="/usr/local/go"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/usr/local/go/pkg/tool/darwin_amd64"
GOVCS=""
GOVERSION="go1.18.4"
GCCGO="gccgo"
GOAMD64="v1"
AR="ar"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
GOMOD=""
GOWORK=""
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -arch x86_64 -m64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=/var/folders/11/nh0v1jld7zd7b9zqm1774gtm0000gn/T/go-build810501022=/tmp/go-build -gno-record-gcc-switches -fno-common"
</pre></details>

I couldn't reproduce it on Linux (`Ubuntu 20.04`) and macOS Catalina (`10.15.7`) but I could very infrequently reproduce it on macOS Big Sur (`11.6.8`) as well.

### What did you do?

Run https://go.dev/play/p/r2qq7xIrBj9 on macOS Big Sur/Monterey via `for i in {1..100}; do go run main.go; echo $i; done`, it usually hangs after ~4-5 runs.

### What did you expect to see?

Command execution and HTTP request-related goroutines terminate, the main goroutine terminates, and the process exits.

### What did you see instead?

TODO:
- The process hangs forever, doesn't terminate.
- It *(the main, parent process)* is consuming 100% CPU.
- On `cmd.Run()` it succesfully `fork`-ed the parent process but the child never `execve`-ed.
- stacktrace
- `dtruss -f` of parent
- `CGO_ENABLED=0` fixes the issue -> something to do with macOS C libs/CGO integration on macOS is buggy? Couldn't reproduce it from C (link)
- check if enabling CGO but switching to Go-based DNS solves it or not

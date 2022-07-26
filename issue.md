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

We couldn't reproduce it on Linux (`Ubuntu 20.04`) and macOS Catalina (`10.15.7`) but we could very infrequently reproduce it on macOS Big Sur (`11.6.8`) as well.

### What did you do?

Run https://go.dev/play/p/3gaoK8kU6Qb on macOS Big Sur/Monterey via `for i in {1..100}; do go run main.go; echo $i; done`, it usually hangs after ~4-5 runs.

### What did you expect to see?

Command execution and HTTP request-related goroutines terminate, the main goroutine terminates, and the process exits.

### What did you see instead?

- The process hangs forever, doesn't terminate.
- It *(the main, parent process)* is consuming 100% CPU.
- On `cmd.Run()` it succesfully `fork`-ed the parent process but the child never `execve`-ed *(it still points to the parent process's binary)*.
- <details><summary>Stacktrace of hung process.</summary><br><pre>SIGQUIT: quit
    PC=0x7ff8061223ea m=0 sigcode=0

    goroutine 0 [idle]:
    runtime.pthread_cond_wait(0x1450a00, 0x14509c0)
        /usr/local/go/src/runtime/sys_darwin.go:448 +0x34
    runtime.semasleep(0xffffffffffffffff)
        /usr/local/go/src/runtime/os_darwin.go:66 +0xad
    runtime.notesleep(0x14507c8)
        /usr/local/go/src/runtime/lock_sema.go:181 +0x85
    runtime.mPark(...)
        /usr/local/go/src/runtime/proc.go:1449
    runtime.stopm()
        /usr/local/go/src/runtime/proc.go:2228 +0x8d
    runtime.findrunnable()
        /usr/local/go/src/runtime/proc.go:2804 +0x865
    runtime.schedule()
        /usr/local/go/src/runtime/proc.go:3187 +0x239
    runtime.exitsyscall0(0xc000166820)
        /usr/local/go/src/runtime/proc.go:3938 +0x15b
    runtime.mcall()
        /usr/local/go/src/runtime/asm_amd64.s:425 +0x43

    goroutine 1 [semacquire]:
    sync.runtime_Semacquire(0xc000013490?)
        /usr/local/go/src/runtime/sema.go:56 +0x25
    sync.(*WaitGroup).Wait(0x1004631?)
        /usr/local/go/src/sync/waitgroup.go:136 +0x52
    main.main()
        /Users/vagrant/hanging-switch-workflow/main.go:46 +0x11d

    goroutine 250 [select]:
    net/http.(*persistConn).writeLoop(0xc000398000)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 260 [select]:
    net/http.(*persistConn).writeLoop(0xc000316120)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 251 [IO wait]:
    internal/poll.runtime_pollWait(0x1774568, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c6f00?, 0xc0003c2000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c6f00, {0xc0003c2000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c6f00, {0xc0003c2000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e388, {0xc0003c2000?, 0xc000010338?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000108fc0, {0xc0003c2000?, 0xc0003ad200?, 0xc00030bd30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039bc80)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039bc80, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000108fc0)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 213 [IO wait]:
    internal/poll.runtime_pollWait(0x1774ec8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0000ca200?, 0xc000436000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0000ca200, {0xc000436000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0000ca200, {0xc000436000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc0000103b0, {0xc000436000?, 0xc000010278?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc0001086c0, {0xc000436000?, 0xc00017a9c0?, 0xc0003a5d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc000179380)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc000179380, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc0001086c0)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 34 [syscall]:
    syscall.syscall(0xc00018cb18?, 0x106de0f?, 0xc00018cbe8?, 0xc0001a6400?)
        /usr/local/go/src/runtime/sys_darwin.go:22 +0x4e
    syscall.readlen(0x14802c0?, 0x144ff20?, 0x23?)
        /usr/local/go/src/syscall/syscall_darwin.go:234 +0x35
    syscall.forkExec({0x1278a7d?, 0xc00019df70?}, {0xc000182120?, 0x1?, 0x1?}, 0xc0001ad000?)
        /usr/local/go/src/syscall/exec_unix.go:221 +0x44c
    syscall.StartProcess(...)
        /usr/local/go/src/syscall/exec_unix.go:255
    os.startProcess({0x1278a7d, 0xe}, {0xc000182120, 0x1, 0x1}, 0xc00018cef8)
        /usr/local/go/src/os/exec_posix.go:54 +0x335
    os.StartProcess({0x1278a7d, 0xe}, {0xc000182120, 0x1, 0x1}, 0x100e045?)
        /usr/local/go/src/os/exec.go:109 +0x5a
    os/exec.(*Cmd).Start(0xc000190420)
        /usr/local/go/src/os/exec/exec.go:425 +0x5e5
    os/exec.(*Cmd).Run(0x1278a7d?)
        /usr/local/go/src/os/exec/exec.go:338 +0x1e
    main.main.func1()
        /Users/vagrant/hanging-switch-workflow/main.go:21 +0xd4
    created by main.main
        /Users/vagrant/hanging-switch-workflow/main.go:16 +0x33

    goroutine 208 [select]:
    net/http.(*persistConn).writeLoop(0xc000398120)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 244 [select]:
    net/http.(*persistConn).writeLoop(0xc000316240)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 254 [select]:
    net/http.(*persistConn).writeLoop(0xc000109320)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 226 [select]:
    net/http.(*persistConn).writeLoop(0xc000109440)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 255 [IO wait]:
    internal/poll.runtime_pollWait(0x1774748, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c6e80?, 0xc0003c6000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c6e80, {0xc0003c6000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c6e80, {0xc0003c6000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e398, {0xc0003c6000?, 0xc000010290?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc0001087e0, {0xc0003c6000?, 0xc0003ad2c0?, 0xc0003d8d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039be00)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039be00, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc0001087e0)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 215 [IO wait]:
    internal/poll.runtime_pollWait(0x1774a18, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c6d80?, 0xc000438000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c6d80, {0xc000438000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c6d80, {0xc000438000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc0000103c0, {0xc000438000?, 0xc000010308?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000108d80, {0xc000438000?, 0xc00017ac60?, 0xc00043fd30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc000179440)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc000179440, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000108d80)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 214 [select]:
    net/http.(*persistConn).writeLoop(0xc0001086c0)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 145 [IO wait]:
    internal/poll.runtime_pollWait(0x1774bf8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c6d00?, 0xc000432000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c6d00, {0xc000432000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c6d00, {0xc000432000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc0000103a0, {0xc000432000?, 0xc00018e330?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000398240, {0xc000432000?, 0xc00017a900?, 0xc0000bad30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc000179200)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc000179200, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000398240)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 257 [IO wait]:
    internal/poll.runtime_pollWait(0x1774838, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c6e00?, 0xc0003c8000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c6e00, {0xc0003c8000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c6e00, {0xc0003c8000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e3a0, {0xc0003c8000?, 0xc000010320?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000108ea0, {0xc0003c8000?, 0xc0003ad320?, 0xc0003dad30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039bec0)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039bec0, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000108ea0)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 247 [IO wait]:
    internal/poll.runtime_pollWait(0x1774478, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c6f80?, 0xc0003be000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c6f80, {0xc0003be000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c6f80, {0xc0003be000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e378, {0xc0003be000?, 0xc0000102a8?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000108900, {0xc0003be000?, 0xc0003ad140?, 0xc0003a2d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039bb00)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039bb00, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000108900)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 245 [IO wait]:
    internal/poll.runtime_pollWait(0x1774298, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c7000?, 0xc0003bc000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c7000, {0xc0003bc000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c7000, {0xc0003bc000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e370, {0xc0003bc000?, 0xc000010350?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc0001090e0, {0xc0003bc000?, 0xc0003ad0e0?, 0xc00031fd30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039ba40)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039ba40, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc0001090e0)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 203 [IO wait]:
    internal/poll.runtime_pollWait(0x17741a8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c7080?, 0xc0003b2000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c7080, {0xc0003b2000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c7080, {0xc0003b2000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e348, {0xc0003b2000?, 0xc0000102c0?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000108a20, {0xc0003b2000?, 0xc000344060?, 0xc0003a3d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039b680)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039b680, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000108a20)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 209 [IO wait]:
    internal/poll.runtime_pollWait(0x17740b8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c7100?, 0xc0003b8000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c7100, {0xc0003b8000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c7100, {0xc0003b8000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e360, {0xc0003b8000?, 0xc000010368?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000109200, {0xc0003b8000?, 0xc00017aba0?, 0xc0001b3d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039b8c0)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039b8c0, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000109200)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 211 [IO wait]:
    internal/poll.runtime_pollWait(0x1774dd8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0001c6c80?, 0xc000434000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0001c6c80, {0xc000434000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0001c6c80, {0xc000434000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc0000103a8, {0xc000434000?, 0xc0002862a0?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000316000, {0xc000434000?, 0xc00017a960?, 0xc00004bd30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc0001792c0)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc0001792c0, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000316000)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 259 [IO wait]:
    internal/poll.runtime_pollWait(0x1774b08, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0000ca300?, 0xc0003ca000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0000ca300, {0xc0003ca000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0000ca300, {0xc0003ca000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e3a8, {0xc0003ca000?, 0xc000286348?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000316120, {0xc0003ca000?, 0xc0003ad380?, 0xc0003d4d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039bf80)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039bf80, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000316120)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 249 [IO wait]:
    internal/poll.runtime_pollWait(0x1774388, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0000ca480?, 0xc0003c0000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0000ca480, {0xc0003c0000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0000ca480, {0xc0003c0000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e380, {0xc0003c0000?, 0xc00018e300?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000398000, {0xc0003c0000?, 0xc0003ad1a0?, 0xc0000b6d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039bbc0)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039bbc0, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000398000)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 253 [IO wait]:
    internal/poll.runtime_pollWait(0x1774928, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0000ca380?, 0xc0003c4000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0000ca380, {0xc0003c4000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0000ca380, {0xc0003c4000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e390, {0xc0003c4000?, 0xc000010380?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000109320, {0xc0003c4000?, 0xc0003ad260?, 0xc000321d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039bd40)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039bd40, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000109320)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 243 [IO wait]:
    internal/poll.runtime_pollWait(0x1774658, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0000ca400?, 0xc0003ba000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0000ca400, {0xc0003ba000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0000ca400, {0xc0003ba000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e368, {0xc0003ba000?, 0xc000286360?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000316240, {0xc0003ba000?, 0xc00017ac00?, 0xc00043ed30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039b980)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039b980, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000316240)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 161 [IO wait]:
    internal/poll.runtime_pollWait(0x1774ce8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0000ca280?, 0xc000342000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0000ca280, {0xc000342000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0000ca280, {0xc000342000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc000286368, {0xc000342000?, 0xc000010398?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000109440, {0xc000342000?, 0xc00017aa20?, 0xc0000bbd30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00032b740)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00032b740, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000109440)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 207 [IO wait]:
    internal/poll.runtime_pollWait(0x1773fc8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0000ca500?, 0xc0003b6000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0000ca500, {0xc0003b6000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0000ca500, {0xc0003b6000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e358, {0xc0003b6000?, 0xc00018e318?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000398120, {0xc0003b6000?, 0xc0003ad3e0?, 0xc0003d6d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039b800)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039b800, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000398120)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 199 [IO wait]:
    internal/poll.runtime_pollWait(0x17c3fb8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc0002b5100?, 0xc0003ae000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc0002b5100, {0xc0003ae000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc0002b5100, {0xc0003ae000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e338, {0xc0003ae000?, 0xc0000102d8?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000108b40, {0xc0003ae000?, 0xc00017ab40?, 0xc0001b4d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039b500)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039b500, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000108b40)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 205 [IO wait]:
    internal/poll.runtime_pollWait(0x17c4198, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc00014e180?, 0xc0003b4000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc00014e180, {0xc0003b4000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc00014e180, {0xc0003b4000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e350, {0xc0003b4000?, 0xc0000102f0?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc000108c60, {0xc0003b4000?, 0xc0003440c0?, 0xc0003dcd30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039b740)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039b740, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc000108c60)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 201 [IO wait]:
    internal/poll.runtime_pollWait(0x17c40a8, 0x72)
        /usr/local/go/src/runtime/netpoll.go:302 +0x89
    internal/poll.(*pollDesc).wait(0xc00014e200?, 0xc0003b0000?, 0x0)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:83 +0x32
    internal/poll.(*pollDesc).waitRead(...)
        /usr/local/go/src/internal/poll/fd_poll_runtime.go:88
    internal/poll.(*FD).Read(0xc00014e200, {0xc0003b0000, 0x1000, 0x1000})
        /usr/local/go/src/internal/poll/fd_unix.go:167 +0x25a
    net.(*netFD).Read(0xc00014e200, {0xc0003b0000?, 0x1006b89?, 0x4?})
        /usr/local/go/src/net/fd_posix.go:55 +0x29
    net.(*conn).Read(0xc00018e340, {0xc0003b0000?, 0xc000010040?, 0x1?})
        /usr/local/go/src/net/net.go:183 +0x45
    net/http.(*persistConn).Read(0xc0001085a0, {0xc0003b0000?, 0xc0003ad440?, 0xc0003d7d30?})
        /usr/local/go/src/net/http/transport.go:1929 +0x4e
    bufio.(*Reader).fill(0xc00039b5c0)
        /usr/local/go/src/bufio/bufio.go:106 +0x103
    bufio.(*Reader).Peek(0xc00039b5c0, 0x1)
        /usr/local/go/src/bufio/bufio.go:144 +0x5d
    net/http.(*persistConn).readLoop(0xc0001085a0)
        /usr/local/go/src/net/http/transport.go:2093 +0x1ac
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1750 +0x173e

    goroutine 212 [select]:
    net/http.(*persistConn).writeLoop(0xc000316000)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 210 [select]:
    net/http.(*persistConn).writeLoop(0xc000398240)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 216 [select]:
    net/http.(*persistConn).writeLoop(0xc000108d80)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 258 [select]:
    net/http.(*persistConn).writeLoop(0xc000108ea0)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 256 [select]:
    net/http.(*persistConn).writeLoop(0xc0001087e0)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 252 [select]:
    net/http.(*persistConn).writeLoop(0xc000108fc0)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 248 [select]:
    net/http.(*persistConn).writeLoop(0xc000108900)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 246 [select]:
    net/http.(*persistConn).writeLoop(0xc0001090e0)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 204 [select]:
    net/http.(*persistConn).writeLoop(0xc000108a20)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 242 [select]:
    net/http.(*persistConn).writeLoop(0xc000109200)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 206 [select]:
    net/http.(*persistConn).writeLoop(0xc000108c60)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 202 [select]:
    net/http.(*persistConn).writeLoop(0xc0001085a0)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    goroutine 200 [select]:
    net/http.(*persistConn).writeLoop(0xc000108b40)
        /usr/local/go/src/net/http/transport.go:2392 +0xf5
    created by net/http.(*Transport).dialConn
        /usr/local/go/src/net/http/transport.go:1751 +0x1791

    rax    0x104
    rbx    0x5044600
    rcx    0x7ff7bfefec58
    rdx    0xb00
    rdi    0x1450a00
    rsi    0x4f0100005000
    rbp    0x7ff7bfefed00
    rsp    0x7ff7bfefec58
    r8     0x0
    r9     0xa0
    r10    0x0
    r11    0x246
    r12    0x16
    r13    0x0
    r14    0x4f0100005000
    r15    0xb00
    rip    0x7ff8061223ea
    rflags 0x247
    cs     0x7
    fs     0x0
    gs     0x0
    exit status 2
    </pre></details>
- `dtruss` on the parent process doesn't report any syscalls. `dtruss` on the child process *(created by `cmd.Run()`)* reports `__semwait_signal(0x903, 0x0, 0x1) = -1 Err#60` every few seconds.
- `CGO_ENABLED=0` fixes the issue.
    - Following up on this, we've tried to reproduce the issue via pure C *(to prove that this is a macOS issue and it has nothing to do with Go)*, trying to use the same syscalls as Go, but we failed to do so *(https://gist.github.com/szabolcsgelencser/b5c0ff22765edc155c6e93cb1f46e338)*. Our latest assumption was that it has something to do with Go's CGO handling *(acquiring OS threads to execute the C code on them, etc)* but we reached the limit of our technical knowledge in this area.
    - Forcing usage of the Go-based DNS resolver *(`cgo` is the default)* via `GODEBUG=netdns=go` fixes the issue.
- `unset`-ig the `__CF_USER_TEXT_ENCODING` environment variable fixes the issue *(similar to/same as https://github.com/golang/go/issues/52086 ?)*.

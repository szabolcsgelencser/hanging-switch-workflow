package main

import (
	"context"
	"net"
	"net/http"
	"os/exec"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "./whatever.exe")
			cmd.Run()
		}()
	}

	var wg2 sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()

			var netTransport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			}
			var netClient = &http.Client{
				Timeout:   time.Second * 10,
				Transport: netTransport,
			}
			netClient.Get("http://google.com")
		}()
	}

	wg.Wait()
	wg2.Wait()
}

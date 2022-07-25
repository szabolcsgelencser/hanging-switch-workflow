package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	go func() {
		var n int
		for {
			ctx := context.Background()
			host := fmt.Sprintf("dns-lookup-%d", n)
			_, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				// fmt.Printf("error: lookup host: %s\n", err)
			} else {
				// fmt.Printf("dns addresses: %+v\n", addrs)
			}
			time.Sleep(100 * time.Millisecond)
			n++
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	var n int
	for {
		select {
		case <-ctx.Done():
			fmt.Println("terminating")
			return
		default:
		}

		n++
		cwd, _ := os.Getwd()
		if err := EnvmanInitAtPath(filepath.Join(cwd, "envstore")); err != nil {
			fmt.Printf("error: envman init: %s\n", err)
		} else if n%1000 == 0 {
			fmt.Printf("envman init succeeded: %d\n", n)
		}
	}
}

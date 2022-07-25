package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
)

func main() {
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

		go func() {
			ctx := context.Background()
			host := fmt.Sprintf("dns-lookup-%d", n)
			_, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				// fmt.Printf("error: lookup host: %s\n", err)
			} else {
				// fmt.Printf("dns addresses: %+v\n", addrs)
			}
		}()

		cwd, _ := os.Getwd()
		if err := EnvmanInitAtPath(filepath.Join(cwd, "envstore")); err != nil {
			fmt.Printf("error: envman init: %s\n", err)
		} else if n%1000 == 0 {
			fmt.Printf("envman init succeeded: %d\n", n)
		}
	}
}

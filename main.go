package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/bitrise-io/go-utils/v2/log"
)

func main() {
	go func() {
		if err := http.ListenAndServe("localhost:50001", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			latency := rand.Float64()
			time.Sleep(time.Duration(latency) * time.Second)

			if rand.Intn(2) == 0 {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})); err != nil {
			fmt.Printf("error: listen and serve: %s\n", err)
		}
	}()

	client := NewDefaultClient(log.NewLogger())
	for n := 0; n < 10; n++ {
		go func() {
			for {
				client.Send(bytes.NewBufferString(`{"k":"v"}`))
				time.Sleep(time.Millisecond)
			}
		}()
	}

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

package main

import (
	"net/http"
)

func main() {
	resp, _ := http.Get("http://google.com")
	resp.Body.Close()
}

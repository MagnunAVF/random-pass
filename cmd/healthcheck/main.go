package main

import (
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:3000/")
	if err != nil || resp.StatusCode >= 500 {
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello, world!")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
	http.ListenAndServe(":9000", nil)
}

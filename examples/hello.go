package main

import (
	"fmt"
	"net/http"

	"github.com/JaCoB1123/web"
)

func hello(val string) string {
	return "hello " + val + "\n"
}

func helloInt(intval int, val string) string {
	return fmt.Sprintf("hello %s (%d)\n", val, intval)
}

func main() {
	server := web.NewServer()
	server.Get("/([-+]?[0-9]*)/(.*)", helloInt)
	server.Get("/(.*)", hello)

	http.ListenAndServe(":9999", server)
}

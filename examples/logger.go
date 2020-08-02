package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JaCoB1123/web"
)

func hello(val string) string {
	return "hello " + val + "\n"
}

func main() {
	f, err := os.Create("server.log")
	if err != nil {
		println(err.Error())
		return
	}
	logger := log.New(f, "", log.Ldate|log.Ltime)

	server := web.NewServer()
	server.Get("/(.*)", hello)
	server.SetLogger(logger)
	http.ListenAndServe("0.0.0.0:9999", server)
}

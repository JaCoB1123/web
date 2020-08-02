package main

import (
	"net/http"

	"github.com/JaCoB1123/web"
)

func hello1(val string) string {
	return "hello1 " + val + "\n"
}

func hello2(val string) string {
	return "hello2 " + val + "\n"
}

func main() {
	server1 := web.NewServer()
	server2 := web.NewServer()

	server1.Get("/(.*)", hello1)
	go http.ListenAndServe("0.0.0.0:9999", server1)
	server2.Get("/(.*)", hello2)
	go http.ListenAndServe("0.0.0.0:8999", server2)
	<-make(chan int)
}

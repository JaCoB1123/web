package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JaCoB1123/web"
)

func hello(ctx *web.Context, num string) {
	flusher, _ := ctx.ResponseWriter.(http.Flusher)
	flusher.Flush()
	n, _ := strconv.ParseInt(num, 10, 64)
	for i := int64(0); i < n; i++ {
		ctx.WriteString("<p>hello world</p>")
		flusher.Flush()
		time.Sleep(1e9)
	}
}

func main() {
	server := web.NewServer()
	server.Get("/([0-9]+)", hello)
	http.ListenAndServe("0.0.0.0:9999", server)
}

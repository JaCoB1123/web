package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/JaCoB1123/web"
)

var form = `<form action="say" method="POST"><input name="said"><input type="submit"></form>`

var users = map[string]string{}

func main() {
	rand.Seed(time.Now().UnixNano())

	server := web.NewServer()
	server.Config.CookieSecret = "7C19QRmwf3mHZ9CPAaPQ0hsWeufKd"
	server.Get("/", func(ctx *web.Context) {
		ctx.Redirect(302, "/said")
	})
	server.Get("/said", func() string { return form })
	server.Post("/say", func(ctx *web.Context) string {
		uid := fmt.Sprintf("%d\n", rand.Int63())
		ctx.SetSecureCookie("user", uid, 3600)
		users[uid] = ctx.Params["said"]
		return `<a href="/final">Click Here</a>`
	})
	server.Get("/final", func(ctx *web.Context) string {
		uid, _ := ctx.GetSecureCookie("user")
		return "You said " + users[uid]
	})
	http.ListenAndServe("0.0.0.0:9999", server)
}

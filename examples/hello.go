package main

import (
	"fmt"

	"github.com/JaCoB1123/web"
)

func hello(val string) string {
	return "hello " + val + "\n"
}

func helloInt(intval int, val string) string {
	return fmt.Sprintf("hello %s (%d)\n", val, intval)
}

func main() {
	web.Get("/([-+]?[0-9]*)/(.*)", helloInt)
	web.Get("/(.*)", hello)
	web.Run("0.0.0.0:9999")
}

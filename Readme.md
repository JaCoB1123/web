# web.go

web.go is the simplest way to write web applications in the Go programming language. It's ideal for writing simple, performant backend web services. 

## Overview

web.go should be familiar to people who've developed websites with higher-level web frameworks like sinatra or web.py. It is designed to be a lightweight web framework that doesn't impose any scaffolding on the user. Some features include:

* Routing to url handlers based on regular expressions
* Secure cookies
* Support for fastcgi and scgi
* Web applications are compiled to native code. This means very fast execution and page render speed
* Efficiently serving static files

## Installation

Make sure you have the a working Go environment. See the [install instructions](http://golang.org/doc/install.html). web.go targets the Go `release` branch.

To install web.go, simply run:

    go get github.com/JaCoB1123/web

To compile it from source:

    git clone git://github.com/JaCoB1123/web.git
    cd web && go build

## Example
```go
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
```

To run the application, put the code in a file called hello.go and run:

    go run hello.go
    
You can point your browser to http://localhost:9999/13/world . 

### Getting parameters

Route handlers may contain a pointer to web.Context as their first parameter. This variable serves many purposes -- it contains information about the request, and it provides methods to control the http connection. For instance, to iterate over the web parameters, either from the URL of a GET request, or the form data of a POST request, you can access `ctx.Params`, which is a `map[string]string`:

```go
package main

import (
    "github.com/JaCoB1123/web"
)
    
func hello(ctx *web.Context, val string) { 
    for k,v := range ctx.Params {
		println(k, v)
	}
}   
    
func main() {
    web.Get("/(.*)", hello)
    web.Run("0.0.0.0:9999")
}
```

In this example, if you visit `http://localhost:9999/?a=1&b=2`, you'll see the following printed out in the terminal:

    a 1
    b 2

## About

Based on the awesome [web.go](https://github.com/hoisie/web) project by Michael Hoisie

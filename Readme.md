[![Build Status](https://cloud.drone.io/api/badges/JaCoB1123/web/status.svg)](https://cloud.drone.io/JaCoB1123/web)

# web

web is a simple way to write web applications in the Go programming language. It's ideal for writing simple backend web services. 

## Overview

web is designed to be a lightweight HTTP-Router that makes dependencies clear. Some features include:

* Function's dependencies clearly visible in signature 
* Routing to url handlers based on regular expressions
* Handlers can return strings to have them written as the response
* Secure cookies

## Known-Issues

If you're looking for the fastest router, this is probably not your best choice, since it uses the reflect-Package to call the handler functions and regular expressions for matching routes. If speed is your main concern, you're probably better off using something like [HttpRouter](https://github.com/julienschmidt/httprouter) (also see the [Go HTTP Router Benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)).

## Installation

Make sure you have the a working Go environment. See the [install instructions](http://golang.org/doc/install.html). web.go targets the Go `release` branch.

To install web.go, simply run:

    go get github.com/JaCoB1123/web

To compile it from source:

    git clone git://github.com/JaCoB1123/web.git
    cd web && go build

## Example

Parameters in the url are declared using regular expressions. Each group can be referenced by adding a parameter to the handler function. 

The following example sets up two routes:

- The first route requires an integer as the first URL element, which is later passed as the first argument to the handler function `helloInt`. Anything after that is accepted as a string (including more slashes) and passed as the second argument.
- The second route accepts anything that doesn't match the first rule as a string and passes that as the only parameter to the handler function `hello`.

```go
package main

import (
        "fmt"

        "github.com/JaCoB1123/web"
)

func main() {
        web.Get("/([-+]?[0-9]*)/(.*)", helloInt)
        web.Get("/(.*)", hello)
        web.Run("0.0.0.0:9999")
}

// hello requires a single string parameter in the url
func hello(val string) string {
        return "hello " + val + "\n"
}

// hello requires an integer and a string parameter in the url
func helloInt(intval int, val string) string {
    return fmt.Sprintf("hello %s (%d)\n", val, intval)
}
```

To run the application, put the code in a file called hello.go and run:

    go run hello.go
    
You can point your browser to http://localhost:9999/13/world.

### Getting parameters

Route handlers may contain a pointer to web.Context as their first parameter. This variable serves many purposes -- it contains information about the request, and it provides methods to control the http connection. This also allows direct access to the `http.ResponseWriter`. For instance, to iterate over the web parameters, either from the URL of a GET request, or the form data of a POST request, you can access `ctx.Params`, which is a `map[string]string`:

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

## Roadmap

Here's a non-exhaustive list of things I'm planning to add:
- Simple JSON-handling by returning structs
- Handling custom types as function-parameters (e.g. using an integer parameter in the URL, but have the function accept a struct, that is loaded from the database)
- Some performance improvements

## About

Based on the awesome [web.go](https://github.com/hoisie/web) project by Michael Hoisie

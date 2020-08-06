package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkWithContextHelloWorldTextonly(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/([-+]?[0-9]*)/(.*)", func(ctx *Context, intval int, val string) {
		fmt.Fprintf(ctx.ResponseWriter, "hello %s (%d)\n", val, intval)
	})
	s.Get("/(.*)", func(ctx *Context, val string) {
		fmt.Fprint(ctx.ResponseWriter, "hello "+val+"\n")
	})
	s.initServer()

	req := buildTestRequest("GET", "/world", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}
func BenchmarkWithContextHelloWorldNumber(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/([-+]?[0-9]*)/(.*)", func(ctx *Context, intval int, val string) {
		fmt.Fprintf(ctx.ResponseWriter, "hello %s (%d)\n", val, intval)
	})
	s.Get("/(.*)", func(ctx *Context, val string) {
		fmt.Fprint(ctx.ResponseWriter, "hello "+val+"\n")
	})
	s.initServer()

	req := buildTestRequest("GET", "/123456/world", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}

func BenchmarkWithContextProcessGet(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/echo/(.*)", func(ctx *Context, val string) {
		fmt.Fprint(ctx.ResponseWriter, val)
	})
	req := buildTestRequest("GET", "/echo/hi", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}

func BenchmarkWithContextProcessPost(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Post("/echo/(.*)", func(ctx *Context, val string) {
		fmt.Fprint(ctx.ResponseWriter, val)
	})
	req := buildTestRequest("POST", "/echo/hi", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}

func BenchmarkWithReturnHelloWorldTextonly(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/([-+]?[0-9]*)/(.*)", func(ctx *Context, intval int, val string) {
		fmt.Fprintf(ctx.ResponseWriter, "hello %s (%d)\n", val, intval)
	})
	s.Get("/(.*)", func(ctx *Context, val string) {
		fmt.Fprint(ctx.ResponseWriter, "hello "+val+"\n")
	})
	s.initServer()

	req := buildTestRequest("GET", "/world", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}
func BenchmarkWithReturnHelloWorldNumber(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/([-+]?[0-9]*)/(.*)", func(intval int, val string) string {
		return fmt.Sprintf("hello %s (%d)\n", val, intval)
	})
	s.Get("/(.*)", func(val string) string {
		return "hello " + val + "\n"
	})
	s.initServer()

	req := buildTestRequest("GET", "/123456/world", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}

func BenchmarkWithReturnProcessGet(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/echo/(.*)", func(s string) string {
		return s
	})
	req := buildTestRequest("GET", "/echo/hi", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}

func BenchmarkWithReturnProcessPost(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Post("/echo/(.*)", func(s string) string {
		return s
	})
	req := buildTestRequest("POST", "/echo/hi", "", nil, nil)
	var buf bytes.Buffer
	iob := ioBuffer{input: nil, output: &buf}
	c := dummyConnection{wroteHeaders: false, req: req, headers: make(map[string][]string), fd: &iob}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Process(&c, req)
	}
}

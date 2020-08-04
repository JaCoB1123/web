package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkHelloWorldTextonly(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/([-+]?[0-9]*)/(.*)", func(intval int, val string) string {
		return fmt.Sprintf("hello %s (%d)\n", val, intval)
	})
	s.Get("/(.*)", func(val string) string {
		return "hello " + val + "\n"
	})
	s.initServer()
	b.ReportAllocs()
	b.ResetTimer()

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
func BenchmarkHelloWorldNumber(b *testing.B) {
	s := NewServer()
	s.SetLogger(log.New(ioutil.Discard, "", 0))
	s.Get("/([-+]?[0-9]*)/(.*)", func(intval int, val string) string {
		return fmt.Sprintf("hello %s (%d)\n", val, intval)
	})
	s.Get("/(.*)", func(val string) string {
		return "hello " + val + "\n"
	})
	s.initServer()
	b.ReportAllocs()
	b.ResetTimer()

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

func BenchmarkProcessGet(b *testing.B) {
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

func BenchmarkProcessPost(b *testing.B) {
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

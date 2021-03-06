package web

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ServerConfig is configuration for server objects.
type ServerConfig struct {
	StaticDir    string
	Addr         string
	Port         int
	CookieSecret string
	RecoverPanic bool
	Profiler     bool
	ColorOutput  bool
}

type typeHandlerDelegate func(reflect.Type, []string, int, *Context) (reflect.Value, error)

// Server represents a web.go server.
type Server struct {
	Config       *ServerConfig
	routes       []*route
	Logger       *log.Logger
	Env          map[string]interface{}
	TypeHandlers []typeHandlerDelegate
	encKey       []byte
	signKey      []byte
}

func NewServer() *Server {
	return &Server{
		Config:       Config,
		Logger:       log.New(os.Stdout, "", log.Ldate|log.Ltime),
		Env:          map[string]interface{}{},
		TypeHandlers: []typeHandlerDelegate{getString, getInt, getContext},
	}
}

func (s *Server) initServer() {
	if s.Config == nil {
		s.Config = &ServerConfig{}
	}

	if s.Logger == nil {
		s.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}

	if s.Config.Profiler {
		s.Get("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		s.Get("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		s.Get("/debug/pprof/heap", pprof.Handler("heap"))
		s.Get("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	if len(s.Config.CookieSecret) > 0 {
		s.Logger.Println("Generating cookie encryption keys")
		s.encKey = genKey(s.Config.CookieSecret, "encryption key salt")
		s.signKey = genKey(s.Config.CookieSecret, "signature key salt")
	}
}

type route struct {
	path         string
	pathRegex    *regexp.Regexp
	method       string
	handler      reflect.Value
	httpHandler  http.Handler
	runner       func() reflect.Value
	argsBuilders []func([]string, *Context) reflect.Value
}

var dummyArgs = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func (route *route) call() reflect.Value {
	if route.httpHandler != nil {
		route.httpHandler.ServeHTTP(nil, nil)
		return reflect.Value{}
	}

	return reflect.Value{}
}

func newRouteFromHandler(pathRegex string, cr *regexp.Regexp, method string, handler http.Handler) *route {
	route := newRoute(pathRegex, cr, method)
	route.httpHandler = handler
	return route
}

func (s *Server) newRouteFromValue(pathRegex string, cr *regexp.Regexp, method string, handler reflect.Value) *route {
	route := newRoute(pathRegex, cr, method)
	route.handler = handler
	route.argsBuilders = []func([]string, *Context) reflect.Value{}

	var args []reflect.Value
	functionType := handler.Type()

	numIn := functionType.NumIn()

	iVal := 1
	for iArg := 0; iArg < numIn; iArg++ {
		arg := functionType.In(iArg)
		iValCopy := iVal

		var err error
		var result reflect.Value
		var typeHandler typeHandlerDelegate
		for i := range s.TypeHandlers {
			typeHandler = s.TypeHandlers[i]
			result, err = typeHandler(arg, dummyArgs, iVal, nil)
			if err != NotSupported {
				break
			}
		}

		route.argsBuilders = append(route.argsBuilders, func(values []string, ctx *Context) reflect.Value {
			result, _ = typeHandler(arg, values, iValCopy, ctx)
			return result
		})

		args = append(args, result)
		if err == NoValueNeeded {
			continue
		}

		iVal++
	}

	return route
}

func newRoute(pathRegex string, cr *regexp.Regexp, method string) *route {
	return &route{
		path:      pathRegex,
		pathRegex: cr,
		method:    method,
	}
}

func (s *Server) addRoute(pathRegex string, method string, handler interface{}) {
	cr, err := regexp.Compile("^" + pathRegex + "$")
	if err != nil {
		s.Logger.Printf("Error in route regex %q\n", pathRegex)
		return
	}

	switch handler.(type) {
	case http.Handler:
		s.routes = append(s.routes, newRouteFromHandler(pathRegex, cr, method, handler.(http.Handler)))
	case reflect.Value:
		fv := handler.(reflect.Value)
		s.routes = append(s.routes, s.newRouteFromValue(pathRegex, cr, method, fv))
	default:
		fv := reflect.ValueOf(handler)
		s.routes = append(s.routes, s.newRouteFromValue(pathRegex, cr, method, fv))
	}
}

// ServeHTTP is the interface method for Go's http server package
func (s *Server) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	s.Process(c, req)
}

// Process invokes the routing system for server s
func (s *Server) Process(c http.ResponseWriter, req *http.Request) {
	route := s.routeHandler(req, c)
	if route != nil {
		route.httpHandler.ServeHTTP(c, req)
	}
}

// Head adds a handler for the 'HEAD' http method for server s.
func (s *Server) Head(route string, handler interface{}) {
	s.addRoute(route, "GET", handler)
}

// Get adds a handler for the 'GET' http method for server s.
func (s *Server) Get(route string, handler interface{}) {
	s.addRoute(route, "GET", handler)
}

// Post adds a handler for the 'POST' http method for server s.
func (s *Server) Post(route string, handler interface{}) {
	s.addRoute(route, "POST", handler)
}

// Put adds a handler for the 'PUT' http method for server s.
func (s *Server) Put(route string, handler interface{}) {
	s.addRoute(route, "PUT", handler)
}

// Delete adds a handler for the 'DELETE' http method for server s.
func (s *Server) Delete(route string, handler interface{}) {
	s.addRoute(route, "DELETE", handler)
}

// Match adds a handler for an arbitrary http method for server s.
func (s *Server) Match(method string, route string, handler interface{}) {
	s.addRoute(route, method, handler)
}

// Add a custom http.Handler
func (s *Server) Handle(route string, method string, httpHandler http.Handler) {
	s.addRoute(route, method, httpHandler)
}

// safelyCall invokes `function` in recover block
func (s *Server) safelyCall(function reflect.Value, args []reflect.Value) (resp []reflect.Value, e interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if !s.Config.RecoverPanic {
				// go back to panic
				panic(err)
			} else {
				e = err
				resp = nil
				s.Logger.Println("Handler crashed with error", err)
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					s.Logger.Println(file, line)
				}
			}
		}
	}()
	return function.Call(args), nil
}

func (s *Server) logRequest(ctx *Context, sTime time.Time) {
	//log the request
	req := ctx.Request
	requestPath := req.URL.Path

	duration := time.Now().Sub(sTime)
	var client string

	// We suppose RemoteAddr is of the form Ip:Port as specified in the Request
	// documentation at http://golang.org/pkg/net/http/#Request
	pos := strings.LastIndex(req.RemoteAddr, ":")
	if pos > 0 {
		client = req.RemoteAddr[0:pos]
	} else {
		client = req.RemoteAddr
	}

	var logEntry bytes.Buffer
	logEntry.WriteString(client)
	logEntry.WriteString(" - " + s.ttyGreen(req.Method+" "+requestPath))
	logEntry.WriteString(" - " + duration.String())
	if len(ctx.Params) > 0 {
		logEntry.WriteString(" - " + s.ttyWhite(fmt.Sprintf("Params: %v\n", ctx.Params)))
	}
	ctx.Server.Logger.Print(logEntry.String())
}

func (s *Server) ttyGreen(msg string) string {
	return s.ttyColor(msg, ttyCodes.green)
}

func (s *Server) ttyWhite(msg string) string {
	return s.ttyColor(msg, ttyCodes.white)
}

func (s *Server) ttyColor(msg string, colorCode string) string {
	if s.Config.ColorOutput {
		return colorCode + msg + ttyCodes.reset
	} else {
		return msg
	}
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{Params: map[string]string{}}
	},
}

// the main route handler in web.go
// Tries to handle the given request.
// Finds the route matching the request, and execute the callback associated
// with it.  In case of custom http handlers, this function returns an "unused"
// route. The caller is then responsible for calling the httpHandler associated
// with the returned route.
func (s *Server) routeHandler(req *http.Request, w http.ResponseWriter) (unused *route) {
	ctx := contextPool.Get().(*Context)
	ctx.Reset(req, s, w)
	defer contextPool.Put(ctx)

	//ignore errors from ParseForm because it's usually harmless.
	req.ParseForm()
	if len(req.Form) > 0 {
		for k, v := range req.Form {
			ctx.Params[k] = v[0]
		}
	}

	tm := time.Now().UTC()
	defer s.logRequest(ctx, tm)

	requestPath := req.URL.Path
	for i := 0; i < len(s.routes); i++ {
		route := s.routes[i]
		cr := route.pathRegex
		//if the methods don't match, skip this handler (except HEAD can be used in place of GET)
		if req.Method != route.method && !(req.Method == "HEAD" && route.method == "GET") {
			continue
		}

		match := cr.FindStringSubmatch(requestPath)
		if match == nil || len(match[0]) != len(requestPath) {
			continue
		}

		// We can not handle custom http handlers here, give back to the caller.
		if route.httpHandler != nil {
			unused = route
			return
		}

		args := make([]reflect.Value, len(route.argsBuilders))
		for i, argBuilder := range route.argsBuilders {
			arg := argBuilder(match, ctx)
			args[i] = arg
		}

		ret, err := s.safelyCall(route.handler, args)

		// set the default content-type
		if ctx.ResponseWriter.Header().Get("Content-Type") == "" {
			ctx.SetHeader("Content-Type", "text/html; charset=utf-8", true)
		}

		if err != nil {
			//there was an error or panic while calling the handler
			ctx.Abort(500, "Server Error")
		}
		if len(ret) == 0 {
			return
		}

		sval := ret[0]

		var content []byte

		if sval.Kind() == reflect.String {
			content = []byte(sval.String())
		} else if sval.Kind() == reflect.Slice && sval.Type().Elem().Kind() == reflect.Uint8 {
			content = sval.Interface().([]byte)
		}
		ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
		_, err = ctx.ResponseWriter.Write(content)
		if err != nil {
			ctx.Server.Logger.Println("Error during write: ", err)
		}
		return
	}

	ctx.Abort(404, "Page not found")
	return
}

var NoValueNeeded = fmt.Errorf("No value needed")
var NotSupported = fmt.Errorf("Type is not supported")

func getString(t reflect.Type, values []string, valueIndex int, ctx *Context) (reflect.Value, error) {
	if t.Kind() != reflect.String {
		return reflect.Value{}, NotSupported
	}

	return reflect.ValueOf(values[valueIndex]), nil
}

func getInt(t reflect.Type, values []string, valueIndex int, ctx *Context) (reflect.Value, error) {
	intVal, err := strconv.Atoi(values[valueIndex])
	if err != nil {
		return reflect.Value{}, err
	}
	switch t.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(intVal)), nil
	case reflect.Int8:
		return reflect.ValueOf(int8(intVal)), nil
	case reflect.Int16:
		return reflect.ValueOf(int16(intVal)), nil
	case reflect.Int32:
		return reflect.ValueOf(int32(intVal)), nil
	case reflect.Int64:
		return reflect.ValueOf(int64(intVal)), nil
	case reflect.Uint:
		return reflect.ValueOf(uint(intVal)), nil
	case reflect.Uint8:
		return reflect.ValueOf(uint8(intVal)), nil
	case reflect.Uint16:
		return reflect.ValueOf(uint16(intVal)), nil
	case reflect.Uint32:
		return reflect.ValueOf(uint32(intVal)), nil
	case reflect.Uint64:
		return reflect.ValueOf(uint64(intVal)), nil
	default:
		return reflect.Value{}, NotSupported
	}
}

func getContext(t reflect.Type, values []string, valueIndex int, ctx *Context) (reflect.Value, error) {
	if t.Kind() != reflect.Ptr || t.Elem() != contextType {
		return reflect.Value{}, NotSupported
	}

	return reflect.ValueOf(ctx), NoValueNeeded
}

// SetLogger sets the logger for server s
func (s *Server) SetLogger(logger *log.Logger) {
	s.Logger = logger
}

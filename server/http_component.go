package server

import (
	"github.com/d-tar/wntr"
	"log"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"
)

//Model And View Renderable
type MavRenderable interface {
	Render(http.ResponseWriter) error
}

//Base web request processing interface
type WebController interface {
	Serve(*http.Request) MavRenderable
}

//Single function request processor interface
//
//  use it to create single-method-request-handlers'
//
//   var routes struct{
//           MyHandler  HandlerFunc `@web-uri:"/my-uri"`
//   }
//
//   routes.MyHandler = functionalHandler
//
//
type HandlerFunc func(*http.Request) MavRenderable

func _() {
	var h HandlerFunc = nil
	var _ WebController = h
}

func (this HandlerFunc) Serve(r *http.Request) MavRenderable {
	return this(r)
}


//WebServer Component
//
//   On PostInit phase it fetches all WebController components
//	that were defined with @web-url  tag and registers them
//	with http.Handler. Then it starts web server
//
//   On PreDestroy phase component destroys listener to stop
// 	http.Server's serving cycle
type WebServerComponent struct {
	listener  net.Listener
	wait      *sync.Cond
	exitError error

	Ctx 	wntr.ConfiguredContext	`inject:"t"`
}
func _() {
	var _ wntr.PostInitable = &WebServerComponent{}
}


/*
Begin Implementation
 */

type requestMapping struct {
	path    string
	handler WebController
}

var gWebControllerType reflect.Type = reflect.TypeOf( (*WebController)(nil)).Elem();


func (this *WebServerComponent) PostInit() error {
	for _, c := range this.requestMappings() {
		log.Println("WebMVC: Mapped ", c.path, " to  ", c.handler, " ")
		http.HandleFunc(c.path, createMavHandler(c.handler))
	}

	var m sync.Mutex
	this.wait = sync.NewCond(&m)
	this.wait.L.Lock()

	s := &http.Server{Addr: ":8080"}

	log.Println("Starting web server...")
	go func() {
		err := listenAndServe(s, this)
		log.Println("WebRoutine done")
		this.exitError = err
		this.wait.Broadcast()
	}()

	return nil
}

func (this *WebServerComponent) requestMappings() []requestMapping {
	r := make([]requestMapping,0)

	for _,ctl := range this.Ctx.FindComponentsByType(gWebControllerType){

		tag := ctl.Tags()
		if uri := tag.Get("@web-uri"); uri != "" {
			r = append(r, requestMapping{
				path: uri,
				handler: ctl.Instance().(WebController),
			})
		}
	}


	return r
}

func createMavHandler(c WebController) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mav := c.Serve(r)
		mav.Render(w)
	}
}

func (this *WebServerComponent) Wait() error {
	this.wait.Wait()
	return this.exitError
}

func (this *WebServerComponent) PreDestroy() {
	log.Println("Closing WebSupport http listener")
	this.listener.Close()
}

//Hack to capture listener object to perform
//server shutdown on component stop
func listenAndServe(srv *http.Server, web *WebServerComponent) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	web.listener = ln
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

//code below was coped from go's stdlib

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

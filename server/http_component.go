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

type WebController interface {
	Serve(*http.Request) MavRenderable
}

type HandlerFunc func(*http.Request) MavRenderable

func (this HandlerFunc) Serve(r *http.Request) MavRenderable {
	return this(r)
}

type WebSupport struct {
	mappings []RequestMapping

	listener  net.Listener
	wait      *sync.Cond
	exitError error
}

type RequestMapping struct {
	path    string
	handler WebController
}

func test_interfaces() {
	var _ wntr.PostInitable = &WebSupport{}
}

func (this *WebSupport) PostInit() error {
	for _, c := range this.mappings {
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

func createMavHandler(c WebController) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mav := c.Serve(r)
		mav.Render(w)
	}
}

func (this *WebSupport) Wait() error {
	this.wait.Wait()
	return this.exitError
}

func (this *WebSupport) PreDestroy() {
	log.Println("Closing WebSupport http listener")
	this.listener.Close()
}

func listenAndServe(srv *http.Server, web *WebSupport) error {
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

/*
 Web request handler autoregistration support
*/

type WebUrlMapper struct { //implements wntr.ComponentLifecycle
	Web *WebSupport
}

//This handler collects each component that was defined with tag @web-url
//and writes this data to WebSupport component
func (h *WebUrlMapper) OnPrepareComponent(c wntr.Component) error {
	if ctl, ok := c.Inst.(WebController); ok {
		tag := reflect.StructTag(c.Tags)
		if uri := tag.Get("@web-uri"); uri != "" {
			h.Web.mappings = append(h.Web.mappings,
				RequestMapping{uri, ctl},
			)
		}
	}
	return nil
}

func (h *WebUrlMapper) OnComponentReady(wntr.Component) error {
	return nil
}

func (h *WebUrlMapper) OnDestroyComponent(wntr.Component) error {

	return nil
}

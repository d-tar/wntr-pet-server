package server

import (
	"fmt"
	"net/http"
	"net/url"
)

type UserController struct{}

func (this *UserController) Serve(r *http.Request) MavRenderable {
	if err := r.ParseForm(); err == nil {
		return MavErr(fmt.Sprintf("Request died :( \n Could this explain you why?\n%v", err))
	}

	if r.Method == "GET" {
		return this.getUser(r.Form)
	} else {
		return MavErr(fmt.Sprintf("Unsupported method %v", r.Method))
	}

}

func (this *UserController) getUser(params url.Values) MavRenderable {
	var a struct {
		Id   uint32
		Name string
	}

	a.Id = 123
	a.Name = "Hello, world"

	return MavOk(a)
}

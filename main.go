package main

import (
	"github.com/d-tar/wntr"
	"github.com/d-tar/wntr-server/server"
	"log"
	"net/http"
)

//This definition provides:
//    	1. Request processing via WebController interface
//	2. ModelAndView basics via MavRenderable and JsonMavRenderable
//	3. Register WebControllers via @web-uri tag on definition
//	4. Stop web-server using context.Stop
//
//
type EnableWebSupport struct {
	Web    server.WebSupport
	Mapper server.DeclRequestMapping //NOTE: Mapper need's to know about WebSupport component
}

//List of concrete handlers
type RegisterWebHandlers struct {
	//Allows to shutdown context and app by GET /shutdown request
	ExitCtrl server.HandlerFunc `@web-uri:"/shutdown"`
	//Sample data handler
	UserCtrl server.UserController `@web-uri:"/user"`
}

//Application singleton
var app struct {
	//Include web support context
	EnableWebSupport
	//Include web-handlers context
	RegisterWebHandlers
	//Get Context to stop by shutdown handler
	Context wntr.Context `inject:"t"`
}

func main() {
	app.EnableWebSupport.Mapper.Web = &app.EnableWebSupport.Web

	app.RegisterWebHandlers.ExitCtrl = shutdownHandler

	if _, err := wntr.FastBoot(&app); err != nil {
		log.Fatal(err)
	}

	if err := app.Web.Wait(); err != nil {
		log.Fatal("Application died with errors\n", err)
	}

}

func shutdownHandler(*http.Request) server.MavRenderable{
	app.Context.Stop()
	return server.MavOk("Shutting down context")
}
package main

import (
	"github.com/d-tar/wntr"
	"github.com/d-tar/wntr-server/server"
	"log"
	"net/http"
)



//List of concrete handlers
type WebRoutes struct {
	//Allows to shutdown context and app by GET /shutdown request
	ExitCtrl server.HandlerFunc `@web-uri:"/shutdown"`
	//Sample data handler
	UserCtrl server.UserController `@web-uri:"/user"`
}

//This definition provides:
//    	1. Request processing via WebController interface
//	2. ModelAndView basics via MavRenderable and JsonMavRenderable
//	3. Register WebControllers via @web-uri tag on definition
//	4. Stop web-server using context.Stop
//
//
type EnableWebSupport struct {
	//First of all we register web handlers
	WebRoutes
	//At top we define WebServerComponent that runs WebController
	Web    server.WebServerComponent
}


//Application singleton
var app struct {
	//Include web support context
	EnableWebSupport
	//Get Context to stop by shutdown handler
	Context wntr.Context `inject:"t"`
}


//Start Application
func main() {
	//Define our web router before we begin
	app.WebRoutes.ExitCtrl = shutdownHandler

	//Create context from 'application structure'
	// then start it, and panic if something go wrong
	wntr.ContextOrPanic(&app)

	if err := app.Web.Wait(); err != nil {
		log.Fatal("Application died with errors\n", err)
	}
}

func shutdownHandler(*http.Request) server.MavRenderable{
	app.Context.Stop()
	return server.MavOk("Shutting down context")
}
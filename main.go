package main

import (
	"github.com/d-tar/wntr"
	"github.com/d-tar/wntr-server/app"
	"github.com/d-tar/wntr-server/app/controllers"
	"github.com/d-tar/wntr-server/app/services"
	"github.com/d-tar/wntr/webmvc"
	"log"
)

type EnableServices struct {
	UsersDao services.UsersDaoLedis
}

//List of concrete handlers
type WebRoutes struct {
	//Enable User API
	controllers.UserApiRoutes

	//Allows to shutdown context and app by GET /shutdown request
	ExitCtrl webmvc.HandlerFunc `@web-uri:"/shutdown"`
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
	webmvc.EnableDefaultWebMvc
}

//Application singleton
var gApp struct {
	//Register application services
	EnableServices
	//Include web support context
	EnableWebSupport
	//Get Context to stop by shutdown handler
	Context wntr.Context `inject:"t"`
}

//Start Application
func main() {
	//Define our web router before we begin
	gApp.WebRoutes.ExitCtrl = shutdownHandler

	//We need custom json view I used to
	gApp.EnableWebSupport.Mvc.SetWebViews(createViewMap())

	//Create context from 'application structure'
	// then start it, and panic if something go wrong
	wntr.ContextOrPanic(&gApp)

	if err := gApp.Web.Wait(); err != nil {
		log.Fatal("Application died with errors\n", err)
	}
}

func shutdownHandler(*webmvc.WebRequest) webmvc.WebResult {
	gApp.Context.Stop()
	return webmvc.WebOk("Shutting down context")
}

func createViewMap() map[string]webmvc.WebView {
	m := make(map[string]webmvc.WebView)
	m["JSON"] = app.NewCustomJsonView()
	m[""] = m["JSON"]
	return m
}

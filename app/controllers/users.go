package controllers

import (
	"github.com/d-tar/wntr-server/app/services"
	"github.com/d-tar/wntr/webmvc"
)

type UserApiRoutes struct {
	UserCtrl UserController
	//Sample data handler
	ListUsers webmvc.HandlerFunc `@web-uri:"/users"    @web-method:"GET"`
	GetUser   webmvc.HandlerFunc `@web-uri:"/user/:id" @web-method:"GET"`
	NewUser   webmvc.HandlerFunc `@web-uri:"/user/new" @web-method:"GET"`
}

//On context setup, let's bind web route methods to controller methods
func (this *UserApiRoutes) PreInit() error {
	this.GetUser = this.UserCtrl.GetUser
	this.NewUser = this.UserCtrl.CreateUser
	this.ListUsers = this.UserCtrl.ListUsers
	return nil
}

//Simple users handler
type UserController struct {
	UserDao services.UsersDao `inject:"t"`
}

func (this *UserController) ListUsers(r *webmvc.WebRequest) webmvc.WebResult {

	u, err := this.UserDao.ListUsers(0, 100)
	if err != nil {
		return webmvc.WebErr(err)
	}

	return webmvc.WebOk(u)
}

func (this *UserController) CreateUser(r *webmvc.WebRequest) webmvc.WebResult {

	u, err := this.UserDao.CreateUser()
	if err != nil {
		return webmvc.WebErr(err)
	}

	return webmvc.WebOk(u)
}

func (this *UserController) GetUser(r *webmvc.WebRequest) webmvc.WebResult {
	if err := r.HttpRequest.ParseForm(); err != nil {
		return webmvc.WebErr(err)
	}

	id := r.NamedParameters["id"]
	if id == "" {
		return webmvc.WebErr("Missing required parameter 'id'")
	}

	u, err := this.UserDao.GetUser(id)
	if err != nil {
		return webmvc.WebErr(err)
	}

	return webmvc.WebOk(u)
}

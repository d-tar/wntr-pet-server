package controllers

import (
	"github.com/d-tar/wntr-server/app/services"
	web "github.com/d-tar/wntr/webmvc"
	"net/http"
)

type UserApiRoutes struct {
	UserCtrl UserController
	//Sample data handler
	ListUsers web.HandlerFunc `@web-uri:"/users" @web-method:"GET"`
	GetUser   web.HandlerFunc `@web-uri:"/user" @web-method:"GET"`
	NewUser   web.HandlerFunc `@web-uri:"/user/new" @web-method:"GET"`
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

func (this *UserController) ListUsers(r *http.Request) web.WebResult {

	u, err := this.UserDao.ListUsers(0,100)
	if err != nil {
		return web.WebErr(err)
	}

	return web.WebOk(u)
}

func (this *UserController) CreateUser(r *http.Request) web.WebResult {

	u, err := this.UserDao.CreateUser()
	if err != nil {
		return web.WebErr(err)
	}

	return web.WebOk(u)
}

func (this *UserController) GetUser(r *http.Request) web.WebResult {
	if err := r.ParseForm(); err != nil {
		return web.WebErr(err)
	}

	id := r.Form.Get("id")
	if id==""{
		return web.WebErr("Missing required parameter 'id'")
	}


	u,err:=this.UserDao.GetUser(id)
	if err != nil {
		return web.WebErr(err)
	}

	return web.WebOk(u)
}

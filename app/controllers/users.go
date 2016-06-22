package controllers

import (
	"fmt"
	"github.com/d-tar/wntr-server/app/services"
	"github.com/d-tar/wntr/webmvc"
)

type UserApiRoutes struct {
	UserCtrl UserController
	//Sample data handler
	ListUsers webmvc.HandlerFunc     `@web-uri:"/users"`
	GetUser   webmvc.SmartWebHandler `@web-uri:"/user/:id" @web-method:"GET"`
	NewUser   webmvc.HandlerFunc     `@web-uri:"/user/new" @web-method:"POST"`
	SaveUser  webmvc.SmartWebHandler `@web-uri:"/test/:id"`
}

//On context setup, let's bind web route methods to controller methods
func (this *UserApiRoutes) PreInit() error {
	this.GetUser = webmvc.AutoHandler(this.UserCtrl.GetUser)
	this.NewUser = this.UserCtrl.CreateUser
	this.ListUsers = this.UserCtrl.ListUsers
	this.SaveUser = webmvc.AutoHandler(this.UserCtrl.UpdateUser)
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

func (this *UserController) GetUser(Request struct {
	Id string `@path-variable:"id"`
}) (*services.User, error) {
	if Request.Id == "" {
		return nil, fmt.Errorf("Missing required parameter 'id'")
	}

	u, err := this.UserDao.GetUser(Request.Id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

//Request Scoped Injection Sample
//State fields are autowired using WebMvc configurer
//and standard autowiring component
func (this *UserController) UpdateUser(state struct {
	UserId   string        `@path-variable:"id"`
	UserData services.User `@request-body:""`
}) (interface{}, error) {

	user, err := this.UserDao.GetUser(state.UserId)
	if err != nil {
		return nil, err
	}

	UserData := state.UserData

	if UserData.Name != "" {
		user.Name = UserData.Name
	}

	if err := this.UserDao.SaveUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

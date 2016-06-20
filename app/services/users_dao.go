package services

import (
	"encoding/json"
	"fmt"
	"github.com/d-tar/wntr"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
	"math/rand"
	"strconv"
	"time"
)

type UsersDao interface {
	GetUser(id string) (*User, error)
	//UpdateUser(*User) error
	//DeleteUser(*User) error
	CreateUser() (*User, error)
	ListUsers(from int, to int) ([]*User, error)
}

type UsersDaoLedis struct {
	ledis *ledis.Ledis
	db    *ledis.DB
}

func _() {
	var _ wntr.PostInitable = &UsersDaoLedis{}
	var _ wntr.PreDestroyable = &UsersDaoLedis{}
	var _ UsersDao = &UsersDaoLedis{}
}

func (this *UsersDaoLedis) PostInit() error {
	cfg := lediscfg.NewConfigDefault()

	l, err := ledis.Open(cfg)
	if err != nil {
		return fmt.Errorf("Unable to create ledis dao %v", err)
	}

	this.ledis = l
	this.db, err = this.ledis.Select(0)
	if err != nil {
		return fmt.Errorf("Unable to create ledis dao %v", err)
	}

	return nil
}

func (this *UsersDaoLedis) PreDestroy() {
	this.ledis.Close()
}

func (this *UsersDaoLedis) CreateUser() (*User, error) {
	buf := make([]byte, 0)
	buf = strconv.AppendInt(buf, time.Now().Unix(), 36)
	buf = strconv.AppendInt(buf, rand.Int63(), 36)

	id := string(buf)

	u := User{
		Id:   id,
		Name: "Doe J.",
	}

	bUser, err := json.Marshal(u)

	if err != nil {
		return nil, fmt.Errorf("Unable to marshall: %v", err)
	}

	if _, err := this.db.SetNX([]byte(id), bUser); err != nil {
		return nil, err
	}

	return &u, nil
}

func (this *UsersDaoLedis) ListUsers(from int, to int) ([]*User, error) {
	result, err := this.db.Scan(ledis.KV, nil, to, false, "")

	if err != nil {
		return nil, err
	}

	r := make([]*User, len(result))

	for i, blob := range result {
		id := string(blob)
		r[i] = &User{
			Id: id,
		}
	}

	return r, nil
}

func (this *UsersDaoLedis) GetUser(id string) (*User, error) {
	bUser, err := this.db.Get([]byte(id))
	if err != nil {
		return nil, err
	}

	if bUser == nil {
		return nil, fmt.Errorf("User %v not exists", id)
	}

	u := &User{}
	if err := json.Unmarshal(bUser, u); err != nil {
		return nil, err
	}
	return u, nil
}

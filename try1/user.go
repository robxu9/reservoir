package reservoir

import (
	"github.com/stretchrcom/goweb/goweb"
)

type User struct {
	name     string
	email    string
	details  string
	projects []string
}

type UserDetails struct {
	key      string
	password string // (sha512sum + name hash)
	apikey   string // encoded (can be decoded)
}

func NewUser() User {

}

func FindUserAPI(apikey string) *User {

}

func FindUserName(name string) *User {

}

func UpdateUser(user *User) error {

}

// User API Controller
// Accessible: http://<apiloc>/user/{name}
// Super: http://<apiloc>/<apikey>/user/{name}
// ...Creation of User or Team should be done via CreateAPIController

type UserAPIController struct{}

func (cr *UserAPIController) Create(cx *goweb.Context) {
	cx.RespondWithData(TestEntity{"1", "Mat", 28})
}

func (cr *UserAPIController) Delete(id string, cx *goweb.Context) {
	cx.RespondWithOK()
}

func (cr *UserAPIController) DeleteMany(cx *goweb.Context) {
	cx.RespondWithStatus(http.StatusForbidden)
}

func (cr *UserAPIController) Read(id string, cx *goweb.Context) {

	if id == "1" {
		cx.RespondWithData(TestEntity{id, "Mat", 28})
	} else if id == "2" {
		cx.RespondWithData(TestEntity{id, "Laurie", 27})
	} else {
		cx.RespondWithNotFound()
	}

}

func (cr *UserAPIController) ReadMany(cx *goweb.Context) {
	cx.RespondWithData([]TestEntity{TestEntity{"1", "Mat", 28}, TestEntity{"2", "Laurie", 27}})
}

func (cr *UserAPIController) Update(id string, cx *goweb.Context) {
	cx.RespondWithData(TestEntity{id, "Mat", 28})
}

func (cr *UserAPIController) UpdateMany(cx *goweb.Context) {
	cx.RespondWithData(TestEntity{"1", "Mat", 28})
}

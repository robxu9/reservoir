package reservoir

import (
	"github.com/stretchrcom/goweb/goweb"
)

type CreateAPIController struct{}

// CreateAPIController mandates all creation of objects, including
// new scheduled signing and build jobs, new release requests,
// new submissions, new users, new teams, etc...
//
// Accessible by: http://<apiloc>/{apikey}/create/{type}/{parameters}/{...}
//
// The nice thing about goweb's flexible framework is that we can monitor
// which creation request wanted by the type.

func (cr *CreateAPIController) Create(cx *goweb.Context) {
	cx.PathParams["apikey"]
	cx.RespondWithData(TestEntity{"1", "Mat", 28})
}

func (cr *CreateAPIController) Delete(id string, cx *goweb.Context) {
	cx.RespondWithOK()
}

func (cr *CreateAPIController) DeleteMany(cx *goweb.Context) {
	cx.RespondWithStatus(http.StatusForbidden)
}

func (cr *CreateAPIController) Read(id string, cx *goweb.Context) {

	if id == "1" {
		cx.RespondWithData(TestEntity{id, "Mat", 28})
	} else if id == "2" {
		cx.RespondWithData(TestEntity{id, "Laurie", 27})
	} else {
		cx.RespondWithNotFound()
	}

}

func (cr *CreateAPIController) ReadMany(cx *goweb.Context) {
	cx.RespondWithData([]TestEntity{TestEntity{"1", "Mat", 28}, TestEntity{"2", "Laurie", 27}})
}

func (cr *CreateAPIController) Update(id string, cx *goweb.Context) {
	cx.RespondWithData(TestEntity{id, "Mat", 28})
}

func (cr *CreateAPIController) UpdateMany(cx *goweb.Context) {
	cx.RespondWithData(TestEntity{"1", "Mat", 28})
}

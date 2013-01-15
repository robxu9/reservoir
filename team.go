package reservoir

import (
	"github.com/stretchrcom/goweb/goweb"
)

type Team struct {
	name     string
	email    string
	members  []TeamMember
	roles    []TeamRole
	projects []string
}

type TeamMember struct {
	user string
	role TeamRole
}

type TeamRole struct {
	name         string
	repositories []TeamRepository
}

type TeamRepository struct {
	repository string
	readable   bool
	writable   bool
	buildable  bool
	adminable  bool
}

type TeamAPIController struct{}

func (cr *TeamAPIController) Create(cx *goweb.Context) {
	cx.RespondWithData(TestEntity{"1", "Mat", 28})
}

func (cr *TeamAPIController) Delete(id string, cx *goweb.Context) {
	cx.RespondWithOK()
}

func (cr *TeamAPIController) DeleteMany(cx *goweb.Context) {
	cx.RespondWithStatus(http.StatusForbidden)
}

func (cr *TeamAPIController) Read(id string, cx *goweb.Context) {

	if id == "1" {
		cx.RespondWithData(TestEntity{id, "Mat", 28})
	} else if id == "2" {
		cx.RespondWithData(TestEntity{id, "Laurie", 27})
	} else {
		cx.RespondWithNotFound()
	}

}

func (cr *TeamAPIController) ReadMany(cx *goweb.Context) {
	cx.RespondWithData([]TestEntity{TestEntity{"1", "Mat", 28}, TestEntity{"2", "Laurie", 27}})
}

func (cr *TeamAPIController) Update(id string, cx *goweb.Context) {
	cx.RespondWithData(TestEntity{id, "Mat", 28})
}

func (cr *TeamAPIController) UpdateMany(cx *goweb.Context) {
	cx.RespondWithData(TestEntity{"1", "Mat", 28})
}

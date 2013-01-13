// Copyright (c) 2013 Robert Xu. New BSD License.
package reservoir

type Owner interface {
	GetName() string
	GetEmail() string
	GetRepositories() []string
}

type User struct {
	name         string
	email        string
	details      string
	repositories []string
}

type UserDetails struct {
	key      string
	password string // (sha512sum + name hash)
	apikey   string // encoded (can be decoded)
}

type Team struct {
	name         string
	email        string
	members      []TeamMember
	roles        []TeamRole
	repositories []string
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

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetRepositories() []string {
	return u.repositories
}

func (t *Team) GetName() string {
	return t.name
}

func (t *Team) GetEmail() string {
	return t.email
}

func (t *Team) GetRepositories() []string {
	return t.repositories
}

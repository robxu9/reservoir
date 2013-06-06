package main

type Group struct {
	Username string
	Users    []*User
	Projects []*Project
}

func (g *Group) GetName() string {
	return g.Username
}

func (g *Group) GetProjects() []*Project {
	return g.Projects
}

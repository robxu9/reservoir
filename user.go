package main

type User struct {
	Username string
	Email    string
	Projects []*Project
}

func (u *User) GetName() string {
	return u.Username
}

func (u *User) GetProjects() []*Project {
	return u.Projects
}

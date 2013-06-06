package main

type Ident interface {
	GetName() string
	GetProjects() []*Project
}

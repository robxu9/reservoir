package main

type Identity interface {
	GetName() string
	GetProjects() []*Project
}

package main

type Project struct {
	Name       string
	Meta       map[string]string
	Repo       map[string]*Repository
	Components map[string]*Component
}

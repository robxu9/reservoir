package main

type Job struct {
	Project      *Project
	Repository   *Repository
	Architecture string
	Component    *Component
	Results      map[string]string // Filenames => filestore id
}

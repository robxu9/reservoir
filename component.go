package main

type Component struct {
	Name          string
	GitRepository string
	Meta          map[string]string
	Jobs          map[string]*Job // Formatted with time/Time.String
}

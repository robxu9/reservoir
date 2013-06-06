package main

type Type struct {
	Name         string
	DisplayName  string
	Description  string
	ApplicableSh string // determine if component is applicable for this repo
	BuildSh      string // build a component
	PublishSh    string // given new components, publish with old ones
}

package main

type Repository struct {
	Name          string
	Type          *Type
	Architectures []string // pass into type's ShBuild and ShPublish as ENV var
}

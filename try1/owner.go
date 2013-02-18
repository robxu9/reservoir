package reservoir

type Owner interface {
	GetName() string
	GetEmail() string
	GetProjects() []string
	IsTeam() bool
	IsUser() bool
}

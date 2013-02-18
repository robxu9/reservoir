package reservoir

type Project struct {
	name    string
	owner   Owner
	private bool
	parent  string // no parent == ""
}

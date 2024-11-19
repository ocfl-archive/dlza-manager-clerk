package graph

type Version struct {
	Name    string
	Address string
	Created string
	Message string
}

type Versions map[string]Version

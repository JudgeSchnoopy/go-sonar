package client

type Response struct {
	Name         string
	Status       string
	Dependencies []Dependency
}

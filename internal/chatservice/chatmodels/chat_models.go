package chatmodels

import "fmt"

// ChatUser data structure
type ChatUser struct {
	Provider string
	ID       string
	Name     string
}

func (u ChatUser) FullyQualifiedName() string {
	return fmt.Sprintf("%s-%s-%s", u.Provider, u.Name, u.ID)
}

// Channel data structure
type Channel struct {
	ID   string
	Name string
}

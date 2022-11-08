package services

type SlackUser struct {
	ID   string
	Name string
}

func NewSlackUser(id, name string) SlackUser {
	return SlackUser{
		ID:   id,
		Name: name,
	}
}

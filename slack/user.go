package slack

type User struct {
	ID   string
	Name string
}

func (u User) String() string {
	return "<@" + u.ID + "|" + u.Name + ">"
}

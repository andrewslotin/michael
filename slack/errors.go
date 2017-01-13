package slack

type NoSuchUserError struct {
	Name string
}

func (e NoSuchUserError) Error() string {
	return "there is no user with username '" + e.Name + "' in the team"
}

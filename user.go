package nflpickem

// User represents a user of the NFL Pickem' Pool
type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Admin     bool   `json:"admin"`
}

func (u User) Equal(other User) bool {
	return u.Email == other.Email
}

type PasswordUpdater interface {
	UpdatePassword(username string, oldPassword string, newPassword string) error
}

type CredentialChecker interface {
	CheckCredentials(username string, password string) (User, error)
}

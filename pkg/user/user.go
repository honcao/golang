package user

// User struct
type User struct {
	FirstName string
	LastName  string
}

// GetFullName return the full name of user
func (u User) GetFullName() string {
	if len(u.FirstName) > 0 && len(u.LastName) > 0 {
		return u.LastName + " " + u.FirstName
	}

	if len(u.FirstName) > 0 {
		return u.FirstName
	}

	if len(u.LastName) > 0 {
		return u.LastName
	}

	return ""
}

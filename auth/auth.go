package auth

//User represents an authenticated user
type User struct {
	Username    string
	DisplayName string
}

//Auth represents an authentication mechanism
type Auth interface {
	//Authenticate authenticates the given credentials and returns the User associated with the account if successful,
	//or nil if not. If an error occurs it is returned.
	Authenticate(username, password string) (user *User, err error)
}

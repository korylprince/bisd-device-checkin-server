package session

import "github.com/korylprince/bisd-device-checkin-server/v2/auth"

//Session represents an authenticated session
type Session auth.User

//Store is a session storage mechanism
type Store interface {
	//Create creates and returns a session id for the given session
	//or an error if one occurred
	Create(s *Session) (id string, err error)
	//Create returns the session for the given id or nil if it doesn't exist
	//or an error if one occurred
	Check(id string) (*Session, error)
}

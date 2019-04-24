package httpapi

import (
	"io"

	"github.com/korylprince/bisd-device-checkin-server/auth"
	"github.com/korylprince/bisd-device-checkin-server/db"
	"github.com/korylprince/bisd-device-checkin-server/session"
)

//Server represents shared resources
type Server struct {
	db           db.DB
	auth         auth.Auth
	sessionStore session.Store
	output       io.Writer
}

//NewServer returns a new server with the given resources
func NewServer(db db.DB, auth auth.Auth, sessionStore session.Store, output io.Writer) *Server {
	return &Server{db: db, auth: auth, sessionStore: sessionStore, output: output}
}

package httpapi

import (
	"io"

	"github.com/korylprince/bisd-device-checkin-server/v2/auth"
	"github.com/korylprince/bisd-device-checkin-server/v2/db"
	"github.com/korylprince/bisd-device-checkin-server/v2/session"
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

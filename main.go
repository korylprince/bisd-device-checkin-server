package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/korylprince/bisd-device-checkin-server/v2/auth/ad"
	"github.com/korylprince/bisd-device-checkin-server/v2/db/sql"
	"github.com/korylprince/bisd-device-checkin-server/v2/httpapi"
	"github.com/korylprince/bisd-device-checkin-server/v2/session/memory"
	auth "github.com/korylprince/go-ad-auth/v3"
)

func main() {
	db, err := sql.New(config.SQLDriver, config.SQLDSN)
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	adAuth := ad.New(
		&auth.Config{
			Server:   config.LDAPServer,
			Port:     config.LDAPPort,
			BaseDN:   config.LDAPBaseDN,
			Security: config.ldapSecurity,
		},
		[]string{config.LDAPGroup},
	)

	sessionStore := memory.New(time.Minute * time.Duration(config.SessionExpiration))

	httpapi.Debug = config.Debug

	s := httpapi.NewServer(db, adAuth, sessionStore, os.Stdout)

	log.Println("Listening on:", config.ListenAddr)

	log.Println(http.ListenAndServe(config.ListenAddr, http.StripPrefix(config.Prefix, s.Router())))
}

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/korylprince/bisd-device-checkin-server/api"
	"github.com/korylprince/bisd-device-checkin-server/httpapi"
	auth "github.com/korylprince/go-ad-auth"
)

func main() {
	db, err := sql.Open(config.SQLDriver, config.SQLDSN)
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	adConfig := &api.AuthConfig{
		ADConfig: &auth.Config{
			Server:   config.LDAPServer,
			Port:     config.LDAPPort,
			BaseDN:   config.LDAPBaseDN,
			Security: config.ldapSecurity,
		},
		Group: config.LDAPGroup,
	}

	s := httpapi.NewMemorySessionStore(time.Minute * time.Duration(config.SessionExpiration))

	r := httpapi.NewRouter(os.Stdout, adConfig, s, db)

	chain := handlers.CompressHandler(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Origin", "X-Session-Key"}),
	)(http.StripPrefix(config.Prefix, r)))

	log.Println("Listening on:", config.ListenAddr)
	log.Println(http.ListenAndServe(config.ListenAddr, chain))
}

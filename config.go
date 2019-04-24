package main

import (
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	auth "gopkg.in/korylprince/go-ad-auth.v2"
)

//Config represents options given in the environment
type Config struct {
	SessionExpiration int `default:"15"` //in minutes

	LDAPServer   string `required:"true"`
	LDAPPort     int    `default:"389" required:"true"`
	LDAPBaseDN   string `required:"true"`
	LDAPGroup    string `required:"true"`
	LDAPSecurity string `default:"none" required:"true"`
	ldapSecurity auth.SecurityType

	SQLDriver string `required:"true"`
	SQLDSN    string `required:"true"`

	ListenAddr string `default:":8080" required:"true"` //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash
	Debug      bool   `default:"false"` //return debugging information to client
}

var config = &Config{}

func init() {
	err := envconfig.Process("INVENTORY", config)
	if err != nil {
		log.Fatalln("Error reading configuration from environment:", err)
	}

	switch strings.ToLower(config.LDAPSecurity) {
	case "", "none":
		config.ldapSecurity = auth.SecurityNone
	case "tls":
		config.ldapSecurity = auth.SecurityTLS
	case "starttls":
		config.ldapSecurity = auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid INVENTORY_LDAPSECURITY:", config.LDAPSecurity)
	}

	if config.SQLDriver == "mysql" && !strings.Contains(config.SQLDSN, "?parseTime=true") {
		log.Fatalln("mysql DSN must contain \"?parseTime=true\"")
	}
}

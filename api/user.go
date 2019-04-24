package api

import (
	"errors"
	"fmt"

	auth "gopkg.in/korylprince/go-ad-auth.v2"
)

//AuthConfig holds configuration for connecting to an authentication source
type AuthConfig struct {
	ADConfig *auth.Config
	Groups   []string
}

//User represents an Active Directory User
type User struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

//Authenticate authenticates the given username and password against the given config,
//returning user information if successful, nil if unsuccessful, or an error if one occurred.
func Authenticate(config *AuthConfig, username, password string) (*User, error) {
	status, entry, groups, err := auth.AuthenticateExtended(config.ADConfig, username, password, []string{"displayName"}, config.Groups)
	if err != nil {
		return nil, fmt.Errorf("Error attempting to authenticate: %v", err)
	}

	if !status {
		return nil, errors.New("User not authenticated")
	}

	if len(groups) == 0 {
		return nil, errors.New("User not in any authorized groups")
	}

	return &User{Username: username, DisplayName: entry.GetAttributeValue("displayName")}, nil
}

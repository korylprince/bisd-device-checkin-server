package ad

import (
	"fmt"

	"github.com/korylprince/bisd-device-checkin-server/auth"
	adauth "github.com/korylprince/go-ad-auth/v3"
)

//Auth represents an Active Directory authentication mechanism
type Auth struct {
	config *adauth.Config
	groups []string
}

//New returns a new *Auth with the given configuration and permissions mapping
func New(config *adauth.Config, groups []string) *Auth {
	return &Auth{config: config, groups: groups}
}

//Authenticate authenticates the given credentials and returns the User associated with the account if successful,
//or nil if not. If an error occurs it is returned.
func (a *Auth) Authenticate(username, password string) (user *auth.User, err error) {
	status, entry, groups, err := adauth.AuthenticateExtended(a.config, username, password, []string{"displayName"}, a.groups)
	if err != nil {
		return nil, fmt.Errorf("Error attempting to authenticate as %s: %v", username, err)
	}

	if !status {
		return nil, nil
	}

	if len(groups) == 0 {
		return nil, nil
	}

	return &auth.User{
		Username:    username,
		DisplayName: entry.GetAttributeValue("displayName"),
	}, nil
}

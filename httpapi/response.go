package httpapi

import "github.com/korylprince/bisd-device-checkin-server/api"

//AuthenticateResponse is a successful authentication response including the session key and User
type AuthenticateResponse struct {
	SessionKey string    `json:"session_key"`
	User       *api.User `json:"user"`
}

//CheckinResponse holds a charge ID if one exists
type CheckinResponse struct {
	ChargeID int64 `json:"charge_id"`
}

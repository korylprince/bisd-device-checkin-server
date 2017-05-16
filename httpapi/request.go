package httpapi

import (
	"fmt"
	"strings"
)

//AuthenticateRequest is an username/password authentication request
type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Charge represents a charge for a damaged/missing device
type Charge struct {
	Description string  `json:"description"`
	Value       float32 `json:"value"`
}

//Charges is a list of Charges
type Charges []Charge

//Marshal marshals charges to the inventory database format
func (charges Charges) Marshal() string {
	s := []string{"|"}
	for _, c := range charges {
		//clean input
		desc := strings.Replace(strings.Replace(strings.TrimSpace(c.Description), "|", "", -1), ":", "", -1)
		s = append(s, fmt.Sprintf("%s:%.2f|", desc, c.Value))
	}
	return strings.Join(s, "")
}

//CheckinRequest is a request to check in a device in the inventory, with any charges that may apply
type CheckinRequest struct {
	Charges Charges `json:"charges"`
}

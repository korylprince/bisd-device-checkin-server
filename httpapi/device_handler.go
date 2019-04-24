package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-device-checkin-server/api"
)

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

//GET /devices/:bagTag
func handleReadDevice(w http.ResponseWriter, r *http.Request) *handlerResponse {
	bagTag := mux.Vars(r)["bagTag"]

	device, err := api.ReadDeviceByBagTag(r.Context(), bagTag)
	if resp := checkAPIError(err); resp != nil {
		return resp
	}
	if device == nil {
		return handleError(http.StatusNotFound, errors.New("Could not find device"))
	}

	return &handlerResponse{Code: http.StatusOK, Body: device}
}

//POST /devices/:bagTag/checkin
func handleCheckinDevice(w http.ResponseWriter, r *http.Request) *handlerResponse {
	type request struct {
		Charges Charges `json:"charges"`
	}

	type response struct {
		ChargeID int64 `json:"charge_id"`
	}

	bagTag := mux.Vars(r)["bagTag"]

	//read charges
	var req *request
	d := json.NewDecoder(r.Body)

	err := d.Decode(&req)
	if err != nil {
		return handleError(http.StatusBadRequest, fmt.Errorf("Could not decode JSON: %v", err))
	}

	device, err := api.ReadDeviceByBagTag(r.Context(), bagTag)
	if resp := checkAPIError(err); resp != nil {
		return resp
	}
	if device == nil {
		return handleError(http.StatusNotFound, errors.New("Could not find device"))
	}

	//validate
	if device.User == "" {
		return handleError(http.StatusBadRequest, fmt.Errorf("There is no user for bag tag: %s", bagTag))
	}

	if device.Status != "Checked Out" {
		return handleError(http.StatusBadRequest, fmt.Errorf("Device is not Checked Out for bag tag: %s", bagTag))
	}

	//update
	user := r.Context().Value(api.UserKey).(*api.User)

	device.Notes = device.Notes + fmt.Sprintf("\n%s %s: Checked in Bag Tag %s from %s\n",
		time.Now().Format("01/02/06"),
		user.DisplayName,
		bagTag,
		device.User,
	)

	//if charges add note
	if req != nil && len(req.Charges) > 0 {
		text := req.Charges.Marshal()
		device.Notes = device.Notes + fmt.Sprintf("\tCharges: %s\n",
			//pretty format
			strings.Replace(strings.Replace(text[1:len(text)-1], "|", ", ", -1), ":", ": $", -1),
		)
	}

	devUser := device.User
	device.User = ""
	device.Status = "Storage"

	err = device.Update(r.Context())
	if resp := checkAPIError(err); resp != nil {
		return resp
	}

	//create charge
	if req != nil && len(req.Charges) > 0 {
		c := &api.Charge{
			InventoryNumber: device.InventoryNumber,
			User:            devUser,
			AmountPaid:      0,
			Charges:         req.Charges.Marshal(),
			Notes:           fmt.Sprintf("Charges created %s by %s.\n", time.Now().Format("01/02/06"), user.DisplayName),
		}

		id, err := api.CreateCharge(r.Context(), c)
		if resp := checkAPIError(err); resp != nil {
			return resp
		}

		return &handlerResponse{Code: http.StatusOK, Body: &response{ChargeID: id}}
	}

	return &handlerResponse{Code: http.StatusOK, Body: &response{ChargeID: 0}}
}

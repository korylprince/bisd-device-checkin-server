package httpapi

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/korylprince/bisd-device-checkin-server/db"
	"github.com/korylprince/bisd-device-checkin-server/session"
)

func (s *Server) readDevice(r *http.Request, tx *sql.Tx) (int, interface{}) {
	type errResponse struct {
		Error string `json:"error"`
	}

	bagTag := mux.Vars(r)["id"]

	device, err := s.db.ReadDevice(tx, bagTag)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if device == nil {
		return http.StatusNotFound, errors.New("Could not find device")
	}

	if device.Status != "Checked Out" {
		return http.StatusBadRequest, &errResponse{Error: "Device not checked out"}
	}

	if strings.TrimSpace(device.User) == "" {
		return http.StatusBadRequest, &errResponse{Error: "Device has no assigned user"}
	}

	return http.StatusOK, device
}

func (s *Server) checkinDevice(r *http.Request, tx *sql.Tx) (int, interface{}) {
	type request struct {
		Charges db.Charges `json:"charges"`
		Missing bool       `json:"missing"`
		Notes   string     `json:"notes"`
	}

	type response struct {
		ChargeID int64 `json:"charge_id"`
	}

	bagTag := mux.Vars(r)["id"]

	//read charges
	req := new(request)

	if err := jsonRequest(r, req); err != nil {
		return http.StatusBadRequest, err
	}

	//read device
	device, err := s.db.ReadDevice(tx, bagTag)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if device == nil {
		return http.StatusNotFound, errors.New("Could not find device")
	}

	//validate
	if device.User == "" {
		return http.StatusBadRequest, fmt.Errorf("There is no user for bag tag: %s", bagTag)
	}

	if device.Status != "Checked Out" {
		return http.StatusBadRequest, fmt.Errorf("Device is not Checked Out for bag tag: %s", bagTag)
	}

	//get inventory user
	user := r.Context().Value(contextKeyUser).(*session.Session)

	//update notes
	var fmtStr string
	if req.Missing {
		fmtStr = "\n%s %s: Bag Tag %s from %s is Missing\n"
	} else {
		fmtStr = "\n%s %s: Checked in Bag Tag %s from %s\n"
	}

	device.Notes += fmt.Sprintf(fmtStr,
		time.Now().Format("01/02/06"),
		user.DisplayName,
		bagTag,
		device.User,
	)

	//if charges add note
	if len(req.Charges) > 0 {
		text := req.Charges.Marshal()
		device.Notes += fmt.Sprintf("\tCharges: %s\n",
			//pretty format
			strings.Replace(strings.Replace(text[1:len(text)-1], "|", ", ", -1), ":", ": $", -1),
		)
	}

	//add extra notes
	if len(req.Notes) > 0 {
		device.Notes += "\tNotes:\n"
		for _, line := range strings.Split(req.Notes, "\n") {
			device.Notes += fmt.Sprintf("\t\t%s\n", line)
		}
	}

	//update fields
	devUser := device.User
	device.User = ""
	device.Status = "Storage"
	if req.Missing {
		device.Status = "Missing"
	}

	//update device
	err = s.db.UpdateDevice(tx, device)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	//create charge
	if req.Charges.Total() > 0 {

		id, err := s.db.CreateCharge(tx, req.Charges, device.InventoryNumber, devUser,
			fmt.Sprintf("Charges created %s by %s.\n%s", time.Now().Format("01/02/06"), user.DisplayName, req.Notes),
		)

		if err != nil {
			return http.StatusInternalServerError, err
		}

		return http.StatusOK, &response{ChargeID: id}
	}

	return http.StatusOK, &response{ChargeID: 0}
}

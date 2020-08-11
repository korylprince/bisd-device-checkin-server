package db

import (
	"database/sql"
	"fmt"
	"strings"
)

//Charge represents a charge for a damaged/missing device
type Charge struct {
	Description string  `json:"description"`
	Amount      float32 `json:"amount"`
}

//Charges is a list of Charges
type Charges []Charge

//Marshal marshals charges to the inventory database format
func (charges Charges) Marshal() string {
	s := []string{"|"}
	for _, c := range charges {
		//clean input
		desc := strings.Replace(strings.Replace(strings.TrimSpace(c.Description), "|", ";", -1), ":", ";", -1)
		s = append(s, fmt.Sprintf("%s:%.2f|", desc, c.Amount))
	}
	return strings.Join(s, "")
}

//Total returns the total charge amount
func (charges Charges) Total() float32 {
	var t float32
	for _, c := range charges {
		t += c.Amount
	}
	return t
}

//Device represents a device in the inventory database
type Device struct {
	ID              int64  `json:"id"`
	InventoryNumber string `json:"inventory_number"`
	SerialNumber    string `json:"serial_number"`
	BagTag          string `json:"bag_tag"`
	Status          string `json:"status"`
	Model           string `json:"model"`
	User            string `json:"user"`
	Notes           string `json:"notes"`
}

//DB represents a database
type DB interface {
	//Begin returns a transaction for the database or an error if one occurred
	Begin() (*sql.Tx, error)

	//CreateCharge creates a charge in the database with the given info and returns the ID of the charges or an error if one occurred
	CreateCharge(tx *sql.Tx, charges Charges, inventoryNumber, user, notes string) (id int64, err error)

	//ReadDevice returns the Device with with the given bagTag or an error if one occurred
	ReadDevice(tx *sql.Tx, bagTag string) (*Device, error)

	//UpdateDevice updates the given Device by the given username or returns an error if one occurred
	UpdateDevice(tx *sql.Tx, device *Device, username string) error
}

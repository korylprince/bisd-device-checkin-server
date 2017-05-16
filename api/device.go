package api

import (
	"context"
	"database/sql"
	"fmt"
)

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

//ReadDeviceByBagTag returns the device with the given bag tag
func ReadDeviceByBagTag(ctx context.Context, bagTag string) (*Device, error) {
	tx := ctx.Value(TransactionKey).(*sql.Tx)

	device := &Device{BagTag: bagTag}

	row := tx.QueryRow("SELECT id, inventory_number, serial_number, status, model, user, notes FROM devices WHERE bag_tag=? AND model = 'C740-C4PE';", bagTag)
	err := row.Scan(
		&(device.ID),
		&(device.InventoryNumber),
		&(device.SerialNumber),
		&(device.Status),
		&(device.Model),
		&(device.User),
		&(device.Notes),
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, &Error{Description: fmt.Sprintf("Could not query Device(%s)", bagTag), Err: err}
	}

	return device, nil
}

//Update updates the status, user, and notes for the given Device
func (d *Device) Update(ctx context.Context) error {
	tx := ctx.Value(TransactionKey).(*sql.Tx)

	_, err := tx.Exec("UPDATE devices SET status=?, user=?, notes=? WHERE id=?",
		d.Status,
		d.User,
		d.Notes,
		d.ID,
	)
	if err != nil {
		return &Error{Description: "Could not update Device", Err: err}
	}

	return nil
}

package sql

import (
	"database/sql"
	"fmt"

	"github.com/korylprince/bisd-device-checkin-server/db"
)

//DB represents a SQL database
type DB struct {
	db *sql.DB
}

//New returns a new DB with the given driver and DSN or an error if one occurred
func New(driver, dsn string) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

//Begin returns a transaction for the database or an error if one occurred
func (d *DB) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

//CreateCharge creates a charge in the database with the given info and returns the ID of the charges or an error if one occurred
func (d *DB) CreateCharge(tx *sql.Tx, charges db.Charges, inventoryNumber, user, notes string) (id int64, err error) {
	res, err := tx.Exec("INSERT INTO charges(inventory_number, user, amount_paid, charges, notes) VALUES(?, ?, ?, ?, ?);",
		inventoryNumber,
		user,
		0.0,
		charges.Marshal(),
		notes,
	)
	if err != nil {
		return 0, fmt.Errorf("Unable to insert Charge: %v", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Unable to fetch Charge ID: %v", err)
	}

	return id, nil
}

//ReadDevice returns the Device with with the given bagTag or an error if one occurred
func (d *DB) ReadDevice(tx *sql.Tx, bagTag string) (*db.Device, error) {
	device := &db.Device{BagTag: bagTag}

	row := tx.QueryRow("SELECT id, inventory_number, serial_number, status, model, user, notes FROM devices WHERE bag_tag=? AND (model = 'C732T-C8VY' OR model = 'Chromebook 3100');", bagTag)

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
		return nil, fmt.Errorf("Unable to query Device(%s): %v", bagTag, err)
	}

	return device, nil
}

//UpdateDevice updates the given Device or returns an error if one occurred
func (d *DB) UpdateDevice(tx *sql.Tx, device *db.Device) error {
	_, err := tx.Exec("UPDATE devices SET status=?, user=?, notes=? WHERE id=?",
		device.Status,
		device.User,
		device.Notes,
		device.ID,
	)
	if err != nil {
		return fmt.Errorf("Unable to update Device: %v", err)
	}

	return nil
}

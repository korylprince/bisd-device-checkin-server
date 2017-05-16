package api

import (
	"context"
	"database/sql"
)

//Charge represents a charge in the inventory database
type Charge struct {
	ID              int64
	InventoryNumber string
	User            string
	AmountPaid      float32
	Charges         string
	Notes           string
}

//CreateCharge creates a new charge in the inventory database and returns its ID
func CreateCharge(ctx context.Context, c *Charge) (id int64, err error) {

	tx := ctx.Value(TransactionKey).(*sql.Tx)

	res, err := tx.Exec("INSERT INTO charges(inventory_number, user, amount_paid, charges, notes) VALUES(?, ?, ?, ?, ?);",
		c.InventoryNumber,
		c.User,
		c.AmountPaid,
		c.Charges,
		c.Notes,
	)
	if err != nil {
		return 0, &Error{Description: "Could not insert Charge", Err: err}
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, &Error{Description: "Could not fetch Charge ID", Err: err}
	}

	return id, nil
}

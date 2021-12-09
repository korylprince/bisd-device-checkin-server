package httpapi

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/korylprince/bisd-device-checkin-server/v2/db"
)

type txReturnHandlerFunc func(*http.Request, *sql.Tx) (int, interface{})

func withTX(db db.DB, next txReturnHandlerFunc) returnHandlerFunc {
	return func(r *http.Request) (int, interface{}) {
		tx, err := db.Begin()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Unable to start database transaction: %v", err)
		}

		status, body := next(r, tx)

		if status != http.StatusOK {
			if err = tx.Rollback(); err != nil {
				if pErr, ok := body.(error); ok {
					return http.StatusInternalServerError, fmt.Errorf("Unable to rollback database transaction: %v; Previous error: HTTP %d %s: %v", err, status, http.StatusText(status), pErr)
				}
				return http.StatusInternalServerError, fmt.Errorf("Unable to rollback database transaction: %v", err)
			}
			return status, body
		}

		if err = tx.Commit(); err != nil {
			if pErr, ok := body.(error); ok {
				return http.StatusInternalServerError, fmt.Errorf("Unable to commit database transaction: %v; Previous error: HTTP %d %s: %v", err, status, http.StatusText(status), pErr)
			}
			return http.StatusInternalServerError, fmt.Errorf("Unable to commit database transaction: %v", err)
		}

		return status, body
	}
}

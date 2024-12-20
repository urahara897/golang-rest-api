package utils

import (
	"database/sql"
	"time"
)

// RetryOnBusy executes a database operation with retries if SQLite returns a "database is locked" error
func RetryOnBusy(operation func() error) error {
	maxRetries := 5
	backoff := 100 * time.Millisecond

	var err error
	for i := 0; i < maxRetries; i++ {
		err = operation()
		
		// Check if the error is a SQLite "database is locked" error
		if err == nil || err != sql.ErrConnDone {
			return err
		}

		// Wait before retrying
		time.Sleep(backoff)
		
		// Exponential backoff
		backoff *= 2
	}

	return err
} 